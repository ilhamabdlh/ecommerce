# Product Service

## Project Structure

Product service is a microservice that handles product management in the e-commerce system.

## Features

- Product CRUD operations
- Stock management
- Category management
- Authentication and authorization
- Metrics and monitoring
- Health checks

## Technologies

- Go 1.21
- Gin Web Framework
- MongoDB
- Docker
- Kubernetes
- Prometheus (metrics)
- Swagger (API documentation)

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker
- MongoDB
- Make

### Installation

1. Clone the repository 
bash
git clone https://github.com/ilhamabdlh/ecommerce.git
2. Navigate to the product service directory
bash
cd services/product-service
3. Install dependencies
bash
make deps
4. Build the service
bash
make build

### Running the Service

#### Locally
bash
make run

### API Documentation

The API documentation is available at `/swagger/index.html` when the service is running.

### Testing

Run tests:
bash
make test

Run linter:
bash
make lint

## API Endpoints

- `GET /api/v1/products` - List all products
- `POST /api/v1/products` - Create a new product
- `GET /api/v1/products/:id` - Get product by ID
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product
- `PATCH /api/v1/products/:id/stock` - Update product stock

## Monitoring

- Health check: `/health`
- Readiness check: `/ready`
- Metrics: `/metrics`

## Contributing

Please read CONTRIBUTING.md for details on our code of conduct, and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details