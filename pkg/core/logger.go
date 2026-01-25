package core

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type loggerKeyType struct{}

var loggerKey = loggerKeyType{}

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		log := LoggerFromContext(r.Context())

		log.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration", duration,
			"remote_ip", r.RemoteAddr,
		)
	})
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

func LoggerContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		reqID := middleware.GetReqID(r.Context())

		logger := slog.With(
			"request_id", reqID,
		)

		ctx := context.WithValue(r.Context(), loggerKey, logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
