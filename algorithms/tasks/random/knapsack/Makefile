build:
	@go mod tidy
	@go mod download
	@go build -o ./bin/main .

run:
	@./bin/main

test:
	go test -v

fmt:
	@go fmt ./...
	@goimports -l -w .
