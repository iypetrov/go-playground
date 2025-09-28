generate:
	@mkdir -p gen/chat && protoc -I idl --go_out=gen --go-grpc_out=require_unimplemented_servers=false:gen --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative idl/chat/chat.proto

chat:
	@go run cmd/chat/*.go
	
fmt:
	@go fmt ./...
	@goimports -l -w .
