package middleware

import (
	"context"
	"e-commerce/utils"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserIDKey  contextKey = "UserID"
	IsAdminKey contextKey = "is_admin"
)

// AuthMiddleWare authenticates the user and sets the user ID and admin status in the context
func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		userID, isAdmin, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		// Set user ID and admin status in the request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, IsAdminKey, isAdmin)

		next.ServeHTTP(w, r.WithContext(ctx)) // Pass context along with the request
	})
}

// AdminMiddleware ensures the user is an admin before proceeding
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(UserIDKey).(int)
		isAdmin, adminOk := r.Context().Value(IsAdminKey).(bool)

		// Ensure the user is authenticated and an admin
		if !ok || !adminOk || !isAdmin {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		// Proceed with the next handler if the user is an admin
		next.ServeHTTP(w, r)
	})
}
