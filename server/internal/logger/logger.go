package logger

import (
	"log"
	"os"

	"golang.org/x/exp/slog"
)

func InitLogger() *slog.Logger {
	var logger *slog.Logger
	file, err := os.OpenFile("/Users/evyaroshen/GolandProjects/GophKeeper/diplom/server/loggers.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger = slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}))

	return logger
}
