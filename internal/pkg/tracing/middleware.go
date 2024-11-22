package tracing

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		tracer := otel.GetTracerProvider().Tracer("")

		propagator := propagation.TraceContext{}
		ctx = propagator.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))

		spanCtx, span := tracer.Start(ctx, c.FullPath())
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.client_ip", c.ClientIP()),
		)

		c.Request = c.Request.WithContext(spanCtx)
		c.Next()

		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
		)
	}
}
