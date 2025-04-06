package authMiddle

import (
	"context"
	"github.com/go-chi/render"
	resp "jurry_dev/internal/lib/api/response"
	"jurry_dev/internal/lib/session"
	"log/slog"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "UserID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const COOKIE_NAME = "sessionId"
		ctx := r.Context()

		cookie, err := r.Cookie(COOKIE_NAME)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("User not authenticated"))
			return
		}

		sessionId := cookie.Value

		data, _ := session.Redis.LoadSession(ctx, sessionId)
		slog.Info("Получаем data.UserID и отправляем в контекст:", data.UserID)
		if data != nil {
			ctx = context.WithValue(ctx, UserIDKey, data.UserID)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
