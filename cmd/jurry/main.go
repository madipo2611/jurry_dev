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
	"jurry_dev/internal/http-server/ha
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
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(corsMiddleware)

	router.Handle("/uploads/", http.StripPrefix("/uploads/",
		http.FileServer(http.Dir("./uploads"))))

	router.Post("/api/login", login.New(log, storage))
	router.Post("/api/register", register.New(log, storage))
	router.Post("/api/post", addPost.New(log, storage))
	router.Get("/api/checkauth", checkauth.New(log))
	router.Post("/api/logout", logout.New(log))
	router.Get("/api/posts", getposts.New(log, storage))
	router.Delete("/api/delpost", delpost.New(log, storage))

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
		w.Header().Set("Access-Control-Allow-Origin", "https://tailly.ru") // Замените на ваш фронтенд
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
