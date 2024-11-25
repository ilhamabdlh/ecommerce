package middleware

import (
	"ecommerce/warehouse-service/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqURI := c.Request.RequestURI
		statusCode := c.Writer.Status()

		utils.Logger.Infof("| %3d | %13v | %15s | %s |",
			statusCode,
			latencyTime,
			reqMethod,
			reqURI,
		)
	}
}
