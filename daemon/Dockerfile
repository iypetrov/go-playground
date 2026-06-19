FROM --platform=$BUILDPLATFORM golang:1.26.2 AS build-stage
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o bin/main .

FROM gcr.io/distroless/base-debian13 AS run-stage
COPY --from=build-stage /app/bin/main /bin/main
CMD ["/bin/main"]
