package main

import (
	"context"
	"flag"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	hrpc "github.com/kamva/hexa-rpc"
	"github.com/kamva/hexa-rpc/examples/simple/hello"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

const (
	service     = "grpc_demo_client"
	environment = "dev"
	id          = 1
)

var serverAddr = new(string)

func init() {
	flag.StringVar(serverAddr, "port", "localhost:9010", "gRPC server port")
}

var logger = hlog.NewPrinterDriver(hlog.DebugLevel)
var translator = hexatranslator.NewEmptyDriver()
var p = hexa.NewContextPropagator(logger, translator)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}

func main() {
	// Navigate to the http://localhost:16686 to access to the Jaeger UI.
	tp, err := tracerProvider("http://localhost:14268/api/traces")
	gutil.PanicErr(err)
	defer func() { tp.Shutdown(context.Background()) }()

	otelPropagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure(), grpc.WithChainUnaryInterceptor(
		// error interceptor must be the first interceptor.
		hrpc.NewErrorInterceptor().UnaryClientInterceptor(),
		hrpc.NewHexaContextInterceptor(p).UnaryClientInterceptor,
		otelgrpc.UnaryClientInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(otelPropagator)),
	))
	gutil.PanicErr(err)
	defer conn.Close()

	client := hello.NewHelloClient(conn)

	// With Hexa context
	ctx := hexa.NewContext(nil, hexa.ContextParams{
		CorrelationId: "my_correlation_id",
		Locale:        "en",
		User:          hexa.NewGuest(),
		Logger:        logger,
		Translator:    translator,
	})

	_, _ = client.SayHello(ctx, &hello.Message{Val: "john"})
	_, _ = client.SayHelloWithErr(ctx, &hello.Message{Val: "john_with_error"})
}
