package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
	hrpc "github.com/Kamva/hexa-rpc"
	"github.com/Kamva/hexa-rpc/examples/simple/hello"
	"github.com/Kamva/hexa/db/mgmadapter"
	"github.com/Kamva/hexa/hexalogger"
	"github.com/Kamva/hexa/hexatranslator"
	"google.golang.org/grpc"
)

var serverAddr = new(string)

func init() {
	flag.StringVar(serverAddr, "port", "localhost:9010", "gRPC server port")
}

var logger = hexalogger.NewPrinterDriver()
var translator = hexatranslator.NewEmptyDriver()
var cei = hexa.NewCtxExporterImporter(hexa.NewUserExporterImporter(mgmadapter.EmptyID), logger, translator)

func main() {
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure(), grpc.WithChainUnaryInterceptor(
		// error interceptor must be the first interceptor.
		hrpc.NewErrorInterceptor().UnaryClientInterceptor(),
		hrpc.NewHexaContextInterceptor(cei).UnaryClientInterceptor,
	))
	if err != nil {
		gutil.PanicErr(err)
	}
	defer conn.Close()

	client := hello.NewHelloClient(conn)

	// With Hexa context
	ctx := hexa.NewCtx(nil, "my_correlation_id", "en", hexa.NewGuest(), logger, translator)
	msg, err := client.SayHello(hrpc.Ctx(ctx), &hello.Message{Val: "mehran"})
	gutil.PanicErr(err)
	fmt.Println(msg.Val)

	// Without hexa context
	msg, err = client.SayHello(context.Background(), &hello.Message{Val: "mehran"})
	gutil.PanicErr(err)
	fmt.Println(msg.Val)

	// Check error converter
	sayHelloWithHexaErr(client)
	sayHelloWithNativeErr(client)

}

func sayHelloWithNativeErr(client hello.HelloClient) {
	_, err := client.SayHelloWithErr(context.Background(), &hello.Message{Val: "john"})
	e, ok := err.(hexa.Error)
	if !ok {
		panic("error is not hexa error")
	}
	errorDetails(e)
}

func sayHelloWithHexaErr(client hello.HelloClient) {
	_, err := client.SayHelloWithErr(context.Background(), &hello.Message{Val: "mehran"})
	e, ok := err.(hexa.Error)
	if !ok {
		panic("error is not hexa error")
	}
	errorDetails(e)
}

func errorDetails(e hexa.Error) {
	localMsg, err := e.Localize(translator)
	if err!=nil{
		fmt.Println("translation err: ",err)
	}

	fmt.Println("--------Hexa error----------")
	fmt.Println("status: ", e.HTTPStatus())
	fmt.Println("code: ", e.Code())
	fmt.Println("data: ", e.Data())
	fmt.Println("local msg: ", localMsg)
	fmt.Println("error: ", e.Error())
	fmt.Println("-----------------------------")
}
