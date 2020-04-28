package main

import (
	"flag"
	"fmt"
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
	hrpc "github.com/Kamva/hexa-rpc"
	"github.com/Kamva/hexa-rpc/examples/simple/hello"
	"github.com/Kamva/hexa-rpc/examples/simple/service"
	"github.com/Kamva/hexa/db/mgmadapter"
	"github.com/Kamva/hexa/hexaconfig"
	"github.com/Kamva/hexa/hexalogger"
	"github.com/Kamva/hexa/hexatranslator"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"net"
)

var port = new(int)

func init() {
	flag.IntVar(port, "port", 9010, "gRPC server port")
}

var logger = hexalogger.NewPrinterDriver()
var translator = hexatranslator.NewEmptyDriver()
var cei = hexa.NewCtxExporterImporter(hexa.NewUserExporterImporter(mgmadapter.EmptyID), logger, translator)
var cfg = hexaconfig.NewMapDriver()

func init() {
	gutil.PanicErr(cfg.Unmarshal(hexa.Map{
		hrpc.GRPCLogVerbosityLevel: int64(0),
	}))
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Set gRPC default logger:
	grpclog.SetLoggerV2(hrpc.NewLogger(logger, cfg))

	// Setup hexa context interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			hrpc.NewHexaContextInterceptor(cei).UnaryServerInterceptor,
			hrpc.NewRequestLogger(logger).UnaryServerInterceptor(hrpc.DefaultLoggerOptions(true)),
			// Error converter must be last interceptor
			hrpc.NewErrorInterceptor().UnaryServerInterceptor(translator),
		)),
	)
	hello.RegisterHelloServer(grpcServer, service.New())
	grpcServer.Serve(lis)
}
