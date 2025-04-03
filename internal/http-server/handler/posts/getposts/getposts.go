package getposts

import (
	"github.com/go-chi/render"
	resp "jurry_dev/internal/lib/api/response"
	"jurry_dev/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	resp.Response
	Data       []sqlite.Posts `json:"data"`
	TotalCount int            `json:"total_count"`
	HasMore    bool           `json:"has_more"`
}

type GetPost interface {
	GetPost(int, int) ([]sqlite.Posts, int, error)
}

func New(log *slog.Logger, getPost *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.posts.addPost.New"

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}

		limit := 5
		offset := (page - 1) * limit

		data, totalCount, err := getPost.GetPost(offset, limit)
		if err != nil {
			log.Error("Error getting data from DB:", op, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		hasMore := offset+limit < totalCount

		render.JSON(w, r, Response{
			Response:   resp.OK(),
			Data:       data,
			TotalCount: totalCount,
			HasMore:    hasMore,
		})
	}
}
