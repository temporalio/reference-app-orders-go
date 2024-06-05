package instrumentation

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type instrumentedResponseWriter struct {
	http.ResponseWriter
	status int
}

func (r *instrumentedResponseWriter) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Middleware logs the status code of each request.
func Middleware(logger *slog.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			iw := instrumentedResponseWriter{w, http.StatusOK}

			h.ServeHTTP(&iw, r)

			route := chi.RouteContext(r.Context()).RoutePattern()

			logger.Info(
				fmt.Sprintf("%d %s %s", iw.status, r.Method, r.URL.Path),
				"status", iw.status,
				"method", r.Method,
				"route", route,
				"path", r.URL.Path,
			)
		})
	}
}

// ServeMetrics starts an HTTP server to serve Prometheus metrics.
func ServeMetrics(endpoint string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr:    endpoint,
		Handler: mux,
	}

	slog.Info("Starting Prometheus exporter", "endpoint", fmt.Sprintf("http://%s/metrics", endpoint))

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("error serving metrics", "error", err)
		}
	}()
}
