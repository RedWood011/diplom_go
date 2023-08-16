package logger

import (
	"os"

	"golang.org/x/exp/slog"
)

func InitLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}
