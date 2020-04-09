package hgrpc

import (
	"context"
	"encoding/json"
	"github.com/Kamva/hexa"
	"github.com/Kamva/tracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ContextKeyHexaCtx is the identifier to set the hexa context as a field in the context of a gRPC method.
const ContextKeyHexaCtx = "__hexa_ctx__"

var (
	ErrInvalidHexaContextPayload = status.Error(codes.Internal, "invalid hexa context payload provided to json marshaller")
)

// HexaContextInterceptor is the gRPC interceptor to pass hexa context through gRPC.
// Note: we do not provide stream interceptors, if you think need it, create PR or issue.
type HexaContextInterceptor struct {
	cei hexa.ContextExporterImporter
}

func (ci *HexaContextInterceptor) UnaryClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	hexaCtx := ctx.Value(ContextKeyHexaCtx)
	if hexaCtx != nil {
		m, err := ci.cei.Export(hexaCtx.(hexa.Context))
		if err != nil {
			return tracer.Trace(err)
		}
		ctxBytes, err := json.Marshal(m)
		if err != nil {
			return tracer.Trace(err)
		}

		ctx = metadata.AppendToOutgoingContext(ctx, ContextKeyHexaCtx, string(ctxBytes))
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}

func (ci *HexaContextInterceptor) UnaryServerInterceptor(c context.Context, req interface{}, info *grpc.UnaryServerInfo, h grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return h(c, req)
	}

	ctxStrArr, ok := md[ContextKeyHexaCtx]
	if ok {
		m := make(hexa.Map)
		err := json.Unmarshal([]byte(ctxStrArr[0]), &m)
		if err != nil {
			return nil, tracer.Trace(ErrInvalidHexaContextPayload)
		}

		hexaCtx, err := ci.cei.Import(m)
		if err != nil {
			return nil, tracer.Trace(err)
		}

		// Set hexa context in the gRPC method context, now we can get it in each gRPC method.
		c = context.WithValue(c, ContextKeyHexaCtx, hexaCtx)
	}

	return h(c, req)
}

// NewHexaContextInterceptor returns new instance of the HexaContextInterceptor.
func NewHexaContextInterceptor(cei hexa.ContextExporterImporter) *HexaContextInterceptor {
	return &HexaContextInterceptor{cei: cei}
}

// Ctx gets Hexa context and returns a context to pass to a gROC method.
func Ctx(ctx hexa.Context) context.Context {
	return context.WithValue(context.Background(), ContextKeyHexaCtx, ctx)
}

var _ grpc.UnaryServerInterceptor = (&HexaContextInterceptor{}).UnaryServerInterceptor
var _ grpc.UnaryClientInterceptor = (&HexaContextInterceptor{}).UnaryClientInterceptor
