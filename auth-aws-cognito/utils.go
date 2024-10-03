package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
