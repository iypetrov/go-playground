package main

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(cfg)

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	mux.Get("/index", Make(func(w http.ResponseWriter, r *http.Request) error {
		return Render(w, "index", nil)
	}))

	mux.Route("/", func(r chi.Router) {
		mux.Post("/upload", Make(func(w http.ResponseWriter, r *http.Request) error {
			r.ParseMultipartForm(10 << 20) // 10MB max size for the image

			name := r.FormValue("name")
			description := r.FormValue("description")

			file, _, err := r.FormFile("image")
			if err != nil {
				return Error{
					StatusCode: http.StatusBadRequest,
					Message:    err.Error(),
				}
			}
			defer file.Close()

			id := uuid.NewString()
			_, err = s3Client.PutObject(r.Context(), &s3.PutObjectInput{
				Bucket: aws.String(os.Getenv("BUCKET_NAME")),
				Key:    aws.String(id + ".png"),
				Body:   file,
			})
			if err != nil {
				return Error{
					StatusCode: http.StatusInternalServerError,
					Message:    err.Error(),
				}
			}

			result := Product{
				Name:        name,
				Description: description,
				ImageURL:    "http://localhost:8080/image/" + id,
			}
			return Render(w, "result", result)
		}))

		mux.With(UUIDFormat).Get("/image/{id}", Make(func(w http.ResponseWriter, r *http.Request) error {
			id, ok := r.Context().Value(UUIDKey).(uuid.UUID)
			if !ok {
				return Error{
					StatusCode: http.StatusBadRequest,
					Message:    "failed to load uuid",
				}
			}

			resp, err := s3Client.GetObject(context.Background(), &s3.GetObjectInput{
				Bucket: aws.String(os.Getenv("BUCKET_NAME")),
				Key:    aws.String(id.String() + ".png"),
			})
			if err != nil {
				return Error{
					StatusCode: http.StatusInternalServerError,
					Message:    err.Error(),
				}
			}
			defer resp.Body.Close()

			w.Header().Set("Content-Type", "image/png")

			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return Error{
					StatusCode: http.StatusInternalServerError,
					Message:    err.Error(),
				}
			}

			return nil
		}))

	})

	http.ListenAndServe(":8080", mux)
}
