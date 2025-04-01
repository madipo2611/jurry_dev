package logout

import (
	"github.com/go-chi/render"
	resp "jurry_dev/internal/lib/api/response"
	"log/slog"
	"net/http"
)

type Response struct {
	resp.Response
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const COOKIE_NAME = "sessionId"

		http.SetCookie(w, &http.Cookie{
			Name:   COOKIE_NAME,
			Value:  "",
			MaxAge: -1,
		})

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
