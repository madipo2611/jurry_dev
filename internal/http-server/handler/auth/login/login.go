package login

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "jurry_dev/internal/lib/api/response"
	"jurry_dev/internal/lib/argon"
	"jurry_dev/internal/lib/logger/sl"
	"jurry_dev/internal/lib/session"
	"jurry_dev/internal/storage"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	resp.Response
}

type Login interface {
	Login(login string) (string, int, error)
}

func New(log *slog.Logger, logins Login) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		session.InitGlobalSession()

		const op = "handler.auth.login.New"
		const COOKIE_NAME = "sessionId"
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

		passDB, userID, err := logins.Login(req.Login)
		if err != nil {
			log.Error("Login not exists", sl.Err(err))
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, resp.Error("Login not exists"))
			return
		}

		pass1 := strings.ReplaceAll(passDB, " ", "")
		hashDec := strings.Split(pass1, ".") // Декодируем строку
		salt, err := hex.DecodeString(hashDec[1])
		if err != nil {
			fmt.Println("salt decode", err)
			fmt.Println("passDB: ", passDB)
			return
		}

		hash, err := hex.DecodeString(hashDec[0])
		if err != nil {
			fmt.Println("hash decode", hash)
			return
		}

		argon := argon.NewArgonHash(2, 32, 20*1024, 64, 4)
		reqPass := []byte(req.Password)
		err = argon.Compare(hash, salt, reqPass)

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
		log.Info("userID", slog.Any(" :", userID))
		sessionId := session.GlobalSession.SetLogin(req.Login, userID)
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
