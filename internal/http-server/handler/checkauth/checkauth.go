package checkauth

import (
	"github.com/go-chi/render"
	"jurry_dev/internal/http-server/middleware/authMiddle"
	resp "jurry_dev/internal/lib/api/response"
	"jurry_dev/internal/storage/sqlite"
	"log/slog"
	"net/http"
)

type User interface {
	GetUser(int) (sqlite.User, error)
}
type Response struct {
	resp.Response
	UserData sqlite.User `json:"userData"`
}

func New(log *slog.Logger, user *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.checkauth.New"

		userID, ok := r.Context().Value(authMiddle.UserIDKey).(int)
		slog.Info("Получаем userID из контекста, checkauth: ", userID)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("User not authenticated"))
			return
		}

		userData, err := user.GetUser(userID)
		if err != nil {
			log.Error("Error getting data from DB:", op, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		log.Debug("Получаем данные userData из бд: ", userData, op)
		responseOK(w, r, userData)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, userData sqlite.User) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		UserData: userData,
	})
}
