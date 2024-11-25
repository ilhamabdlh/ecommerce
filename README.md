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

### 5. Warehouse Service (Port: 8085)
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
   curl -X POST http://localhost:8085/api/v1/auth/token \
     -H "Content-Type: application/json" \
     -d '{
       "user_id": "user123",
       "roles": ["admin"]
     }'

   # Save token
   export TOKEN="<token-from-response>"

   # Create warehouse
   curl -X POST http://localhost:8085/api/v1/warehouses \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{
       "name": "Main Warehouse",
       "location": "Jakarta"
     }'

   # Get all warehouses
   curl http://localhost:8085/api/v1/warehouses \
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
| Warehouse Service | http://localhost:8085 | http://localhost:8085/swagger/index.html |

### Monitoring

| Service | URL | Credentials |
|---------|-----|-------------|
| Prometheus | http://localhost:9090 | N/A |
| Grafana | http://localhost:3000 | admin/admin |
| Service Metrics | http://localhost:8084/metrics | N/A |

### Development Commands

### API Documentation

#### Postman Collection

[![Run in Postman](https://run.pstmn.io/button.svg)](https://warped-eclipse-691260.postman.co/workspace/binary~a677502a-2330-4346-9ae5-c38535abca72/collection/36495103-05fa9abb-2a53-4234-a472-b7bf1031a696?action=share&creator=36495103)

Collection includes:

1. Order Service Endpoints:
   - POST `/api/v1/orders` - Create new order
   - GET `/api/v1/orders/{order_id}` - Get order by ID
   - GET `/api/v1/orders/user` - Get user orders
   - GET `/api/v1/orders/shop/{shop_id}` - Get shop orders
   - POST `/api/v1/orders/{order_id}/process` - Process order
   - POST `/api/v1/orders/{order_id}/complete` - Complete order
   - POST `/api/v1/orders/{order_id}/cancel` - Cancel order

2. Product Service Endpoints:
   - POST `/api/v1/products` - Create product
   - GET `/api/v1/products` - Get all products
   - GET `/api/v1/products/{product_id}` - Get product by ID
   - PUT `/api/v1/products/{product_id}` - Update product
   - PATCH `/api/v1/products/{product_id}/stock` - Update stock
   - DELETE `/api/v1/products/{product_id}` - Delete product

3. Shop Service Endpoints:
   - GET `/api/v1/shops/{shop_id}` - Get shop by ID
   - GET `/api/v1/shops` - Get all shops
   - POST `/api/v1/shops` - Create shop
   - PUT `/api/v1/shops/{shop_id}` - Update shop
   - POST `/api/v1/shops/{shop_id}/warehouses/{warehouse_id}` - Add warehouse to shop
   - DELETE `/api/v1/shops/{shop_id}/warehouses/{warehouse_id}` - Remove warehouse from shop
   - DELETE `/api/v1/shops/{shop_id}` - Delete shop

4. User Service Endpoints:
   - POST `/api/v1/auth/register` - Register new user
   - POST `/api/v1/auth/login` - Login user
   - GET `/api/v1/users/me` - Get current user
   - PUT `/api/v1/users/me` - Update user profile

5. Warehouse Service Endpoints:
   - PUT `/api/v1/warehouses/{warehouse_id}/stock` - Update stock
   - POST `/api/v1/warehouses` - Create warehouse
   - POST `/api/v1/auth/token` - Generate JWT token
   - GET `/api/v1/warehouses` - Get all warehouses
   - GET `/api/v1/warehouses/{warehouse_id}` - Get warehouse by ID
   - PUT `/api/v1/warehouses/{id}/stock` - Update stock
   - POST `/api/v1/warehouses/transfer` - Transfer stock
   - PUT `/api/v1/warehouses/{id}/{status}` - Update warehouse status (activate/deactivate)

