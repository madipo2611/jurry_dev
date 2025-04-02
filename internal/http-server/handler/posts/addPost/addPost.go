package addPost

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	resp "jurry_dev/internal/lib/api/response"
	"jurry_dev/internal/lib/logger/sl"
	"jurry_dev/internal/lib/session"
	"jurry_dev/internal/lib/utils"
	_ "jurry_dev/internal/storage"
	"jurry_dev/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Request struct {
	Text  string `json:"text"`
	Image string `json:"image"`
}

type Response struct {
	resp.Response
}

type AddPost interface {
	AddPost(text string, image string, user string) (bool, error)
}

func New(log *slog.Logger, addPost *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		const op = "handler.posts.addPost.New"
		const COOKIE_NAME = "sessionId"

		var req Request
		//распарсить запрос
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		// Удаляем префикс, если он есть (например, data:image/jpeg;base64,)
		if strings.HasPrefix(req.Image, "data:image/jpeg;base64,") {
			req.Image = strings.TrimPrefix(req.Image, "data:image/jpeg;base64,")
		}

		cookie, _ := r.Cookie(COOKIE_NAME)
		sessionId := cookie.Value
		data, err := session.Redis.LoadSession(ctx, sessionId)
		if err != nil {
			log.Error("Ошибка получения данных из сессии:", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		imgBase64, err := base64.StdEncoding.DecodeString(req.Image)
		if err != nil {
			log.Error("failed to decode image base64", op, sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode image"))
			return
		}
		log.Info("userID", slog.Any("serID :", data.UserID))
		randName := utils.GenerateId()
		currentDate := time.Now().Format("2006-01-02")

		upPath := fmt.Sprintf("./uploads/%s/%d", currentDate, data.UserID)

		err = os.MkdirAll(upPath, os.ModePerm)
		if err != nil {
			log.Error("failed to create directory", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to create directory"))
			return
		}

		fileName := fmt.Sprintf("%s/%s.jpg", upPath, randName)

		err = os.WriteFile(fileName, imgBase64, os.ModePerm)
		if err != nil {
			log.Error("failed to write file", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to write file"))
			return
		}

		reg, err := addPost.AddPost(req.Text, fileName, data.UserID)
		if err != nil {
			log.Error("Post not created", sl.Err(err))
			w.WriteHeader(http.StatusBadGateway)
			render.JSON(w, r, resp.Error("Post not created"))
			return
		}
		if reg == false {
			w.WriteHeader(http.StatusBadGateway)
			render.JSON(w, r, resp.Error("Post not created"))
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
