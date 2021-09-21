package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
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
	"google.golang.org/grpc/grpclog"
)

const (
	service     = "grpc_demo_hello_service"
	environment = "dev"
	id          = 1
)

var port = new(int)

func init() {
	flag.IntVar(port, "port", 9010, "gRPC server port")
}

var logger = hlog.NewPrinterDriver(hlog.DebugLevel)
var translator = hexatranslator.NewEmptyDriver()
var cei = hexa.NewContextPropagator(logger, translator)

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

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Set gRPC default logger:
	grpclog.SetLoggerV2(hrpc.NewLogger(logger, 0))

	errOptions := hrpc.ErrInterceptorOptions{
		Logger:       logger,
		Translator:   translator,
		ReportErrors: true,
	}

	// Setup hexa context interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(otelPropagator)),
			hrpc.NewHexaContextInterceptor(cei).UnaryServerInterceptor,
			hrpc.NewRequestLogger(logger).UnaryServerInterceptor(hrpc.DefaultLoggerOptions(true)),
			hrpc.NewErrorInterceptor().UnaryServerInterceptor(errOptions),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(hrpc.RecoverHandler)),
		)),
	)
	hello.RegisterHelloServer(grpcServer, hello.New())
	_ = grpcServer.Serve(lis)
}
