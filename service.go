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
	return c.Value(ContextKeyHexaCtx).(hexa.Context)
}
