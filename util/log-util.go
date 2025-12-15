package util

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"selfstudy/crawl/product/configuration"
	"strings"
	"sync"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"
	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	//b, err := json.MarshalIndent(fields, "", "  ")
	//if err != nil {
	//	return err
	//}
	msg := color.CyanString(r.Message)

	// h.l.Println(timeToString(r.Time, Format_yyyy_mm_dd_space_hh_dot_mm_dot_ss), level, msg, color.WhiteString(string(b)))

	h.l.Println(timeToString(r.Time, Format_yyyy_mm_dd_space_hh_dot_mm_dot_ss), level, msg)
	return nil
}

func NewPrettyHandler(
	out io.Writer,
	opts PrettyHandlerOptions,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewTextHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 1),
	}

	return h
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
	LevelTrace:      "TRACE",
	LevelFatal:      "FATAL",
	slog.LevelDebug: "DEBUG",
	slog.LevelInfo:  "INFO",
	slog.LevelError: "ERROR",
	slog.LevelWarn:  "WARN",
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
						timeToString(a.Value.Time(), Format_yyyy_mm_dd_space_hh_dot_mm_dot_ss),
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
	logCommon(slog.LevelInfo, msg, args...)
}

func logError(msg string, args ...any) {
	logCommon(slog.LevelError, msg, args...)
}

func logDebug(msg string, args ...any) {
	logCommon(slog.LevelDebug, msg, args...)
}

func logWarn(msg string, args ...any) {
	logCommon(slog.LevelWarn, msg, args...)
}

var LogInfo = logInfo
var LogError = logError
var LogDebug = logDebug
var LogWarn = logWarn
