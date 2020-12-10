package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/kamva/hexa"
	hrpc "github.com/kamva/hexa-rpc"
	"github.com/kamva/hexa-rpc/examples/simple/hello"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var port = new(int)

func init() {
	flag.IntVar(port, "port", 9010, "gRPC server port")
}

var logger = hlog.NewPrinterDriver(hlog.DebugLevel)
var translator = hexatranslator.NewEmptyDriver()
var cei = hexa.NewContextPropagator(logger, translator)

func main() {
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
			hrpc.NewHexaContextInterceptor(cei).UnaryServerInterceptor,
			hrpc.NewRequestLogger(logger).UnaryServerInterceptor(hrpc.DefaultLoggerOptions(true)),
			hrpc.NewErrorInterceptor().UnaryServerInterceptor(errOptions),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(hrpc.RecoverHandler)),
		)),
	)
	hello.RegisterHelloServer(grpcServer, hello.New())
	_ = grpcServer.Serve(lis)
}
