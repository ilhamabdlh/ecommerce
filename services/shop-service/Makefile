.PHONY: all build test clean docker swagger

# Variables
SERVICE_NAME=shop-service
DOCKER_IMAGE=ecommerce/$(SERVICE_NAME)
DOCKER_TAG=latest

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

all: test build

build:
	$(GOBUILD) -o bin/$(SERVICE_NAME) -v ./main.go

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f bin/$(SERVICE_NAME)

docker:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f ../../docker/Dockerfile.shop-service ../..

run: prepare
	$(GOBUILD) -o bin/$(SERVICE_NAME) -v ./main.go
	./bin/$(SERVICE_NAME)

prepare:
	mkdir -p logs
	mkdir -p bin

deps:
	$(GOMOD) download

tidy:
	$(GOMOD) tidy

swagger:
	swag init -g main.go

docker-compose-up:
	cd ../../docker && docker-compose up --build -d

docker-compose-down:
	cd ../../docker && docker-compose down

lint:
	golangci-lint run

.DEFAULT_GOAL := all 