build:
	go build -o build/web main.go

database:
	docker run --name db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5433:5432 -d postgres
	sleep 1;
	docker exec db psql -U postgres -d postgres -c "create database test;"

run:
	go run main.go

format:
	go fmt ./...

tidy:
	go mod tidy

.DEFAULT_GOAL := build

.PHONY: database build