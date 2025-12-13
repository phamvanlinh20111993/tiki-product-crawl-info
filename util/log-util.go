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

var LevelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
}

func getLoggerInstance() *slog.Logger {
	once.Do(func() {
		opts := &slog.HandlerOptions{
			Level:     slog.LevelDebug, // check config
			AddSource: true,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					level := a.Value.Any().(slog.Level)
					levelLabel, exists := LevelNames[level]
					if !exists {
						levelLabel = level.String()
					}

					a.Value = slog.StringValue(levelLabel)
				}

				return a
			},
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

var LogInfo = logInfo
var LogError = logError
var LogDebug = logDebug
