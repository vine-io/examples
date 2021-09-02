package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vine-io/vine"
	log "github.com/vine-io/vine/lib/logger"

	proto "github.com/vine-io/examples/pubsub/proto"
)

// send events using the publisher
func sendEv(topic string, p vine.Event) {
	t := time.NewTicker(time.Second)

	for _ = range t.C {
		// create new event
		ev := &proto.Event{
			Id:        uuid.New().String(),
			Timestamp: time.Now().Unix(),
			Message:   fmt.Sprintf("Messaging you all day on %s", topic),
		}

		log.Infof("publishing %+v\n", ev)

		// publish an event
		if err := p.Publish(context.Background(), ev); err != nil {
			log.Infof("error publishing: %v", err)
		}
	}
}

func main() {
	// create a service
	service := vine.NewService(
		vine.Name("go.vine.cli.pubsub"),
	)
	// parse command line
	service.Init()

	// create publisher
	pub1 := vine.NewEvent("example.topic.pubsub.1", service.Client())
	pub2 := vine.NewEvent("example.topic.pubsub.2", service.Client())

	// pub to topic 1
	go sendEv("example.topic.pubsub.1", pub1)
	// pub to topic 2
	go sendEv("example.topic.pubsub.2", pub2)

	// block forever
	select {}
}
