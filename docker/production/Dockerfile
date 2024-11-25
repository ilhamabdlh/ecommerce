FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o warehouse-service main.go

# Create final image
FROM alpine:latest

WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/warehouse-service .
COPY --from=builder /app/config/config.go ./config/

# Create necessary directories
RUN mkdir -p /app/logs

# Set timezone
ENV TZ=UTC

# Expose ports
EXPOSE 8084
EXPOSE 9090

# Run the application
CMD ["./warehouse-service"] 