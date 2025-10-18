package middleware

import (
	"context"
	users "dijam-ecommerce/internal/user"
	"dijam-ecommerce/pkg/apierrors"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type contextKey string

const (
	// UserIDKey is the key for user ID in the request context
	UserIDKey contextKey = "userID"
)

func AuthMiddleware(us *users.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				apierrors.Unauthorized(w, r, errors.New("invalid authorization format"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				apierrors.Unauthorized(w, r, errors.New("invalid authorization format"))
				return
			}
			tokenString := parts[1]
			claims, err := us.ValidateToken(tokenString)
			if err != nil {
				apierrors.Unauthorized(w, r, errors.New("invalid or expired token"))
			}

			userIDStr, ok := claims["sub"].(string)
			if !ok {
				apierrors.Unauthorized(w, r, errors.New("invalid token claims"))
				return
			}
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				apierrors.Unauthorized(w, r, errors.New("invalid user ID in token"))
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	return userID, ok
}
