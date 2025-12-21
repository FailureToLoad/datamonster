package logger

import (
	"context"
	"log/slog"
	"os"
)

func Setup() {
	base := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(NewContextHandler(base))
	slog.SetDefault(logger)
}

func attrsToAny(fields []slog.Attr) []any {
	if len(fields) == 0 {
		return nil
	}
	args := make([]any, len(fields))
	for i := range fields {
		args[i] = fields[i]
	}
	return args
}

func Debug(ctx context.Context, msg string, fields ...slog.Attr) {
	slog.DebugContext(ctx, msg, attrsToAny(fields)...)
}

func Info(ctx context.Context, msg string, fields ...slog.Attr) {
	slog.InfoContext(ctx, msg, attrsToAny(fields)...)
}

func Warn(ctx context.Context, msg string, fields ...slog.Attr) {
	slog.WarnContext(ctx, msg, attrsToAny(fields)...)
}

func Error(ctx context.Context, msg string, fields ...slog.Attr) {
	slog.ErrorContext(ctx, msg, attrsToAny(fields)...)
}
