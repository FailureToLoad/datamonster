package request

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
)

const maxRequestBodySize = 1 << 20 // 1 MB

func DecodeJSON(rc io.ReadCloser, data any) error {
	defer func(rc io.ReadCloser) {
		err := rc.Close()
		if err != nil {
			slog.Log(context.Background(), slog.LevelWarn, "error closing ReadCloser", slog.Any("error", err))
		}
	}(rc)

	limitedReader := io.LimitReader(rc, maxRequestBodySize)
	decoder := json.NewDecoder(limitedReader)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(data); err != nil {
		return err
	}

	if decoder.More() {
		return errors.New("request body exceeds maximum size")
	}

	return nil
}
