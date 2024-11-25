.PHONY: test build run docker-build docker-run

# Go related variables
BINARY_NAME=user-service
MAIN_FILE=main.go

# Docker related variables
DOCKER_IMAGE=ecommerce/user-service
DOCKER_TAG=latest

# Build the project
build:
	go build -o ${BINARY_NAME} ${MAIN_FILE}

# Run the project
run:
	go run ${MAIN_FILE}

# Run tests
test:
	go test -v ./...

# Run integration tests
test-integration:
	go test -v ./tests/integration/...

# Build docker image
docker-build:
	docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} -f ../../docker/Dockerfile.user-service ../..

# Run docker container
docker-run:
	docker run -p 8081:8081 ${DOCKER_IMAGE}:${DOCKER_TAG}

# Generate swagger documentation
swagger:
	swag init

# Clean up
clean:
	go clean
	rm -f ${BINARY_NAME}

# Install dependencies
deps:
	go mod download

# Run linter
lint:
	golangci-lint run

# Run all tests and lint
check: test lint 