package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func writeJSON(rw http.ResponseWriter, status int, data any) {
	js, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		slog.Error("could not marshal json response", slog.Any("error", jsonErr), slog.Int("targetStatus", status))
	}

	rw.WriteHeader(status)
	rw.Header().Set("Content-Type", "application/json")
	_, writeErr := rw.Write(js)
	if writeErr != nil {
		slog.Error("could not write response", slog.Any("error", writeErr), slog.Int("targetStatus", status))
	}
}

func BadRequest(rw http.ResponseWriter, reason string, err error) {
	slog.Error("bad request", slog.Any("error", err))
	writeJSON(rw, http.StatusBadRequest, reason)
}

func InternalServerError(rw http.ResponseWriter, reason string, err error) {
	slog.Error("internal server error", slog.Any("error", err))
	writeJSON(rw, http.StatusInternalServerError, reason)
}

func Unauthorized(rw http.ResponseWriter, err error) {
	slog.Error("unauthorized", slog.Any("error", err))
	writeJSON(rw, http.StatusUnauthorized, "Unauthorized")
}

func NoContent(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNoContent)
}

func OK(rw http.ResponseWriter, data any) {
	writeJSON(rw, http.StatusOK, data)
}
