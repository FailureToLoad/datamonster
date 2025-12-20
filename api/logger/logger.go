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

func Debug(ctx context.Context, msg string) {
	slog.DebugContext(ctx, msg)
}

func Info(ctx context.Context, msg string) {
	slog.InfoContext(ctx, msg)
}

func Warn(ctx context.Context, msg string) {
	slog.WarnContext(ctx, msg)
}

func Error(ctx context.Context, msg string) {
	slog.ErrorContext(ctx, msg)
}
