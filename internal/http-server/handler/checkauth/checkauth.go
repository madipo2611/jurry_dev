package checkauth

import (
	"github.com/go-chi/render"
	"jurry_dev/internal/http-server/middleware/authMiddle"
	resp "jurry_dev/internal/lib/api/response"
	"log/slog"
	"net/http"
)

type Response struct {
	resp.Response
	UserID int `json:"userID"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(authMiddle.UserIDKey).(int)
		slog.Info("Получаем userID из контекста, checkauth: ", userID)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("User not authenticated"))
			return
		}
		responseOK(w, r, userID)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, userID int) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		UserID:   userID,
	})
}
