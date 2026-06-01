package main

import (
	"log/slog"
	"os"
	"taskflow/internal/app"
	"taskflow/internal/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("falha ao carregar configuração", "error", err)
		os.Exit(1)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		slog.Error("falha ao inicializar aplicação", "error", err)
		os.Exit(1)
	}

	slog.Info("servidor iniciado", "port", cfg.Port)
	if err := application.Router.Run(":" + cfg.Port); err != nil {
		slog.Error("erro no servidor", "error", err)
		os.Exit(1)
	}
}
