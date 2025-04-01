package checkauth

import (
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

		cookie, err := r.Cookie(COOKIE_NAME)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("Пользователь не авторизован"))
			return
		}
		_, userID := session.GlobalSession.Get(cookie.Value)
		if userID == 0 {
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
