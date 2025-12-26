package response

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/failuretoload/datamonster/logger"
	"github.com/failuretoload/datamonster/request"
)

func BadRequest(ctx context.Context, rw http.ResponseWriter, err error) {
	logger.Error(ctx, "bad request", slog.Any("error", err))
	writeError(ctx, rw, http.StatusBadRequest)
}

func InternalServerError(ctx context.Context, rw http.ResponseWriter, err error) {
	slog.Error("internal server error", slog.Any("error", err))
	writeError(ctx, rw, http.StatusInternalServerError)
}

func Unauthorized(ctx context.Context, rw http.ResponseWriter, err error) {
	slog.Error("unauthorized", slog.Any("error", err))
	writeError(ctx, rw, http.StatusUnauthorized)
}

func NotFound(ctx context.Context, rw http.ResponseWriter, err error) {
	slog.Error("not found", slog.Any("error", err))
	writeError(ctx, rw, http.StatusNotFound)
}

func NoContent(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNoContent)
}

func OK(ctx context.Context, rw http.ResponseWriter, data any) {
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	js, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		InternalServerError(ctx, rw, jsonErr)
		return
	}

	_, writeErr := rw.Write(js)
	if writeErr != nil {
		InternalServerError(ctx, rw, writeErr)
		return
	}
}

type errorResponse struct {
	Status        string `json:"status"`
	CorrelationID string `json:"correlationId"`
}

func writeError(ctx context.Context, rw http.ResponseWriter, status int) {
	statusString := http.StatusText(status)
	er := errorResponse{
		Status: statusString,
	}

	if cid := request.CorrelationID(ctx); cid != "" {
		er.CorrelationID = cid
	}

	rw.WriteHeader(status)
	js, jsonErr := json.Marshal(er)
	if jsonErr != nil {
		logger.Error(ctx, fmt.Sprintf("could not marshal %s response: %v", statusString, jsonErr))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, writeErr := rw.Write(js)
	if writeErr != nil {
		logger.Error(ctx, fmt.Sprintf("could not write %s response: %v", statusString, writeErr))
	}
}
