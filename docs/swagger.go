package docs

import "github.com/swaggo/swag"

var SwaggerInfo = &swag.Spec{
	Version:     "1.0",
	Host:        "localhost:8080",
	BasePath:    "/api/v1",
	Title:       "E-Commerce API",
	Description: "API Server for E-Commerce Application",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
