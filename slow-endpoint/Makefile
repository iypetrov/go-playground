build:
	@go build -o ./bin/main

run:
	@./bin/main

test:
	go test -v

fmt:
	@go fmt ./...
	@goimports -l -w .

curl:
	@curl -fsSm 5 http://127.0.0.1:8080/healthz || echo test
