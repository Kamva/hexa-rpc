package hrpc

import (
	"context"
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorInterceptor implements a gRPC interceptor to
//convert error into status and reverse.
type ErrorInterceptor struct {
}

// NewErrorInterceptor returns new instance of the ErrorInterceptor
func NewErrorInterceptor() *ErrorInterceptor {
	return &ErrorInterceptor{}
}

// UnaryServerInterceptor returns unary server interceptor to convert Hexa error to status.
// Note: error interceptor must be last interceptor in chained interceptor.
func (i ErrorInterceptor) UnaryServerInterceptor(t hexa.Translator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, rErr := handler(ctx, req)

		if rErr == nil {
			return resp, rErr
		}

		// If error implements the GRPCStatus interface, we dont convert it to
		if _, ok := rErr.(interface{ GRPCStatus() *status.Status }); ok {
			return resp, rErr
		}

		baseErr, ok := gutil.CauseErr(rErr).(hexa.Error)
		if !ok {
			baseErr = ErrUnknownError.SetError(rErr)
		}
		return resp, Status(baseErr, t).Err()
	}
}

// UnaryClientInterceptor returns client interceptor to convert status to Hexa error.
// Note: error interceptor must be first client interceptor.
func (i ErrorInterceptor) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil || status.Convert(err).Code() == codes.OK {
			return err
		}
		s := status.Convert(err)
		return Error(s)
	}
}
