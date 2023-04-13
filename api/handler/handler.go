package handler

import (
	"context"
	hello "examples/api/proto"
)

var _ hello.HelloHandler = (*HelloWorld)(nil)

type HelloWorld struct {
}

func (h *HelloWorld) Echo(ctx context.Context, req *hello.EchoReq, rsp *hello.EchoRsp) error {
	rsp.Reply = "reply: " + req.Name
	return nil
}
