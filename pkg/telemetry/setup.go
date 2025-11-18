package telemetry

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

// Setup initializes OpenTelemetry tracing provider with OTLP exporter to Tempo.
// Controlled by env vars:
//
//	TRACING_ENABLED (default: true)
//	OTEL_EXPORTER_OTLP_ENDPOINT (default: http://tempo:4317)
//	OTEL_EXPORTER_OTLP_PROTOCOL (grpc|http, default: grpc)
//	OTEL_SAMPLER (always_on|always_off|parentbased_always_on, default: parentbased_always_on)
//	OTEL_ENV (maps to deployment.environment)
func Setup(serviceName string) (func(context.Context) error, error) {
	if strings.EqualFold(os.Getenv("TRACING_ENABLED"), "false") {
		// still set a noop provider to keep instrumentation happy
		otel.SetTracerProvider(sdktrace.NewTracerProvider())
		otel.SetTextMapPropagator(propagation.TraceContext{})
		return func(ctx context.Context) error { return nil }, nil
	}

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		// default to tempo service name in docker-compose
		endpoint = "http://tempo:4317"
	}
	proto := strings.ToLower(strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")))
	if proto == "" {
		proto = "grpc"
	}
	var exp *otlptrace.Exporter
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch proto {
	case "http", "http/protobuf":
		exp, err = otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(strings.TrimPrefix(endpoint, "http://")), otlptracehttp.WithInsecure())
	default: // grpc
		// Strip http:// if present for grpc dialer
		cleaned := strings.TrimPrefix(endpoint, "http://")
		cleaned = strings.TrimPrefix(cleaned, "https://")
		exp, err = otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(cleaned), otlptracegrpc.WithInsecure())
	}
	if err != nil {
		return nil, err
	}

	env := os.Getenv("OTEL_ENV")
	if env == "" {
		env = os.Getenv("ENVIRONMENT")
		if env == "" {
			env = "dev"
		}
	}

	res, _ := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.DeploymentEnvironment(env),
	))

	// sampler
	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))
	switch strings.ToLower(os.Getenv("OTEL_SAMPLER")) {
	case "always_on":
		sampler = sdktrace.AlwaysSample()
	case "always_off":
		sampler = sdktrace.NeverSample()
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	shutdown := func(ctx context.Context) error {
		return tp.Shutdown(ctx)
	}
	return shutdown, nil
}

// MuxMiddleware returns gorilla/mux middleware to create spans per HTTP request.
func MuxMiddleware(service string) func(h http.Handler) http.Handler {
	return otelmux.Middleware(service)
}

// NewHTTPClient returns an http.Client with otelhttp transport for propagating trace headers.
func NewHTTPClient() *http.Client {
	return &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
}

// OtelTransport returns an http.RoundTripper that injects trace context into outbounds.
func OtelTransport() http.RoundTripper {
	return otelhttp.NewTransport(http.DefaultTransport)
}

// RegisterPostgres returns the base driver unchanged if SQL instrumentation is unavailable.
// Keeping the signature allows services to call this without code changes.
func RegisterPostgres(baseDriver string) (string, error) { return baseDriver, nil }
