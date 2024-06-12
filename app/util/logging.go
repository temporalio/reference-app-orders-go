package util

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type instrumentedResponseWriter struct {
	http.ResponseWriter
	status int
}

func (r *instrumentedResponseWriter) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware is a http middleware that logs the request method, status code and path.
func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		iw := instrumentedResponseWriter{w, http.StatusOK}

		next.ServeHTTP(&iw, r)

		level := slog.LevelDebug
		if iw.status >= 500 {
			level = slog.LevelError
		}

		logger.Log(
			context.Background(), level,
			fmt.Sprintf("%d %s %s", iw.status, r.Method, r.URL.Path),
			"method", r.Method, "status", iw.status, "path", r.URL.Path,
		)
	})
}
