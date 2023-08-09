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

/*func MiddlewareLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		if err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("method = %s, patch = %s,  IP = %s, statusCode = %d, timeResponse=%v", c.Method(), c.Path(), c.IP(), c.Response().StatusCode(), time.Since(start)))
		return nil
	}
}*/
