package metrics

import (
	"auth/internal/config"
	"context"
	"fmt"

	"net/http"

	"go.opentelemetry.io/otel"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var tracer = otel.Tracer("TASKS_SERVICE")

type MetricsContainer struct {
	*jaeger.Exporter
	*tracesdk.TracerProvider
}

func New() (*MetricsContainer, error) {
	cfg, _ := config.GetConfig()

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(cfg.JaegerAddress)),
	)
	if err != nil {
		return nil, fmt.Errorf("init jaeger failed: %w", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
		)))

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":"+cfg.JaegerPort, nil)

	return &MetricsContainer{
		Exporter:       exp,
		TracerProvider: tp,
	}, nil

}

func (mc *MetricsContainer) Stop(ctx context.Context) error {
	err := mc.Exporter.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("jaeger shutdown failed: %w", err)
	}

	err = mc.TracerProvider.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("tracer proviser shutdown failed: %w", err)
	}

	return nil
}
