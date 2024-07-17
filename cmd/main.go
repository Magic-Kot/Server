package main

import (
	"3.Server/internal/config"
	"3.Server/internal/http-server/handlers/redirect"
	"3.Server/internal/http-server/handlers/url/save"
	"3.Server/internal/lib/logger/sl"
	"3.Server/internal/storage/postgresql"
	"3.Server/pkg/postg"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	//fmt.Println(cfg)

	log := setupLogger(cfg.Env)

	log.Info("starting server", slog.String("env", cfg.Env))
	log.Debug("debug message are enabled")

	// postgres://postgres:12345@localhost:5438
	pool, err := postg.NewClient(context.TODO(), 5, "postgres", "12345", "localhost", "5438", "postgres")
	if err != nil {
		log.Error("failed to init storage:", sl.Err(err))
		os.Exit(1)
	}

	//fmt.Println(pool)

	repository := postgresql.NewRepository(pool, log)

	/*id, err := repository.Create(context.TODO(), "Youtube", "https://www.youtube.com/") // https://www.google.ru/
	if err != nil {
		log.Error("", sl.Err(err))
	}

	fmt.Println(id)*/

	/*dbURL, err := repository.Get(context.Background(), "Yandex")
	if err != nil {
		log.Error("", sl.Err(err))
	}

	fmt.Println(dbURL)*/

	/*aliases, err := repository.GetAll(context.TODO())
	if err != nil {
		log.Error("", sl.Err(err))
	}

	fmt.Println(aliases)*/

	/*err = repository.Delete(context.TODO(), "Youtube")
	if err != nil {
		log.Error("repository.Delete:", sl.Err(err))
	}*/

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger) // будет логировать все входящие запросы (добавит строчку)
	//router.Use(mwLogger.New(log)) // самописный middleware
	router.Use(middleware.Recoverer) // при возникновении паники будет внутри хендрера, будет востановлена паника
	router.Use(middleware.URLFormat) // форматирует url

	// authorization
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
			//cfg.HTTPServer.User: cfg.HTTPServer.Password, // добавление остальных пользователей и паролей для них
		}))

		r.Post("/", save.New(log, repository))
		// TODO: add DELETE /url/{id}
	})

	// handlers
	router.Post("/url", save.New(log, repository))
	router.Get("/{alias}", redirect.New(log, repository))
	//router.Delete("/{alias}", redirect.New(log, repository))

	log.Info("starting server", slog.String("address", cfg.Address))

	//done := make(chan os.Signal, 1)
	//signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
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
