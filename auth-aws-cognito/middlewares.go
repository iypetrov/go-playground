package main

import (
	"context"
	"os"
	"net/http"
)

func AuthClient(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userCookie, err := ReadCookie(r, "APP_COOKIE", []byte(os.Getenv("APP_SECRET")))
		if err != nil {
			return
		}
		ctx := context.WithValue(r.Context(), "UserCookie", userCookie)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}