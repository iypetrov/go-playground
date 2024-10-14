package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

var UUIDKey = "UUID_KEY"

func UUIDFormat(next http.Handler) http.Handler {
	var val uuid.UUID
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if len(id) != 0 {
			i, err := uuid.Parse(id)
			if err != nil {
				err = writeErr(w, Error{
					StatusCode: http.StatusBadRequest,
					Message:    err.Error(),
				})
				if err != nil {
					return
				}
			}
			val = i
		}
		ctx := context.WithValue(r.Context(), UUIDKey, val)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeErr(w http.ResponseWriter, err Error) error {
	w.WriteHeader(err.StatusCode)
	return json.NewEncoder(w).Encode(err)
}
