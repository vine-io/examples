package handler

import (
	"context"

	hello "github.com/vine-io/examples/api/proto"
)

var _ hello.HelloHandler = (*HelloWorld)(nil)

type HelloWorld struct {
}

func (h HelloWorld) Echo(ctx context.Context, req *hello.EchoReq, rsp *hello.EchoRsp) error {
	rsp.Reply = "reply: " + req.Name
	return nil
}
