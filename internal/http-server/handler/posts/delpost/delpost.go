package delpost

import (
	"github.com/go-chi/render"
	resp "jurry_dev/internal/lib/api/response"
	"jurry_dev/internal/lib/logger/sl"
	"jurry_dev/internal/storage/sqlite"
	"log/slog"
	"net/http"
)

type Request struct {
	Id int `json:"id"`
}

type Response struct {
	resp.Response
}

type DelPost interface {
	DelPost(int) error
}

func New(log *slog.Logger, delPost *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.posts.addPost.New"

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

		err = delPost.DelPost(req.Id)
		if err != nil {
			log.Error("Error getting data from DB:", op, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
