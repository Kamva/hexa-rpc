package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/kamva/hexa"
	hrpc "github.com/kamva/hexa-rpc"
	"github.com/kamva/hexa-rpc/examples/simple/hello"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/hexa/sr"
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

var r = sr.New()

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

	hexaSrv :=hrpc.NewHexaService(hrpc.NewGRPCHealth("my_health_server", fmt.Sprintf(":%d", *port)),lis,grpcServer)
	r.Register("grpc_server", hexaSrv)
	go func() { sr.ShutdownBySignals(r, time.Second*30) }()
	_ = hexaSrv.(hexa.Runnable).Run()
}
