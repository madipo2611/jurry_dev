package login

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "jurry_dev/internal/lib/api/response"
	"jurry_dev/internal/lib/logger/sl"
	"jurry_dev/internal/lib/session"
	"jurry_dev/internal/storage"
	"log/slog"
	"net/http"
	"time"
)

var inMemorySession *session.Session

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	resp.Response
}

type Login interface {
	Login(login string, password string) (bool, error)
}

func New(log *slog.Logger, logins Login) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handler.auth.login.New"
		const COOKIE_NAME = "sessionId"
		inMemorySession = session.NewSession()
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		//распарсить запрос
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			//возвращаем json с ответом ошибки
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		auth, err := logins.Login(req.Login, req.Password)
		if errors.Is(err, storage.ErrLoginNotFound) {
			log.Error("login or password not found", sl.Err(err))
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, resp.Error("invalid username or password"))
			return
		}
		if err != nil {
			log.Error("authorization error", sl.Err(err))
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, resp.Error("authorization error"))
			return
		}
		if auth == false {
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, resp.Error("invalid username or password"))
			return
		}
		sessionId := inMemorySession.SetLogin(req.Login)
		log.Info("sessionId set", slog.String("sessionId", sessionId))
		cookie := &http.Cookie{
			Name:    COOKIE_NAME,
			Domain:  "tailly.ru",
			Value:   sessionId,
			Expires: time.Now().Add(5 * time.Minute),
		}
		http.SetCookie(w, cookie)
		log.Info("cookie set", slog.String("cookie", cookie.Value))
		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
