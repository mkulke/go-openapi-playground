package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/mkulke/go-openapi-playground/api"
	otelmiddleware "go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// func getCorrelationID(ctx echo.Context) string {
// 	for key, values := range ctx.Request().Header {
// 		if strings.ToLower(key) == "x-correlation-id" {
// 			for _, value := range values {
// 				return value
// 			}
// 		}
// 	}
// 	return uuid.New().String()
// }

const (
	serviceName    = "go-openapi-playground"
	version        = "0.0.1"
	jaegerEndpoint = "http://localhost:14268/api/traces"
)

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("go-openapi-playground"),
		semconv.ServiceVersionKey.String("0.0.1"),
	)
}

func newTracerProvider() (*sdktrace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))

	if err != nil {
		// fmt.Printf("Failed to create Jaeger trace exporter: %v\n", err)
		// os.Exit(1)
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(newResource()),
	)

	return tp, nil
}

func main() {
	var port = flag.Int("port", 8080, "Port for test HTTP server")
	flag.Parse()

	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Jaeger tracing
	tp, err := newTracerProvider()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error init tracing logic\n: %s", err)
		os.Exit(1)
	}
	ctx := context.Background()
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)

	// This is how you set up a basic Echo router
	e := echo.New()

	// Add otel instrumentation
	hn, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot get Hostname\n: %s", err)
		os.Exit(1)
	}
	e.Use(otelmiddleware.Middleware(hn))

	// Log all requests
	e.Use(echomiddleware.Logger())

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	e.Use(middleware.OapiRequestValidator(swagger))

	// We now register our petStore above as the handler for the interface
	api.RegisterHandlers(e, &api.Api{})

	// And we serve HTTP until the world ends.
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", *port)))
}
