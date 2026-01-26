package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
)

var (
	tp         *sdktrace.TracerProvider
	tracerOnce sync.Once
)

func InitTracer() {
	tracerOnce.Do(func() {
		endpoint := os.Getenv("JAEGER_URL")
		serviceName := os.Getenv("SERVICE_NAME")

		//test connectiion manual

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		conn, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("❌ Cannot connect to OTLP endpoint (%s): %v\n💥 Shutting down application.", endpoint, err)
		}
		_ = conn.Close()

		//end test connection manual

		exporter, err := otlptracegrpc.New(ctx,
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithTimeout(5*time.Second),
		)
		if err != nil {
			log.Fatalf("❌ Failed to create OTLP exporter: %v\n💥 Shutting down application.", err)
		}

		res, err := resource.New(ctx,
			resource.WithAttributes(
				semconv.ServiceNameKey.String(serviceName),
				semconv.ServiceVersionKey.String("1.0.0"),
			),
		)
		if err != nil {
			log.Fatalf("❌ Failed to create resource: %v\n💥 Shutting down application.", err)
		}

		tp = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter,
				sdktrace.WithBatchTimeout(2*time.Second),
				sdktrace.WithMaxExportBatchSize(100),
			),
			sdktrace.WithResource(res),
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
		)

		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

		fmt.Printf("✅ Tracing initialized with OTLP gRPC: %s (service: %s)\n", endpoint, serviceName)
	})
}

// TAMBAHAN: Test connection ke exporter
func testExporterConnection(ctx context.Context, exporter sdktrace.SpanExporter) error {
	// Create a test span untuk verify connection
	testTP := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
	)

	tracer := testTP.Tracer("test")
	_, span := tracer.Start(ctx, "connection-test")
	span.End()

	// Force export
	return testTP.ForceFlush(ctx)
}

func ShutdownTracer(ctx context.Context) error {
	if tp != nil {
		log.Println("🔄 Shutting down tracer...")
		return tp.Shutdown(ctx)
	}
	return nil
}

// TAMBAHAN: Helper function untuk check apakah tracing enabled
func IsTracingEnabled() bool {
	return tp != nil
}
