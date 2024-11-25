# E-commerce Microservices

A modern e-commerce platform built with microservices architecture using Go, MongoDB, Docker, and Kubernetes.

## Services Overview

### 1. User Service (Port: 8081)
Handles user management and authentication:
- User registration and login
- Profile management
- JWT authentication
- Role-based access control (RBAC)

### 2. Product Service (Port: 8082)
Manages product catalog:
- Product CRUD operations
- Product categories
- Product search
- Stock availability check

### 3. Shop Service (Port: 8083)
Manages shop operations:
- Shop registration
- Shop management
- Shop products
- Shop analytics

### 4. Order Service (Port: 8084)
Handles order processing:
- Order creation
- Order status management
- Payment integration
- Order history

### 5. Warehouse Service (Port: 8084)
Manages inventory and stock:
- Stock management
- Multi-warehouse support
- Stock transfer between warehouses
- Stock level monitoring
- Automated stock release

## Technology Stack

### Backend
- Go (Gin Framework)
- MongoDB
- JWT Authentication
- Swagger Documentation

### Infrastructure
- Docker
- Kubernetes
- Prometheus & Grafana
- Consul (Service Discovery)

### Tools & Libraries
- Air (Hot Reload)
- Zap (Logging)
- Prometheus (Metrics)
- Circuit Breaker
- Rate Limiter

## Features

- Microservices Architecture
- RESTful APIs
- JWT Authentication
- Role-Based Access Control
- API Documentation (Swagger)
- Monitoring & Metrics
- Distributed Logging
- Circuit Breaker Pattern
- Rate Limiting
- Service Discovery
- Load Balancing
- Horizontal Scaling

## Quick Start Guide

### Prerequisites
1. Install required tools:
   ```bash
   # Install Go
   brew install go  # For MacOS
   # or download from https://golang.org/dl/

   # Install Docker
   brew install docker  # For MacOS
   # or follow https://docs.docker.com/get-docker/

   # Install Kubernetes tools
   brew install kubectl
   brew install minikube

   # Install MongoDB
   brew install mongodb-community  # For MacOS
   # or use Docker: docker run -d -p 27017:27017 mongo
   ```

2. Clone the repository:
   ```bash
   git clone <repository-url>
   cd ecommerce
   ```

### Running Individual Services

1. Start MongoDB:
   ```bash
   # Using local MongoDB
   brew services start mongodb-community

   # Or using Docker
   docker run -d -p 27017:27017 --name mongodb mongo
   ```

2. Run Warehouse Service:
   ```bash
   cd services/warehouse-service

   # Install dependencies
   make deps

   # Generate Swagger docs
   make swagger

   # Run service
   make run
   # or for development with hot reload:
   make dev
   ```

3. Test the API:
   ```bash
   # Generate JWT token
   curl -X POST http://localhost:8084/api/v1/auth/token \
     -H "Content-Type: application/json" \
     -d '{
       "user_id": "user123",
       "roles": ["admin"]
     }'

   # Save token
   export TOKEN="<token-from-response>"

   # Create warehouse
   curl -X POST http://localhost:8084/api/v1/warehouses \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{
       "name": "Main Warehouse",
       "location": "Jakarta"
     }'

   # Get all warehouses
   curl http://localhost:8084/api/v1/warehouses \
     -H "Authorization: Bearer $TOKEN"
   ```

### Running with Docker Compose

1. Start all services:
   ```bash
   docker-compose -f docker/docker-compose.yml up -d
   ```

2. Start specific service:
   ```bash
   docker-compose -f docker/docker-compose.yml up -d warehouse-service
   ```

3. View logs:
   ```bash
   docker-compose -f docker/docker-compose.yml logs -f warehouse-service
   ```

4. Stop services:
   ```bash
   docker-compose -f docker/docker-compose.yml down
   ```

### Development Environment

1. Start development environment:
   ```bash
   cd docker/development
   docker-compose up -d
   ```

2. View logs:
   ```bash
   docker-compose logs -f
   ```

3. Stop development environment:
   ```bash
   docker-compose down
   ```

### Kubernetes Deployment

1. Start Minikube:
   ```bash
   minikube start
   ```

2. Apply Kubernetes configurations:
   ```bash
   kubectl create namespace ecommerce
   kubectl apply -f k8s/warehouse-service/
   ```

3. Verify deployment:
   ```bash
   kubectl get pods -n ecommerce
   kubectl get services -n ecommerce
   ```

4. Access service:
   ```bash
   # Get service URL
   minikube service warehouse-service -n ecommerce --url
   ```

### Accessing Services

| Service | Local URL | Swagger Documentation |
|---------|-----------|---------------------|
| User Service | http://localhost:8081 | http://localhost:8081/swagger/index.html |
| Product Service | http://localhost:8082 | http://localhost:8082/swagger/index.html |
| Shop Service | http://localhost:8083 | http://localhost:8083/swagger/index.html |
| Order Service | http://localhost:8084 | http://localhost:8084/swagger/index.html |
| Warehouse Service | http://localhost:8084 | http://localhost:8084/swagger/index.html |

### Monitoring

| Service | URL | Credentials |
|---------|-----|-------------|
| Prometheus | http://localhost:9090 | N/A |
| Grafana | http://localhost:3000 | admin/admin |
| Service Metrics | http://localhost:8084/metrics | N/A |

### Development Commands
