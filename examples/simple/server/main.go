package main

import (
	"flag"
	"fmt"
	"github.com/Kamva/hexa"
	hrpc "github.com/Kamva/hexa-rpc"
	"github.com/Kamva/hexa-rpc/examples/simple/hello"
	"github.com/Kamva/hexa-rpc/examples/simple/service"
	"github.com/Kamva/hexa/db/mgmadapter"
	"github.com/Kamva/hexa/hexalogger"
	"github.com/Kamva/hexa/hexatranslator"
	"google.golang.org/grpc"
	"log"
	"net"
)

var port = new(int)

func init() {
	flag.IntVar(port, "port", 9090, "gRPC server port")
}

var logger = hexalogger.NewPrinterDriver()
var translator = hexatranslator.NewEmptyDriver()
var cei = hexa.NewCtxExporterImporter(hexa.NewUserExporterImporter(mgmadapter.EmptyID), logger, translator)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	hexaCtxInt := hrpc.NewHexaContextInterceptor(cei)
	// Setup hexa context interceptor
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(hexaCtxInt.UnaryServerInterceptor))
	hello.RegisterHelloServer(grpcServer, service.New())
	grpcServer.Serve(lis)
}
