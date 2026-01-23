package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/go-chi/chi/v5"
	"github.com/godruoyi/go-snowflake"
	"github.com/iypetrov/o11y/database"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	"github.com/rs/zerolog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var (
	gaugeNewUsersTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "o11y_new_users_total",
		Help: "Total number of new users",
	})
)

func initTracer(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("init exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp, nil
}

func main() {
	ctx := context.Background()
	serviceName := "our-o11y-service"

	// Logger
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	// Tracing
	tp, err := initTracer(ctx, serviceName)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init tracer")
	}
	defer func() {
		_ = tp.Shutdown(ctx)
	}()

	tracer := otel.Tracer(serviceName)

	// Snowflake
	snowflake.SetMachineID(1)

	// Database
	db, err := otelsql.Open(
		"postgres",
		"postgres://user:pass@localhost:5432/o11y?sslmode=disable",
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	// Migrations
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal().Err(err).Msg("failed to set dialect")
	}
	if err := goose.Up(db, "sql/migrations"); err != nil {
		log.Fatal().Err(err).Msg("failed to apply migrations")
	}

	// Router
	r := chi.NewRouter()
	r.Use(otelchi.Middleware(serviceName, otelchi.WithChiRoutes(r)))
	r.Handle("/metrics", promhttp.Handler())

	r.Post("/user", func(w http.ResponseWriter, r *http.Request) {
		queries := database.New(db)

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, span := tracer.Start(
			r.Context(),
			"CreateUser",
			trace.WithAttributes(
				attribute.String("user.name", user.Name),
				attribute.Int("user.age", user.Age),
			),
		)
		defer span.End()

		result, err := queries.CreateUser(ctx, database.CreateUserParams{
			ID:   int64(snowflake.ID()),
			Name: user.Name,
			Age:  int64(user.Age),
		})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Err(err).Msg("failed to encode response")
			return
		}

		gaugeNewUsersTotal.Inc()
	})

	log.Info().Msg("starting server on :3000")
	log.Fatal().Err(http.ListenAndServe(":3000", r))
}
