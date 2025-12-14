package util

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"selfstudy/crawl/product/configuration"
	"strings"
	"sync"
)

var (
	logger *slog.Logger
	once   sync.Once
)

const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
}

var Levels = map[string]slog.Level{
	"TRACE": LevelTrace,
	"FATAL": LevelFatal,
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"ERROR": slog.LevelError,
	"WARN":  slog.LevelWarn,
}

func getLoggerInstance() *slog.Logger {
	once.Do(func() {
		var loggerConfig = configuration.GetLoggerConfig()
		opts := &slog.HandlerOptions{
			Level:     Levels[strings.ToUpper(loggerConfig.Level)],
			AddSource: loggerConfig.IsAddSource,
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

func logCommon(logLevel slog.Level, msg string, val ...any) {
	var msgFormat strings.Builder
	var attrs []slog.Attr
	for _, v := range val {
		attrValue, isSlogAttr := v.(slog.Attr)
		if !isSlogAttr {
			attrs = append(attrs, attrValue)
			continue
		}

		errorValue, isError := v.(error)
		if isError {
			msgFormat.WriteString(" ")
			msgFormat.WriteString(errorValue.Error())
			continue
		}

		msgFormat.WriteString(" ")
		msgFormat.WriteString(fmt.Sprintf("%v", v))
	}

	getLoggerInstance().LogAttrs(
		context.Background(),
		logLevel,
		msg+msgFormat.String(),
		attrs...)
}

func logInfo(msg string, args ...any) {
	logCommon(slog.LevelInfo, msg, args)
}

func logError(msg string, args ...any) {
	logCommon(slog.LevelError, msg, args)
}

func logDebug(msg string, args ...any) {
	logCommon(slog.LevelDebug, msg, args)
}

func logWarn(msg string, args ...any) {
	logCommon(slog.LevelWarn, msg, args)
}

var LogInfo = logInfo
var LogError = logError
var LogDebug = logDebug
var LogWarn = logWarn
