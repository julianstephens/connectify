package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"

	"github.com/julianstephens/connectify/server/internal/config"
)

var (
	logger          = config.GetLogger()
	ErrInvalidToken = errors.New("invalid token")
)

// validateToken validates the JWT token and returns the claims if valid.
func validateToken(token string) (jwt.MapClaims, error) {
	k, err := keyfunc.NewDefaultCtx(context.Background(), []string{fmt.Sprintf("https://%s/.well-known/jwks.json", config.AppConfig.Auth0Domain)})
	if err != nil {
		return nil, fmt.Errorf("failed to create a keyfunc from the jwks URL: %v", err)
	}

	parsed, err := jwt.Parse(token, k.Keyfunc)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, fmt.Errorf("malformed token: %w", ErrInvalidToken)
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, fmt.Errorf("invalid token signature: %w", ErrInvalidToken)
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, fmt.Errorf("token expired: %w", ErrInvalidToken)
		default:
			return nil, fmt.Errorf("failed to parse jwt: %v", err)
		}
	}

	if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
		iss, err := claims.GetIssuer()
		if err != nil || iss != fmt.Sprintf("https://%s/", config.AppConfig.Auth0Domain) {
			return nil, fmt.Errorf("invalid token issuer: %w", ErrInvalidToken)
		}

		aud, err := claims.GetAudience()
		if err != nil || !slices.Contains(aud, config.AppConfig.Auth0Audience) {
			return nil, fmt.Errorf("invalid token audience: %w", ErrInvalidToken)
		}

		return claims, nil
	} else {
		return nil, fmt.Errorf("failed to parse token claims: %w", ErrInvalidToken)
	}
}

func AuthGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Error("missing or invalid Authorization header")
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" {
			logger.Error("missing token")
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		_, err := validateToken(tokenString)
		if err != nil {
			logger.Error(err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
