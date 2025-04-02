package checkauth

import (
	"fmt"
	"github.com/go-chi/render"
	resp "jurry_dev/internal/lib/api/response"
	"jurry_dev/internal/lib/session"
	"log/slog"
	"net/http"
)

type Response struct {
	resp.Response
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const COOKIE_NAME = "sessionId"
		ctx := r.Context()
		cookie, err := r.Cookie(COOKIE_NAME)
		if err != nil {
			fmt.Println("Не удалось получить куку с именем:", COOKIE_NAME)
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("Кука не найдена. Пользователь не авторизован."))
			return
		}
		sessionId := cookie.Value

		data, _ := session.Redis.LoadSession(ctx, sessionId)
		if data == nil {
			fmt.Println("Сессия не найдена:", data)
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("Пользователь не авторизован"))
			return
		}

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
