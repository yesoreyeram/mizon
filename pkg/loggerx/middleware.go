// HTTP logging middleware
package loggerx

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Config struct {
	LogRequestBody bool
	MaxBody        int
}

type statusRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	n, err := sr.ResponseWriter.Write(b)
	sr.size += n
	return n, err
}

// Middleware returns a standard http middleware that logs requests and responses.
func Middleware(cfg Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			var bodyPreview string
			if cfg.LogRequestBody && (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch) && r.Body != nil {
				data, _ := io.ReadAll(io.LimitReader(r.Body, int64(cfg.MaxBody)))
				bodyPreview = string(data)
				r.Body = io.NopCloser(bytes.NewReader(data))
			}
			rec := &statusRecorder{ResponseWriter: w, status: 200}
			Debugw("request", map[string]interface{}{
				"method": r.Method,
				"uri":    r.URL.RequestURI(),
				"remote": r.RemoteAddr,
				"ua":     r.UserAgent(),
				"len":    r.Header.Get("Content-Length"),
				"body":   bodyPreview,
			})
			next.ServeHTTP(rec, r)
			dur := time.Since(start)
			Infow("response", map[string]interface{}{
				"method":   r.Method,
				"path":     r.URL.Path,
				"status":   rec.status,
				"size":     rec.size,
				"duration": dur.String(),
			})
		})
	}
}

// Env helpers
func EnvBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func EnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
