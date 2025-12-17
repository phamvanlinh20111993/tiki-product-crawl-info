package logger

import (
	"context"
	"io"
	"log"
	"log/slog"
	"selfstudy/crawl/product/util"

	"github.com/fatih/color"
)

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

	h.l.Println(util.TimeToString(r.Time, util.Format_yyyy_mm_dd_space_hh_dot_mm_dot_ss), level, msg)
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
