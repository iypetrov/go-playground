build:
	@air -c .air.toml

test:
	go test

test_bench:
	go test -bench=. 

fmt:
	@go fmt ./...
	@goimports -l -w .