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
			ClientId: aws.String(os.Getenv("COGNITO_CLIENT_ID")),
			Username: aws.String(req.Email),
			Password: aws.String(req.Password),
		})
		if err != nil {
			return Error{
				StatusCode: http.StatusBadRequest,
				Message:    GetAuthProviderError(err),
			}
		}

		w.Write([]byte("check your email to verify your email"))

		return nil 
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
			ClientId:         aws.String(os.Getenv("COGNITO_CLIENT_ID")),
			Username:         aws.String(req.Email),
			ConfirmationCode: aws.String(req.Code),
		})
		if err != nil {
			return Error{
				StatusCode: http.StatusBadRequest,
				Message:    GetAuthProviderError(err),
			}
		}

		w.Write([]byte("your account is confirmed"))

		return nil 
	}))

	mux.Post("/login", Make(func(w http.ResponseWriter, r *http.Request) error {
		var req LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return Error{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
		}

		result, err := cognitoClient.InitiateAuth(r.Context(), &cip.InitiateAuthInput{
			ClientId:       aws.String(os.Getenv("COGNITO_CLIENT_ID")),
			AuthFlow:       "USER_PASSWORD_AUTH",
			AuthParameters: map[string]string{"USERNAME": req.Email, "PASSWORD": req.Password},
		})
		if err != nil {
			return Error{
				StatusCode: http.StatusBadRequest,
				Message:    GetAuthProviderError(err),
			}
		}

		cookie := UserCookie{
			Email:        req.Email,
			AccessToken:  *result.AuthenticationResult.AccessToken,
			RefreshToken: *result.AuthenticationResult.RefreshToken,
		}
		err = WriteCookie(w, "APP_COOKIE", cookie, []byte(os.Getenv("APP_SECRET")))
		w.Write([]byte("you are logged in successfully"))

		return nil
	}))

	mux.With(AuthClient).Get("/client", Make(func(w http.ResponseWriter, r *http.Request) error {
		userCookie, ok := r.Context().Value("UserCookie").(UserCookie)
		if !ok {
			return Error{
				StatusCode: http.StatusUnauthorized,
				Message:    "not valid cookie",
			}
		}	
		
		w.Write([]byte("hello client " + userCookie.Email))

		return nil 
	}))

	http.ListenAndServe(":8080", mux)
}
