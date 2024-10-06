package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"

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

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	cognitoClient := cip.NewFromConfig(cfg)

	mux.Post("/register", Make(func(w http.ResponseWriter, r *http.Request) error {
		var req RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return Error{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
		}

		_, err = cognitoClient.SignUp(r.Context(), &cip.SignUpInput{
			ClientId: aws.String(os.Getenv("COGNITO_APP_CLIENT_ID")),
			Username: aws.String(req.Email),
			Password: aws.String(req.Password),
		})
		if err != nil {
			return Error{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
			}
		}

		w.Write([]byte("check your email to verify your email"))

		return err
	}))

	mux.Post("/verification-code", Make(func(w http.ResponseWriter, r *http.Request) error {
		var req VerificationCodeRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return Error{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
		}

		_, err = cognitoClient.ConfirmSignUp(r.Context(), &cip.ConfirmSignUpInput{
			ClientId:         aws.String(os.Getenv("COGNITO_APP_CLIENT_ID")),
			Username:         aws.String(req.Email),
			ConfirmationCode: aws.String(req.Code),
		})
		if err != nil {
			return Error{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
			}
		}

		w.Write([]byte("your account is confirmed"))

		return err
	}))
	http.ListenAndServe(":8080", mux)
}
