package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/go-playground/form"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
)

type Error struct {
	StatusCode int
	Message    any
}

func (e Error) Error() string {
	return fmt.Sprintf("API error: %d", e.StatusCode)
}

func AWSError(authErr error) string {
	parts := strings.Split(authErr.Error(), ",")

	lastPart := parts[len(parts)-1]
	finalParts := strings.Split(lastPart, ":")

	return strings.TrimSpace(finalParts[len(finalParts)-1])
}

func Make(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			var apiErr Error
			if errors.As(err, &apiErr) {
				w.WriteHeader(apiErr.StatusCode)
				log.Printf("error %d: %s", apiErr.StatusCode, apiErr.Message)
			}
		}
	}
}

func Render(w http.ResponseWriter, name string, data interface{}) {
	template.Must(template.ParseGlob("templates/*.html")).ExecuteTemplate(w, name, data)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	cognitoClient := cip.NewFromConfig(cfg)

	decoder := form.NewDecoder()

	mux := chi.NewRouter()
	mux.Route("/p", func(mux chi.Router) {
		mux.Get("/register", func(w http.ResponseWriter, r *http.Request) {
			Render(w, "register", nil)
		})
		mux.Get("/confirm", func(w http.ResponseWriter, r *http.Request) {
			Render(w, "confirm", nil)
		})
		mux.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			Render(w, "login", nil)
		})
		mux.Get("/home", func(w http.ResponseWriter, r *http.Request) {
			Render(w, "home", nil)
		})
	})
	mux.Route("/api/v0", func(mux chi.Router) {
		mux.Post("/register", Make(func(w http.ResponseWriter, r *http.Request) error {
			err := r.ParseForm()
			if err != nil {
				return err
			}

			var req RegisterRequest
			err = decoder.Decode(&req, r.Form)
			if err != nil {
				log.Panic(err)
			}

			output, err := cognitoClient.SignUp(r.Context(), &cip.SignUpInput{
				ClientId: aws.String(os.Getenv("AWS_COGNITO_CLIENT_ID")),
				Username: aws.String(req.Email),
				Password: aws.String(req.Password),
			})
			if err != nil {
				return Error{
					StatusCode: http.StatusBadRequest,
					Message:    AWSError(err),
				}
			}
			log.Println(output)

			w.Write([]byte("check your email to verify your email"))
			return nil
		}))
	})
	http.ListenAndServe(":8080", mux)
}
