package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	hrpc "github.com/kamva/hexa-rpc"
	"github.com/kamva/hexa-rpc/examples/simple/hello"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"google.golang.org/grpc"
)

var serverAddr = new(string)

func init() {
	flag.StringVar(serverAddr, "port", "localhost:9010", "gRPC server port")
}

var logger = hlog.NewPrinterDriver(hlog.DebugLevel)
var translator = hexatranslator.NewEmptyDriver()
var p = hexa.NewContextPropagator(logger, translator)

func main() {
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure(), grpc.WithChainUnaryInterceptor(
		// error interceptor must be the first interceptor.
		hrpc.NewErrorInterceptor().UnaryClientInterceptor(),
		hrpc.NewHexaContextInterceptor(p).UnaryClientInterceptor,
	))
	if err != nil {
		gutil.PanicErr(err)
	}
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
	msg, err := client.SayHello(ctx, &hello.Message{Val: "mehran"})
	gutil.PanicErr(err)
	fmt.Println(msg.Val)
	hlog.Debug("end of call 1")

	// Without hexa context
	msg, err = client.SayHello(context.Background(), &hello.Message{Val: "mehran"})
	gutil.PanicErr(err)
	fmt.Println(msg.Val)
	hlog.Debug("end of call 2")

	// Check error converter
	sayHelloWithHexaErr(client)
	hlog.Debug("end of call 3")

	sayHelloWithNativeErr(client)
	hlog.Debug("end of call 4")

	sayHelloWithPanic(client)
	hlog.Debug("end of call 5")

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

func sayHelloWithPanic(client hello.HelloClient) {
	_, err := client.SayHelloWithErr(context.Background(), &hello.Message{Val: "panic"})
	e, ok := err.(hexa.Error)
	if !ok {
		panic("error is not hexa error")
	}
	errorDetails(e)
}

func errorDetails(e hexa.Error) {
	localMsg, err := e.Localize(translator)
	if err != nil {
		fmt.Println("translation err: ", err)
	}

	fmt.Println("--------Hexa error----------")
	fmt.Println("status: ", e.HTTPStatus())
	fmt.Println("id: ", e.ID())
	fmt.Println("data: ", e.Data())
	fmt.Println("local msg: ", localMsg)
	fmt.Println("error: ", e.Error())
	fmt.Println("-----------------------------")
}
