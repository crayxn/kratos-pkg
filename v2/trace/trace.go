package trace

import (
	"context"
	kc "github.com/go-kratos/kratos/v2/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type Tracer struct {
	Exporter tracesdk.SpanExporter
}

func NewTracer(kc kc.Config) (*Tracer, func()) {
	exp := func() tracesdk.SpanExporter {
		driver, _ := kc.Value("trace.driver").String()
		grpc, _ := kc.Value("trace.grpc").Bool()
		switch driver {
		case "jaeger":
			endpoint, _ := kc.Value("trace.endpoint").String()
			var exporter tracesdk.SpanExporter
			var err error
			if grpc {
				exporter, err = otlptracegrpc.New(context.Background(), otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithInsecure())
			} else {
				exporter, err = otlptracehttp.New(context.Background(), otlptracehttp.WithEndpoint(endpoint), otlptracehttp.WithInsecure())
			}
			if err != nil {
				panic(err)
			}
			return exporter
		default:
			return &tracetest.NoopExporter{}
		}
	}()
	return &Tracer{exp}, func() { _ = exp.Shutdown(context.Background()) }
}

func (t *Tracer) InitOTELProvider(attrs ...attribute.KeyValue) {
	tp := tracesdk.NewTracerProvider(
		// 将基于父span的采样率设置为100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// 始终确保在生产中批量处理
		tracesdk.WithBatcher(t.Exporter),
		// 在资源中记录有关此应用程序的信息
		tracesdk.WithResource(resource.NewSchemaless(attrs...)),
	)
	otel.SetTracerProvider(tp)
}
