package logger

import (
	"calculator/services/api/pkg/types"
	"log/slog"
	"os"
)

func New(cfg types.Config) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(cfg.LogLevel),
	}))
}
