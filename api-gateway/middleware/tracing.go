package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TracingMiddleware(serviceName string) gin.HandlerFunc {
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := propagation.TraceContext{}

	return func(c *gin.Context) {
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		spanCtx, span := tracer.Start(ctx, c.FullPath(),
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
			),
		)
		defer span.End()

		c.Request = c.Request.WithContext(spanCtx)
		c.Next()

		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
		)
	}
}
