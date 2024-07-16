package main

import (
	"3.Server/internal/config"
	"3.Server/internal/lib/logger/sl"
	"3.Server/internal/storage/postgresql"
	"3.Server/pkg/postg"
	"context"
	"fmt"
	"log/slog"
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
	id, err := repository.Create(context.TODO(), "Youtube", "https://www.youtube.com/") // https://www.google.ru/
	if err != nil {
		log.Error("", sl.Err(err))
	}

	fmt.Println(id)

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
