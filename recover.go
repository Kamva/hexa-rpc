package hrpc

import (
	"errors"
	"fmt"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
)

// RecoverHandler handle handle recovered data from panics in the gRPC server
func RecoverHandler(r interface{}) error {
	e, ok := r.(error)
	if ok {
		return e
	}
	return errors.New(fmt.Sprint(e))
}

// Assertion
var _ grpc_recovery.RecoveryHandlerFunc = RecoverHandler
