module github.com/ilhamabdlh/ecommerce/api-gateway

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	google.golang.org/grpc v1.58.3
	github.com/ilhamabdlh/ecommerce/proto v0.0.0
)

replace github.com/ilhamabdlh/ecommerce/proto => ../proto 
