package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/paymentintent"
	"github.com/stripe/stripe-go/v80/webhook"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	mux.Route("/", func(r chi.Router) {
		r.Handle("/web/*", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

		r.Get("/checkout", Make(func(w http.ResponseWriter, r *http.Request) error {
			product := Product{
				Name:        "stripe t-shirt",
				Description: "nice t-shirt",
				Price:       19.99,
			}
			return Render(w, "checkout", product)
		}))

		r.Get("/checkout/result", Make(func(w http.ResponseWriter, r *http.Request) error {
			return Render(w, "result", Product{})
		}))

		r.Route("/payments", func(r chi.Router) {
			r.Get("/config", Make(func(w http.ResponseWriter, r *http.Request) error {
				WriteJson(w, 200, struct {
					PublishableKey string `json:"publishableKey"`
				}{
					PublishableKey: os.Getenv("STRIPE_PUBLISHABLE_KEY"),
				})

				return nil
			}))

			r.Post("/intent", Make(func(w http.ResponseWriter, r *http.Request) error {
				var req Product 
				err := json.NewDecoder(r.Body).Decode(&req)
				if err != nil {
					return Error{
						StatusCode: http.StatusBadRequest,
						Message:    err.Error(),
					}
				}

				params := &stripe.PaymentIntentParams{
					Amount: stripe.Int64(int64(req.Price * 100)),	
					Currency: stripe.String(string(stripe.CurrencyUSD)),
					AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
						Enabled: stripe.Bool(true),
					},
				}

				pi, err := paymentintent.New(params)
				if err != nil {
					if stripeErr, ok := err.(*stripe.Error); ok {
						return Error{
							StatusCode: http.StatusBadRequest,
							Message:    fmt.Errorf("payments provider error: %v", stripeErr.Error()),
						}
					} else {
						return Error{
							StatusCode: http.StatusInternalServerError,
							Message:    err,
						}
					}
				}

				WriteJson(w, 200, struct {
					ClientSecret string `json:"clientSecret"`
				}{
					ClientSecret: pi.ClientSecret,
				})

				return nil
			}))

			r.Post("/webhook", Make(func(w http.ResponseWriter, r *http.Request) error {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					return Error{
						StatusCode: http.StatusInternalServerError,
						Message:    err,
					}
				}
				defer r.Body.Close()

				event, err := webhook.ConstructEvent(
					body, 
					r.Header.Get("Stripe-Signature"), 
					os.Getenv("STRIPE_WEBHOOK_SECRET"),
				)
				if err != nil {
					return Error{
						StatusCode: http.StatusInternalServerError,
						Message:    err,
					}
				}

				if event.Type == "payment_intent.succeeded" {
					fmt.Println("we got the money")
				}

				WriteJson(w, 200, nil)

				return nil
			}))
		})
	})

	http.ListenAndServe(":8080", mux)
}
