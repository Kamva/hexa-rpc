package hrpc

import (
	"encoding/json"
	"errors"
	"github.com/Kamva/hexa"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"net/http"
)

const hexaToStatusError = "error on converting Hexa error into Status with message: "
const statusToHexaError = "error on converting gRPC Status into Hexa error with message: "

// Status gets a Hexa error and converts it to gRPC Status
// Implementation Details:
// - Convert http status to gRPC code
// - Set localized message and data.
func Status(hexaErr hexa.Error, t hexa.Translator) *status.Status {
	if hexaErr == nil {
		return nil
	}

	code := CodeFromHTTPStatus(hexaErr.HTTPStatus())
	localMsg, _ := hexaErr.Localize(t)
	data, _ := json.Marshal(hexaErr.Data())

	s := status.New(code, hexaErr.Error())
	s, err := s.WithDetails(&ErrorDetails{
		Status:           int32(hexaErr.HTTPStatus()),
		Code:             hexaErr.Code(),
		LocalizedMessage: localMsg,
		Data:             string(data),
	})
	if err != nil {
		grpclog.Infof(hexaToStatusError, err.Error())
	}
	return s
}

// Error gets a gRPC status and converts it to Hexa error
func Error(status *status.Status) hexa.Error {
	if status == nil {
		return nil
	}
	httpStatus := http.StatusInternalServerError
	code := ""
	localizedMsg := ""
	data := hexa.Map{}
	for _, detail := range status.Details() {
		if d, ok := detail.(*ErrorDetails); ok {
			httpStatus = int(d.Status)
			code = d.Code
			localizedMsg = d.LocalizedMessage
			err := json.Unmarshal([]byte(d.Data), &data)
			if err != nil {
				grpclog.Info(statusToHexaError, err)
			}
		}
	}
	return hexa.NewLocalizedError(httpStatus, code, localizedMsg, errors.New(status.Message())).SetData(data)
}

// HTTPStatusFromCode converts a gRPC error code into the corresponding HTTP response status.
// See: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
// Note: We got this function from the [gRPC gateway](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/runtime/errors.go)
func HTTPStatusFromCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		// Note, this deliberately doesn't translate to the similarly named '412 Precondition Failed' HTTP response status.
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	}

	grpclog.Infof("Unknown gRPC error code: %v", code)
	return http.StatusInternalServerError
}

// CodeFromHTTPStatus converts a https status into corresponding gRPC error code.
// Note: error mapping from http status to hRPC code is not good, do not use this
// function as you can.
func CodeFromHTTPStatus(status int) codes.Code {
	switch status {
	case http.StatusOK:
		return codes.OK
	case http.StatusRequestTimeout:
		return codes.Canceled
	case http.StatusInternalServerError:
		//return codes.Unknown
		return codes.Internal
	case http.StatusBadRequest:
		// Note: this deliberately doesn't translate to
		// return codes.InvalidArgument
		return codes.FailedPrecondition
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	}

	grpclog.Infof("unsupported http status ", status)
	return codes.Unknown
}
