module github.com/ilhamabdlh/ecommerce/warehouse-service

go 1.21

require (
	go.mongodb.org/mongo-driver v1.17.1
	google.golang.org/grpc v1.58.3
	github.com/ilhamabdlh/ecommerce/proto v0.0.0
)

replace github.com/ilhamabdlh/ecommerce/proto => ../proto 
