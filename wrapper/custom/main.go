package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/vine-io/examples/wrapper/pb"
	"github.com/vine-io/vine"
	"github.com/vine-io/vine/core/client"
	"github.com/vine-io/vine/core/registry"
	log "github.com/vine-io/vine/lib/logger"
	"github.com/vine-io/vine/lib/trace"
)

type hello struct {
}

func (h hello) Echo(ctx context.Context, request *pb.Request, response *pb.Response) error {
	response.Result = request.Name
	return nil
}

func main() {
	s := vine.NewService(
		vine.WrapCall(LoggerWrapper(), SubWrapper()),
	)

	s.Init()

	pb.RegisterHelloHandler(s.Server(), &hello{})

	go func() {
		time.Sleep(time.Second)
		cc := pb.NewHelloService(s.Name(), s.Client())
		cc.Echo(context.TODO(), &pb.Request{"Client"})
	}()

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

func LoggerWrapper() client.CallWrapper {
	return func(fn client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {

			fmt.Println("logger wrapper: before call")
			err := fn(ctx, node, req, rsp, opts)
			fmt.Println("logger wrapper: after call")

			return err
		}
	}
}

func SubWrapper() client.CallWrapper {
	return func(fn client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {

			fmt.Println(trace.FromContext(ctx))

			fmt.Println("sub wrapper: before call")
			err := fn(ctx, node, req, rsp, opts)
			fmt.Println("sub wrapper: after call")

			return err
		}
	}
}
