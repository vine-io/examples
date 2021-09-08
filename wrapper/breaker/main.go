package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sony/gobreaker"
	pb "github.com/vine-io/examples/wrapper/pb"
	bk "github.com/vine-io/plugins/wrapper/breaker/gobreaker"
	"github.com/vine-io/vine"
	log "github.com/vine-io/vine/lib/logger"
)

type hello struct{}

func (h hello) Echo(ctx context.Context, request *pb.Request, response *pb.Response) error {
	response.Result = request.Name
	return nil
}

func main() {
	bs := gobreaker.Settings{
		Name:          "breaker",
		MaxRequests:   10,
		Interval:      time.Second * 10,
		Timeout:       time.Second * 15,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Info("%s %s => %s", name, from.String(), to.String())
		},
	}
	breaker := bk.NewCustomClientWrapper(bs, bk.BreakService)

	s := vine.NewService(
		vine.Name("helloworld"),
		vine.WrapClient(breaker),
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
