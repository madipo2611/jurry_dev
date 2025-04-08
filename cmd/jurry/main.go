package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"jurry_dev/internal/config"
	"jurry_dev/internal/http-server/handler/auth/login"
	"jurry_dev/internal/http-server/handler/auth/register"
	"jurry_dev/internal/http-server/handler/checkauth"
	"jurry_dev/internal/http-server/handler/logout"
	"jurry_dev/internal/http-server/handler/posts/addPost"
	"jurry_dev/internal/http-server/handler/posts/delpost"
	"jurry_dev/internal/http-server/handler/posts/getposts"
	"jurry_dev/internal/http-server/handler/user"
	"jurry_dev/internal/http-server/middleware/authMiddle"
	"jurry_dev/internal/lib/logger/sl"
	"jurry_dev/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("starting server", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.PS)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
	}
	_ = storage

	router := chi.NewRouter()
	// Глобальные middleware (применяются ко ВСЕМ роутам)
	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.URLFormat,
		corsMiddleware,
	)

	router.Post("/login", login.New(log, storage))
	router.Post("/register", register.New(log, storage))
	router.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	router.Group(func(r chi.Router) {
		r.Use(authMiddle.AuthMiddleware) // Middleware проверки авторизации

		r.Post("/post", addPost.New(log, storage))
		r.Get("/checkauth", checkauth.New(log, storage))
		r.Post("/logout", logout.New(log))
		r.Get("/posts", getposts.New(log, storage))
		r.Delete("/delpost", delpost.New(log, storage))
		r.Get("/user", user.MeHandler(log, storage))
	})

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	//TODO: run server

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}
	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Замените на ваш фронтенд
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Обработка preflight запросов
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
