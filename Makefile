build:
	go build -o build/web main.go

database:
	docker run --name db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5433:5432 -d postgres
	sleep 1;
	docker exec db psql -U postgres -d postgres -c "create database test;"

run:
	go run main.go

generate:
	protoc --go_out=. --go-grpc_out=. audit/proto/service.proto

format:
	go fmt ./...

tidy:
	go mod tidy

grpc-info:
	grpcurl -plaintext localhost:8071 list
	@echo
	grpcurl -plaintext localhost:8071 describe audit.Audit
	@echo
	grpcurl -plaintext localhost:8071 describe .audit.CreateOrderRequest
	@echo
	grpcurl -plaintext localhost:8071 describe .audit.Order

.DEFAULT_GOAL := build

.PHONY: database build