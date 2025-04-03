package getposts

import (
	"github.com/go-chi/render"
	resp "jurry_dev/internal/lib/api/response"
	_ "jurry_dev/internal/storage"
	"jurry_dev/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	resp.Response
	Data []sqlite.Posts `json:"data"`
}

type GetPost interface {
	GetPost(int, int) ([]sqlite.Posts, error)
}

func New(log *slog.Logger, getPost *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handler.posts.addPost.New"
		page, _ := strconv.Atoi(r.URL.Query().Get("page")) // номер страницы
		limit := 5                                         // лимит постов за раз (как в Instagram)
		offset := (page - 1) * limit
		if offset < 0 {
			offset = 0
		}

		data, err := getPost.GetPost(offset, limit)
		if err != nil {
			log.Error("Error get data it is DB:", op, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		responseOK(w, r, data)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, data []sqlite.Posts) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Data:     data,
	})
}
