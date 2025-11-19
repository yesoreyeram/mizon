module product-generator

go 1.23.0

require mizon/loggerx v0.0.0

require (
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
)

replace mizon/loggerx => ../../pkg/loggerx
