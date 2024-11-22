module github.com/yourusername/ecommerce/user-service

go 1.21

require (
	go.mongodb.org/mongo-driver v1.17.1
	google.golang.org/grpc v1.58.3
	github.com/yourusername/ecommerce/proto v0.0.0
	golang.org/x/crypto v0.29.0
)

replace github.com/yourusername/ecommerce/proto => ../proto 