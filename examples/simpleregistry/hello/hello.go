package hello

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
)

type helloService struct {
}

func (s *helloService) SayHello(c context.Context, m *Message) (*Message, error) {
	msg := fmt.Sprintf("Hello %s without hexa context :)", m.Val)

	if hexa.CtxCorrelationId(c) != "" {
		fmt.Println("correlation id: ", hexa.CtxCorrelationId(c))
		msg = fmt.Sprintf("Hello %s with hexa context and correlation id: %s", m.Val, hexa.CtxCorrelationId(c))
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
