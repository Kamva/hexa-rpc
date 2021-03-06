package hello

import (
	"context"
	"errors"
	"fmt"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	hrpc "github.com/kamva/hexa-rpc"
	"github.com/kamva/tracer"
	"net/http"
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

func (s *helloService) SayHello(c context.Context, m *Message) (*Message, error) {
	ctx := s.ctx(c)

	msg := fmt.Sprintf("Hello %s without hexa context :)", m.Val)
	if ctx != nil {
		fmt.Println("correlation id: ", ctx.CorrelationID())
		msg = fmt.Sprintf("Hello %s with hexa context and correlation id: %s", m.Val, ctx.CorrelationID())
	}

	return &Message{Val: msg}, nil
}

func (s *helloService) SayHelloWithErr(c context.Context, m *Message) (*Message, error) {
	if m.Val == "john" {
		return nil, tracer.Trace(errors.New("name must be john"))
	}
	if m.Val == "panic" {
		panic("name must be john from panic")
	}
	data := hexa.Map{
		"a": "b",
		"c": "d",
	}
	err := hexa.NewLocalizedError(http.StatusNotFound, "rpc.example.code", "localized message :)", errors.New("example error"))
	err = err.SetData(data)
	return nil, err
}

func New() HelloServer {
	return &helloService{}
}

var _ HelloServer = &helloService{}
