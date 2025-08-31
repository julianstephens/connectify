package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// LoggingMiddlewareLogrus returns a middleware that logs requests using the provided logrus.Logger.
// It logs method, path (RequestURI), status, bytes written, remote address, duration and request-id.
func LoggingMiddlewareLogrus(logger *logrus.Logger) func(http.Handler) http.Handler {
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

			logger.WithFields(logrus.Fields{
				"method":      r.Method,
				"uri":         r.RequestURI,
				"status":      rw.status,
				"bytes":       rw.written,
				"remote_addr": r.RemoteAddr,
				"duration":    duration.String(),
				"request_id":  reqID,
			}).Info("http_request")
		})
	}
}
