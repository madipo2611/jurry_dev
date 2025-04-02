package logout

import (
	"fmt"
	"github.com/go-chi/render"
	resp "jurry_dev/internal/lib/api/response"
	"jurry_dev/internal/lib/logger/sl"
	"jurry_dev/internal/lib/session"
	"log/slog"
	"net/http"
)

type Response struct {
	resp.Response
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		const COOKIE_NAME = "sessionId"

		cookie, err := r.Cookie(COOKIE_NAME)
		if err != nil {
			fmt.Println("Не удалось получить куку с именем:", COOKIE_NAME)
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("Кука не найдена. Пользователь не авторизован."))
			return
		}
		sessionId := cookie.Value

		err = session.Redis.DeleteSession(ctx, sessionId)
		if err != nil {
			slog.Error("Ошибка удаления сессии:", sl.Err(err))
			w.WriteHeader(http.StatusUnauthorized)
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
