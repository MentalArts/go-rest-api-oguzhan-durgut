package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func Telemetry() gin.HandlerFunc {
	tracer := otel.GetTracerProvider().Tracer("gin-middleware")

	return func(c *gin.Context) {
		carrier := propagation.HeaderCarrier(c.Request.Header)
		ctx := c.Request.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

		spanName := c.FullPath()
		if spanName == "" {
			spanName = c.Request.URL.Path
		}

		ctx, span := tracer.Start(
			ctx,
			spanName,
			trace.WithAttributes(
				attribute.String("http.request.method", c.Request.Method),
				attribute.String("http.route", c.FullPath()),
				attribute.String("http.target", c.Request.URL.Path),
			),
		)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		status := c.Writer.Status()
		span.SetAttributes(
			attribute.Int("http.response.status_code", status),
			attribute.String("http.response.status_text", http.StatusText(status)),
			attribute.Int64("http.response.written_bytes", int64(c.Writer.Size())),
		)

		if status >= 400 {
			span.SetStatus(codes.Error, http.StatusText(status))
		}
	}
}
