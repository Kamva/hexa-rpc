package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	hrpc "github.com/kamva/hexa-rpc"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/hexa/probe"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var l = hlog.NewPrinterDriver(hlog.DebugLevel)
var t = hexatranslator.NewEmptyDriver()

const addr = ":2323"
const probeAddr = ":2300"

type ThreadChecker struct {
}

func (t ThreadChecker) HealthIdentifier() string {
	return "thread_checker"
}

func (t ThreadChecker) LivenessStatus(ctx context.Context) hexa.LivenessStatus {
	return hexa.StatusAlive
}

func (t ThreadChecker) ReadinessStatus(ctx context.Context) hexa.ReadinessStatus {
	return hexa.StatusReady
}

func (t ThreadChecker) HealthStatus(ctx context.Context) hexa.HealthStatus {
	return hexa.HealthStatus{
		Id: t.HealthIdentifier(),
		Tags: map[string]string{
			"thread": fmt.Sprint(runtime.NumGoroutine()),
		},
		Alive: hexa.StatusAlive,
		Ready: hexa.StatusReady,
	}
}

var _ hexa.Health = &ThreadChecker{}

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

	hr := hexa.NewHealthReporter()
	hr.AddToChecks(hrpc.NewGRPCHealth("grpc_server", addr))
	hr.AddToChecks(ThreadChecker{})
	ps := probe.NewServer(&http.Server{Addr: probeAddr}, http.NewServeMux())
	probe.RegisterHealthHandlers(ps, hr)
	gutil.PanicErr(ps.Run())

	gutil.PanicErr(server.Serve(listener))
}
