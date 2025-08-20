package sl

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PrettyHandler struct {
	opts PrettyHandlerOptions
	slog.Handler
	logger *log.Logger
	attrs  []slog.Attr
}

func (opts PrettyHandlerOptions) NewPrettyLogger(out io.Writer) *PrettyHandler {
	return &PrettyHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		logger:  log.New(out, "", 0),
	}
}

func (handler *PrettyHandler) Handle(_ context.Context, record slog.Record) error {
	level := record.Level.String() + ":"

	switch record.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]any, record.NumAttrs())

	record.Attrs(func(attr slog.Attr) bool {
		fields[attr.Key] = attr.Value.Any()

		return true
	})

	for _, attr := range handler.attrs {
		fields[attr.Key] = attr.Value.Any()
	}

	var buffer []byte
	var err error

	if len(fields) > 0 {
		buffer, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := record.Time.Format("[15:05:05.000]")
	msg := color.CyanString(record.Message)

	handler.logger.Println(
		timeStr,
		level,
		msg,
		color.WhiteString(string(buffer)),
	)

	return nil
}

func (handler *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: handler.Handler,
		logger:  handler.logger,
		attrs:   attrs,
	}
}

func (handler *PrettyHandler) WithGroup(name string) slog.Handler {
	// TODO: implement
	return &PrettyHandler{
		Handler: handler.Handler.WithGroup(name),
		logger:  handler.logger,
	}
}
