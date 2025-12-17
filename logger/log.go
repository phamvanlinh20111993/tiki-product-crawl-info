package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/util"
	"strings"
	"sync"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

var (
	logger *slog.Logger
	once   sync.Once
)

const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace:      "TRACE", // -8
	slog.LevelDebug: "DEBUG", // -4
	slog.LevelInfo:  "INFO",  // 0
	slog.LevelWarn:  "WARN",  // 4
	slog.LevelError: "ERROR", //8
	LevelFatal:      "FATAL", // 12
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
		slogOpts := &slog.HandlerOptions{
			Level:     Levels[strings.ToUpper(loggerConfig.Level)],
			AddSource: loggerConfig.IsAddSource,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					a.Value = slog.StringValue(
						util.TimeToString(a.Value.Time(), util.Format_yyyy_mm_dd_space_hh_dot_mm_dot_ss),
					)
				}
				if a.Key == slog.LevelKey {
					level := a.Value.Any().(slog.Level)
					levelLabel, exists := LevelNames[level]
					if !exists {
						levelLabel = level.String()
					}

					a.Value = slog.StringValue(levelLabel)
				}
				if a.Key == slog.MessageKey {
					a.Key = "message"
				}
				return a
			},
		}

		//opts := PrettyHandlerOptions{
		//	SlogOpts: *slogOpts,
		//}
		//handler := NewPrettyHandler(os.Stdout, opts)
		handler := slog.NewTextHandler(os.Stdout, slogOpts)
		logger = slog.New(handler)
	})
	return logger
}

func logCommon(logLevel slog.Level, msg string, args ...any) {
	var msgFormat strings.Builder
	var attrs []slog.Attr

	for _, v := range args {
		attrValue, isSlogAttr := v.(slog.Attr)
		if isSlogAttr {
			attrs = append(attrs, attrValue)
			continue
		}

		errorValue, isError := v.(error)
		if isError {
			msgFormat.WriteString(" ")
			msgFormat.WriteString(errorValue.Error())
			continue
		}

		strValue, isStr := v.(string)
		if isStr {
			msgFormat.WriteString(" ")
			msgFormat.WriteString(strValue)
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
	var logLevel = Levels[strings.ToUpper(configuration.GetLoggerConfig().Level)]
	if logLevel <= slog.LevelInfo {
		logCommon(slog.LevelInfo, msg, args...)
	}
}

func logError(msg string, args ...any) {
	var logLevel = Levels[strings.ToUpper(configuration.GetLoggerConfig().Level)]
	if logLevel <= slog.LevelError {
		logCommon(slog.LevelError, msg, args...)
	}
}

func logDebug(msg string, args ...any) {
	var logLevel = Levels[strings.ToUpper(configuration.GetLoggerConfig().Level)]
	if logLevel <= slog.LevelDebug {
		logCommon(slog.LevelDebug, msg, args...)
	}
}

func logWarn(msg string, args ...any) {
	var logLevel = Levels[strings.ToUpper(configuration.GetLoggerConfig().Level)]
	if logLevel <= slog.LevelWarn {
		logCommon(slog.LevelWarn, msg, args...)
	}
}

var LogInfo = logInfo
var LogError = logError
var LogDebug = logDebug
var LogWarn = logWarn
