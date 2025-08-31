package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// LoggingMiddlewareZap returns a middleware that logs requests using the provided zap.Logger.
// It logs method, path (RequestURI), status, bytes written, remote address, duration and request-id.
func LoggingMiddlewareZap(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w}

			next.ServeHTTP(rw, r)

			// Default status to 200 if none set
			if rw.status == 0 {
				rw.status = http.StatusOK
			}
			duration := time.Since(start)
			reqID := GetReqID(r.Context())

			logger.Info("http_request",
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.Int("status", rw.status),
				zap.Int("bytes", rw.written),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Duration("duration", duration),
				zap.String("request_id", reqID),
			)
		})
	}
}
