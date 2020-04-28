package hrpc

import (
	"encoding/json"
	"errors"
	"github.com/Kamva/hexa"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"net/http"
)

const convertErrMsg = "error on converting Hexa error gRPC Status with message: "
const convertStatusMsg = "error on converting gRPC Status into Hexa error with message: "

// Status gets a Hexa error and converts it to gRPC Status
// Implementation Details:
// - Convert http status to gRPC code
// - Set localized message and data.
func Status(hexaErr hexa.Error, t hexa.Translator) *status.Status {
	if hexaErr == nil {
		return nil
	}

	code := CodeFromHTTPStatus(hexaErr.HTTPStatus())
	localMsg, err := hexaErr.Localize(t)
	if err != nil {
		grpclog.Info(convertErrMsg, err.Error())
	}
	msg := &errdetails.LocalizedMessage{Message: localMsg}
	s := status.New(code, hexaErr.Error())

	// Append localized message
	s, err = s.WithDetails(msg)
	if err != nil {
		grpclog.Info(convertErrMsg, err.Error())
	}

	// Append data
	help := &errdetails.Help{Links: linkData(hexaErr)}
	s, err = s.WithDetails(help)
	if err != nil {
		grpclog.Infof(convertErrMsg, err.Error())
		return s
	}
	return s

}

// Error gets a gRPC status and converts it to Hexa error
func Error(status *status.Status) hexa.Error {
	if status == nil {
		return nil
	}
	code := ""
	localizedMsg := ""
	data := hexa.Map{}
	for _, detail := range status.Details() {
		switch t := detail.(type) {
		case *errdetails.LocalizedMessage:
			localizedMsg = t.Message
		case *errdetails.Help:
			data, code = extractLinkData(t)
		}
	}
	return hexa.NewLocalizedError(code, hexa.TranslateKeyEmptyMessage, localizedMsg).
		SetHTTPStatus(HTTPStatusFromCode(status.Code())).
		SetError(errors.New(status.Message())).
		SetData(data)
}

func linkData(e hexa.Error) []*errdetails.Help_Link {
	data, err := json.Marshal(e.Data())
	if err != nil {
		grpclog.Infof(convertErrMsg, err.Error())
	}
	return []*errdetails.Help_Link{
		{
			Description: "code",
			Url:         e.Code(),
		},
		{
			Description: "data",
			Url:         string(data),
		},
	}
}

func extractLinkData(help *errdetails.Help) (data hexa.Map, code string) {
	data = make(hexa.Map)
	for _, link := range help.Links {
		switch link.Description {
		case "code":
			code = link.Url
		case "data":
			err := json.Unmarshal([]byte(link.Url), &data)
			if err != nil {
				grpclog.Info(convertStatusMsg, err)
			}
		}
	}

	return
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
