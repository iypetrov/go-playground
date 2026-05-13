build:
	@go build -o ./bin/main

run:
	@docker run --rm --name otel-col-test -d -p 4317:4317 -p 8888:8888 -v $$(pwd)/otel-config.yaml:/etc/otel-config.yaml -e OTEL_LOG_LEVEL=debug otel/opentelemetry-collector-contrib:latest --config /etc/otel-config.yaml
	@./bin/main

teardown:
	@docker rm -f otel-col-test

metrics-plugin:
	@curl http://localhost:2021/metrics

test:
	go test -v

fmt:
	@go fmt ./...
	@goimports -l -w .
