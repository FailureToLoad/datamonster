package logger

import "log/slog"

func SettlementID(id string) slog.Attr {
	return slog.String("settlement_id", id)
}

func ErrorField(e error) slog.Attr {
	return slog.Any("error", e)
}
