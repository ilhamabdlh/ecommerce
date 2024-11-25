.PHONY: all build run test clean docker-build docker-run swagger lint help

# Go parameters
BINARY_NAME=warehouse-service
MAIN_FILE=main.go
DOCKER_IMAGE=warehouse-service
DOCKER_TAG=latest

# Environment variables
export GO111MODULE=on

all: clean lint test build

help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build files"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"
	@echo "  make swagger       - Generate swagger documentation"
	@echo "  make lint          - Run linter"
	@echo "  make mock          - Generate mocks"
	@echo "  make all           - Clean, lint, test, and build"

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) $(MAIN_FILE)

run:
	@echo "Running $(BINARY_NAME)..."
	go run $(MAIN_FILE)

test:
	@echo "Running tests..."
	go test -v -race -cover ./...

clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean
	rm -f $(BINARY_NAME)

docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f docker/Dockerfile.warehouse-service .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8084:8084 -p 9090:9090 $(DOCKER_IMAGE):$(DOCKER_TAG)

swagger:
	@echo "Generating Swagger documentation..."
	swag init -g $(MAIN_FILE)

lint:
	@echo "Running linter..."
	golangci-lint run

mock:
	@echo "Generating mocks..."
	mockgen -source=internal/domain/warehouse.go -destination=mocks/mock_warehouse.go

# Development commands
dev:
	@echo "Running in development mode..."
	air

deps:
	@echo "Installing dependencies..."
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/cosmtrek/air@latest

# Docker Compose commands
compose-up:
	@echo "Starting services with Docker Compose..."
	docker-compose -f docker/development/docker-compose.yml up -d

compose-down:
	@echo "Stopping services..."
	docker-compose -f docker/development/docker-compose.yml down

# Kubernetes commands
k8s-apply:
	@echo "Applying Kubernetes manifests..."
	kubectl apply -f k8s/warehouse-service/

k8s-delete:
	@echo "Deleting Kubernetes resources..."
	kubectl delete -f k8s/warehouse-service/

# Database commands
db-migrate:
	@echo "Running database migrations..."
	# Add your migration command here

db-seed:
	@echo "Seeding database..."
	# Add your seeding command here

# Default target
.DEFAULT_GOAL := help 