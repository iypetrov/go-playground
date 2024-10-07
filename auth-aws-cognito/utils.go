package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Error struct {
	StatusCode int `json:"statusCode"`
	Message    any `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("API error: %d", e.StatusCode)
}

func Make(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			var apiErr Error
			if errors.As(err, &apiErr) {
				err := WriteJson(w, apiErr.StatusCode, apiErr)
				if err != nil {
					return
				}
			}
		}
	}
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func GetAuthProviderError(authErr error) string {
	parts := strings.Split(authErr.Error(), ",")

	lastPart := parts[len(parts)-1]
	finalParts := strings.Split(lastPart, ":")

	return strings.TrimSpace(finalParts[len(finalParts)-1])
}
