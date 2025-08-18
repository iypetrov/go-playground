generate:
	@mkdir -p gen/common
	@mkdir -p gen/exchange
	@mkdir -p gen/payments
	@protoc -I idl --go_out=gen --go-grpc_out=require_unimplemented_servers=false:gen --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative idl/common/types.proto
	@protoc -I idl --go_out=gen --go-grpc_out=require_unimplemented_servers=false:gen --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative idl/events/events.proto
	@protoc -I idl --go_out=gen --go-grpc_out=require_unimplemented_servers=false:gen --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative idl/tickets/tickets.proto
	@protoc -I idl --go_out=gen --go-grpc_out=require_unimplemented_servers=false:gen --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative idl/payments/payments.proto

events:
	@go run cmd/events/*.go
	
tickets:
	@go run cmd/tickets/*.go

payments:
	@go run cmd/payments/*.go

fmt:
	@go fmt ./...
	@goimports -l -w .
