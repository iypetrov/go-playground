IMAGE     ?= iypetrov/daemon
TAG       ?= latest
PLATFORMS ?= linux/amd64,linux/arm64

build:
	@go build -o ./bin/main

run:
	@./bin/main

test:
	go test -v

fmt:
	@go fmt ./...
	@goimports -l -w .

docker-build:
	@docker buildx build --platform $(PLATFORMS) -t $(IMAGE):$(TAG) .

docker-push:
	@docker buildx build --platform $(PLATFORMS) -t $(IMAGE):$(TAG) --push .

.PHONY: build fmt test
