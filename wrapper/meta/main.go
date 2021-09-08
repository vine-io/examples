package main

import (
	"context"
	"fmt"
	"os"
	"time"

	pb "github.com/vine-io/examples/wrapper/pb"
	"github.com/vine-io/vine"
	"github.com/vine-io/vine/core/client"
	"github.com/vine-io/vine/core/registry"
	"github.com/vine-io/vine/core/server"
	log "github.com/vine-io/vine/lib/logger"
	"github.com/vine-io/vine/util/context/metadata"
)

type hello struct{}

func (h hello) Echo(ctx context.Context, request *pb.Request, response *pb.Response) error {
	response.Result = request.Name
	return nil
}

func main() {
	s := vine.NewService(
		vine.Name("helloworld"),
		vine.WrapCall(CallWrapper()),
		vine.WrapHandler(HandlerWrapper()),
	)

	s.Init()

	pb.RegisterHelloHandler(s.Server(), &hello{})

	go func() {
		time.Sleep(time.Second)
		cc := pb.NewHelloService(s.Name(), s.Client())
		rsp, err := cc.Echo(context.TODO(), &pb.Request{""})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(rsp)
		os.Exit(0)
	}()

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

func CallWrapper() client.CallWrapper {
	return func(fn client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			md, _ := metadata.FromContext(ctx)
			// 追加 client=wrapper
			md.Set("client", "wrapper")
			ctx = metadata.NewContext(ctx, md)
			return fn(ctx, node, req, rsp, opts)
		}
	}
}

func HandlerWrapper() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			val, _ := metadata.Get(ctx, "client")
			fmt.Println("client: ", val)
			return fn(ctx, req, rsp)
		}
	}
}
