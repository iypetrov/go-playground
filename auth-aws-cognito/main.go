package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	mux.Get("/", Make(func(w http.ResponseWriter, r *http.Request) error {
		_, err := w.Write([]byte("Hello World!"))
		return err
	}))

	http.ListenAndServe(":8080", mux)
}
