package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	pb "github.com/vine-io/examples/wrapper/pb"
	"github.com/vine-io/vine"
	"github.com/vine-io/vine/core/client"
	"github.com/vine-io/vine/core/client/grpc"
	"github.com/vine-io/vine/core/registry"
	"github.com/vine-io/vine/core/server"
	log "github.com/vine-io/vine/lib/logger"
	"github.com/vine-io/vine/lib/trace"
	"github.com/vine-io/vine/lib/trace/memory"
	"github.com/vine-io/vine/util/wrapper"
)

type hello struct {}

func (h hello) Echo(ctx context.Context, request *pb.Request, response *pb.Response) error {
	ctx, span := trace.DefaultTracer.Start(ctx, "echo")
	defer trace.DefaultTracer.Finish(span)

	response.Result = request.Name

	ctx = trace.ToContext(ctx, uuid.NewString(), uuid.NewString())

	return nil
}

func main() {
	s := vine.NewService(
		vine.WrapHandler(HandlerWrapper()),
	)

	s.Init()

	pb.RegisterHelloHandler(s.Server(), &hello{})

	go func() {
		time.Sleep(time.Second)
		cli := grpc.NewClient(client.WrapCall(CallWrapper()))
		cli = wrapper.TraceCall(s.Name(), memory.NewTracer(), cli)
		cc := pb.NewHelloService(s.Name(), cli)
		cc.Echo(context.TODO(), &pb.Request{"Client"})
	}()

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

func CallWrapper() client.CallWrapper {
	return func(fn client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			traceID, parentID, ok := trace.FromContext(ctx)
			if ok {
				fmt.Printf("call: tarceID=%s parentID=%s\n", traceID, parentID)
			}
			return fn(ctx, node, req, rsp, opts)
		}
	}
}

func HandlerWrapper() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			traceID, parentID, ok := trace.FromContext(ctx)
			if ok {
				fmt.Printf("handle: tarceID=%s parentID=%s\n", traceID, parentID)
			}
			return fn(ctx, req, rsp)
		}
	}
}
