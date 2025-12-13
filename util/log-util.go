package util

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

var (
	logger *slog.Logger
	once   sync.Once
)

func getLoggerInstance() *slog.Logger {
	once.Do(func() {
		opts := &slog.HandlerOptions{
			Level:     slog.LevelDebug, // check config
			AddSource: true,
		}
		handler := slog.NewJSONHandler(os.Stdout, opts)
		logger = slog.New(handler)
	})
	return logger
}

func logInfo(msg string) {
	getLoggerInstance().LogAttrs(
		context.Background(),
		slog.LevelInfo,
		msg)
}

func logError(msg string) {
	getLoggerInstance().LogAttrs(
		context.Background(),
		slog.LevelError,
		msg)
}

func logDebug(msg string) {
	getLoggerInstance().LogAttrs(
		context.Background(),
		slog.LevelDebug,
		msg)
}
