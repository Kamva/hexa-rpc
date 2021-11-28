package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/kamva/gutil"
	hrpc "github.com/kamva/hexa-rpc"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var l = hlog.NewPrinterDriver(hlog.DebugLevel)
var t = hexatranslator.NewEmptyDriver()

const (
	addr        = ":2323"
	service     = "hexa-demo"
	environment = "dev"
	id          = 1
)

func initMeter() {
	config := prometheus.Config{}
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		controller.WithResource(resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id))),
	)

	exporter, err := prometheus.New(config, c)

	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}
	global.SetMeterProvider(exporter.MeterProvider())

	http.HandleFunc("/", exporter.ServeHTTP)
	go func() {
		_ = http.ListenAndServe(":2222", nil)
	}()

	fmt.Println("Prometheus server running on :2222")
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "send-request" {

		sendHealthRequest()
		return
	}

	runServer()
}

func sendHealthRequest() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithChainUnaryInterceptor(
		// error interceptor must be the first interceptor.
		hrpc.NewErrorInterceptor().UnaryClientInterceptor(),
	))
	if err != nil {
		gutil.PanicErr(err)
	}
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)

	msg, err := client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	gutil.PanicErr(err)
	fmt.Println(msg.Status)
	hlog.Debug("end of call for healthCheck, bye :)")
}

func runServer() {
	initMeter()

	// Set gRPC default logger:
	grpclog.SetLoggerV2(hrpc.NewLogger(l, 0))

	errOptions := hrpc.ErrInterceptorOptions{
		Logger:       l,
		Translator:   t,
		ReportErrors: true,
	}

	metricsOpts := hrpc.MetricsOptions{
		MeterProvider: global.GetMeterProvider(),
		ServerName:    "lab",
	}

	// Setup hexa context interceptor
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			(&hrpc.Metrics{}).UnaryServerInterceptor(metricsOpts),
			hrpc.NewErrorInterceptor().UnaryServerInterceptor(errOptions),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(hrpc.RecoverHandler)),
		)),
	)

	listener, err := net.Listen("tcp", addr)
	gutil.PanicErr(err)

	grpc_health_v1.RegisterHealthServer(server, hrpc.NewHealthServer()) // Use health server as an example server.
	gutil.PanicErr(server.Serve(listener))
}
