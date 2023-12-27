package web

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func DecodeJson(rc io.ReadCloser, data interface{}) error {
	defer rc.Close()
	decoder := json.NewDecoder(rc)
	return decoder.Decode(data)
}

type ctxSettlementIdKey string

const SettlementIdKey ctxSettlementIdKey = "settlementId"

func SettlementIdExtractor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		settlementIdString := chi.URLParam(r, "id")
		if settlementIdString != "" {
			settlementId, convErr := strconv.Atoi(settlementIdString)
			if convErr != nil {
				MakeJsonResponse(w, http.StatusBadRequest, "settlement id should be a number")
				return
			}
			ctx := context.WithValue(r.Context(), SettlementIdKey, settlementId)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
