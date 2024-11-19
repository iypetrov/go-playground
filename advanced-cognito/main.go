package main

import (
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
)

func Render(w http.ResponseWriter, name string, data interface{}) {
	template.Must(template.ParseGlob("templates/*.html")).ExecuteTemplate(w, name, data)
}

func main() {
	mux := chi.NewRouter()
	mux.Route("/p", func(r chi.Router) {
		r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
			Render(w, "register", nil)
		})
		r.Get("/confirm", func(w http.ResponseWriter, r *http.Request) {
			Render(w, "confirm", nil)
		})
		r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			Render(w, "login", nil)
		})
		r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
			Render(w, "home", nil)
		})
	})
	mux.Route("/api/v0", func(r chi.Router) {
		mux.Get("/register", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello world!"))
		})
	})
	http.ListenAndServe(":8080", mux)
}
