.PHONY: build run test migrate clean docker-up docker-down

build:
	go build -o gym-strongcode cmd/main.go

run:
	go run cmd/main.go

test:
	go test -v ./tests/...

test-unit:
	go test -v ./tests/unit/...

test-integration:
	go test -v ./tests/integration/...

migrate:
	go run cmd/migrate/main.go

clean:
	rm -f gym-strongcode
	rm -rf logs/*

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

swagger:
	swag init -g cmd/main.go

lint:
	golangci-lint run

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out