package register

import (
	"encoding/hex"
	"errors"
	"github.com/go-chi/render"
	resp "jurry_dev/internal/lib/api/response"
	"jurry_dev/internal/lib/argon"
	"jurry_dev/internal/lib/logger/sl"
	"jurry_dev/internal/storage"
	"log/slog"
	"net/http"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
}

type Response struct {
	resp.Response
}

type Register interface {
	Register(login string, password string, gender string) (bool, error)
}

func New(log *slog.Logger, register Register) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handler.auth.register.New"

		var req Request
		//распарсить запрос
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			//возвращаем json с ответом ошибки
			w.WriteHeader(http.StatusSeeOther)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		pass := []byte(req.Password)

		argon := argon.NewArgonHash(2, 32, 20*1024, 64, 4)
		hashSalt, err := argon.GenerateHash(pass, nil)
		if err != nil {
			log.Error("failed to generate hash salt", sl.Err(err))
			w.WriteHeader(http.StatusSeeOther)
			render.JSON(w, r, resp.Error("failed to generate hash salt"))
			return
		}
		hashHex := hex.EncodeToString(hashSalt.Hash)
		saltHex := hex.EncodeToString(hashSalt.Salt)

		password := hashHex + "." + saltHex
		log.Info("hash", slog.Any("hash", hashSalt.Hash))
		log.Info("salt", slog.Any("salt", hashSalt.Salt))
		log.Info("password", slog.Any("password", password))

		reg, err := register.Register(req.Login, password, req.Gender)
		if errors.Is(err, storage.ErrLoginNotFound) {
			log.Error("login or password not found", sl.Err(err))
			w.WriteHeader(http.StatusSeeOther)
			render.JSON(w, r, resp.Error("invalid username or password"))
			return
		}
		if err != nil {
			log.Error("Login exists", sl.Err(err))
			w.WriteHeader(http.StatusSeeOther)
			render.JSON(w, r, resp.Error("Login exists"))
			return
		}
		if reg == false {
			w.WriteHeader(http.StatusSeeOther)
			render.JSON(w, r, resp.Error("Account already exists"))
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
