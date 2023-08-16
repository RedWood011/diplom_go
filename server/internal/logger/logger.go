package logger

import (
	"os"

	"golang.org/x/exp/slog"
)

func InitLogger() *slog.Logger {
	var logger *slog.Logger

	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	return logger
}
