BINARY=taskflow

.PHONY: setup run dev test build docker-up docker-down clean

setup:
	cp .env.example .env
	go mod tidy
	docker-compose up -d db

run:
	go run ./cmd/server

test:
	go test ./tests/... -v -count=1

build:
	go build -o $(BINARY).exe ./cmd/server

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

clean:
	docker-compose down -v
	rm -f $(BINARY).exe
