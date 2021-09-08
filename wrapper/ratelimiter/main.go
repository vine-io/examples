package main

import (
	"context"

	pb "github.com/vine-io/examples/wrapper/pb"
	ub "github.com/vine-io/plugins/wrapper/ratelimiter/uber"
	"github.com/vine-io/vine"
	log "github.com/vine-io/vine/lib/logger"
)

type hello struct{}

func (h hello) Echo(ctx context.Context, request *pb.Request, response *pb.Response) error {
	response.Result = request.Name
	return nil
}

func main() {
	handler := ub.NewHandlerWrapper(1000)

	s := vine.NewService(
		vine.Name("helloworld"),
		vine.WrapHandler(handler),
	)

	s.Init()

	pb.RegisterHelloHandler(s.Server(), &hello{})

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
