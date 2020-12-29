package main

import (
	"context"
	"fmt"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	hrpc "github.com/kamva/hexa-rpc"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var l = hlog.NewPrinterDriver(hlog.DebugLevel)
var t = hexatranslator.NewEmptyDriver()

const addr = "127.0.0.1:2323"

func main() {
	p := hexa.NewContextPropagator(l, t)
	// Set gRPC default logger:
	grpclog.SetLoggerV2(hrpc.NewLogger(l, 0))

	errOptions := hrpc.ErrInterceptorOptions{
		Logger:       l,
		Translator:   t,
		ReportErrors: true,
	}

	// Setup hexa context interceptor
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			hrpc.NewHexaContextInterceptor(p).UnaryServerInterceptor,
			hrpc.NewRequestLogger(l).UnaryServerInterceptor(hrpc.DefaultLoggerOptions(true)),
			hrpc.NewErrorInterceptor().UnaryServerInterceptor(errOptions),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(hrpc.RecoverHandler)),
		)),
	)

	listener, err := net.Listen("tcp", addr)
	gutil.PanicErr(err)

	grpc_health_v1.RegisterHealthServer(server, hrpc.NewHealthServer())
	go func() {
		server.Serve(listener)
	}()

	time.Sleep(time.Second) // Wait to run the server

	health := hrpc.NewGRPCHealth("my_grpc_server", addr)
	fmt.Println(health.HealthStatus(context.Background()))

	server.GracefulStop()
}
