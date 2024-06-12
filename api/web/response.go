package web

import (
	"encoding/json"
	"net/http"
)

func MakeJsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(data)
		_, _ = w.Write(body)
	}
}
