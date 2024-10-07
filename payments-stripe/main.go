package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	mux.Get("/payments", Make(func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("uhuu, your payments"))

		return nil
	}))

	http.ListenAndServe(":8080", mux)
}
