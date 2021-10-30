package hrpc

import (
	"context"

	"github.com/kamva/hexa"
)

// Service is just a base struct to use in hexa services(it's optional and you can drop it).
type Service struct {
}

// Ctx method extract the hexa context from the context of a gRPC service
func (r Service) Ctx(c context.Context) hexa.Context {
	// We ignore error, user can check if ctx is nil or not.
	ctx, _ := hexa.NewContextFromRawContext(c)
	return ctx
}

func (r Service) Error(ctx context.Context, err error) *ErrorDetails {
	return NewErrDetails(HexaErrFromErr(err), r.Ctx(ctx).Translator())
}
