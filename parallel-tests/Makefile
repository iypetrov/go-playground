build:
	@go build -o ./bin/main

run:
	@./bin/main

fmt:
	@go fmt ./...
	@goimports -l -w .

test:
	go test -v

vgt:
	@GOMAXPROCS=1 go test -count=1 -parallel 128 -p 16 -json ./... | vgt --print-html > foo.html
	@sudo -u ipetrov xdg-open foo.html
