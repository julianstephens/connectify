package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey string

const ctxKeyRequestID ctxKey = "request_id"
const RequestIDHeader = "X-Request-ID"

// RequestIDMiddleware ensures each request has a request id. If the incoming request
// already has X-Request-ID it uses that, otherwise it generates a new UUIDv4.
// The request id is added to the request context and also written to the response header.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(RequestIDHeader)
		if id == "" {
			id = uuid.New().String()
		}
		// set on response for clients
		w.Header().Set(RequestIDHeader, id)
		// add to context for downstream handlers/logging
		ctx := context.WithValue(r.Context(), ctxKeyRequestID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetReqID extracts request id from context. Returns empty string if not present.
func GetReqID(ctx context.Context) string {
	v := ctx.Value(ctxKeyRequestID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
