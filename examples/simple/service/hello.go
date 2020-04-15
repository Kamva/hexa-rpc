package service

import (
	"context"
	"fmt"
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
	hrpc "github.com/Kamva/hexa-rpc"
	"github.com/Kamva/hexa-rpc/examples/simple/hello"
)

type helloService struct {
}

func (s *helloService) ctx(c context.Context) hexa.Context {
	hexaCtx := c.Value(hrpc.ContextKeyHexaCtx)
	if gutil.IsNil(hexaCtx) {
		return nil
	}
	return hexaCtx.(hexa.Context)
}

func (s *helloService) SayHello(c context.Context, m *hello.Message) (*hello.Message, error) {
	ctx := s.ctx(c)

	msg := fmt.Sprintf("Hello %s without hexa context :)", m.Val)
	if ctx != nil {
		fmt.Println("correlation id: ", ctx.CorrelationID())
		msg = fmt.Sprintf("Hello %s with hexa context and correlation id: %s", m.Val, ctx.CorrelationID())
	}

	return &hello.Message{Val: msg}, nil
}

func New() hello.HelloServer {
	return &helloService{}
}

var _ hello.HelloServer = &helloService{}
