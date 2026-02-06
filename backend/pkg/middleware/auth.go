package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type AuthConfig struct {
	JWKSURL string
	Issuer  string
}

type AuthMiddleware struct {
	jwks   keyfunc.Keyfunc
	issuer string
	logger *slog.Logger
}

// UserClaims holds standard OIDC claims we care about
type UserClaims struct {
	jwt.RegisteredClaims
	Email string   `json:"email,omitempty"`
	Roles []string `json:"roles,omitempty"`
	// Add other claims as needed (e.g. metadata)
}

// Key for Context
type contextKey string

const UserContextKey contextKey = "user"

// NewAuthMiddleware initializes the JWKS fetcher and returns the middleware
func NewAuthMiddleware(ctx context.Context, cfg AuthConfig, logger *slog.Logger) (*AuthMiddleware, error) {
	// Create the JWKS from the URL.
	// This will fetch the keys immediately and cache them.
	// It handles refresh automatically based on Cache-Control headers or errors.
	k, err := keyfunc.NewDefault([]string{cfg.JWKSURL})
	if err != nil {
		return nil, fmt.Errorf("failed to create JWKS from URL %s: %w", cfg.JWKSURL, err)
	}

	return &AuthMiddleware{
		jwks:   k,
		issuer: cfg.Issuer,
		logger: logger,
	}, nil
}

// Handler is the actual middleware function
func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.logger.Warn("Missing Authorization header", "path", r.URL.Path)
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Warn("Invalid Authorization header format", "path", r.URL.Path)
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// 2. Parse and Validate Token
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, m.jwks.Keyfunc)
		if err != nil {
			m.logger.Warn("Token validation failed", "error", err, "path", r.URL.Path)
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// 3. Verify Claims (Issuer)
		if !token.Valid {
			m.logger.Warn("Token is invalid", "path", r.URL.Path)
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			m.logger.Error("Failed to cast claims", "path", r.URL.Path)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Optional: Verify Issuer strictly if configured
		// Note: keyfunc handles signature, but we check business logic claims here
		if m.issuer != "" && claims.Issuer != m.issuer {
			m.logger.Warn("Token issuer mismatch", "expected", m.issuer, "got", claims.Issuer)
			http.Error(w, "Unauthorized: Invalid issuer", http.StatusUnauthorized)
			return
		}

		// 4. Inject into Context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
