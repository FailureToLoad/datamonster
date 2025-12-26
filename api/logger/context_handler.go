package logger

import (
	"context"
	"log/slog"

	"github.com/failuretoload/datamonster/request"
	"github.com/gofrs/uuid/v5"
)

type ContextHandler struct {
	next slog.Handler
}

func NewContextHandler(next slog.Handler) slog.Handler {
	if next == nil {
		panic("next handler cannot be nil")
	}
	return &ContextHandler{next: next}
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *ContextHandler) Handle(ctx context.Context, record slog.Record) error {
	correlationID := request.CorrelationID(ctx)
	if correlationID != "" {
		record.AddAttrs(slog.String("correlationID", correlationID))
	}

	userID := request.UserID(ctx)
	if userID != "" {
		record.AddAttrs(slog.String("userID", userID))
	}

	settlementID := request.SettlementID(ctx)
	if settlementID != uuid.Nil {
		record.AddAttrs(slog.String("settlementID", settlementID.String()))
	}

	return h.next.Handle(ctx, record)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{next: h.next.WithAttrs(attrs)}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{next: h.next.WithGroup(name)}
}
