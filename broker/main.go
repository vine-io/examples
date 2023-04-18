package main

import (
	"context"
	"fmt"
	"time"

	"github.com/vine-io/vine/core/broker"

	"github.com/vine-io/plugins/registry/etcd"
	"github.com/vine-io/vine/core/broker/http"
	log "github.com/vine-io/vine/lib/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	topic := "go.vine.topic.foo"

	//初始化etcd
	etcdClient, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	etcdRegistry := etcd.NewEtcdRegistry(etcdClient)
	if err := etcdRegistry.Init(); err != nil {
		log.Fatal(err)
	}

	b := http.NewBroker(broker.Registry(etcdRegistry))

	if err := b.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := b.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	go func() {
		// receive message from broker
		b.Subscribe(topic, func(p broker.Event) error {
			fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)
			return nil
		})
	}()

	go func() {
		<-time.After(time.Second * 1)
		// publish message to broker
		b.Publish(context.TODO(), topic, &broker.Message{Header: map[string]string{"a": "b"}, Body: []byte("hello world")})
	}()

	time.Sleep(time.Second * 2)

	b.Disconnect()
}
