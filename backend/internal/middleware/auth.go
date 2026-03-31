package middleware

import (
	"context"
	"net/http"

	"kanban/internal/service"
)

func Auth(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			session, err := authService.ValidateSession(r.Context(), cookie.Value)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", session.UserID)
			ctx = context.WithValue(ctx, "session_id", session.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
