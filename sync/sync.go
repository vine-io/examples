package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/vine-io/plugins/sync/etcd"
	"github.com/vine-io/vine/lib/sync"
)

func rd() int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(10)
}

func main() {
	s := etcd.NewSync()
	if err := s.Init(); err != nil {
		log.Fatalln(err)
	}

	done := make(chan struct{}, 1)

	id := uuid.NewString()
	role1, err := s.Leader(context.Background(), "leader", sync.LeaderNS("default"), sync.LeaderId(id), sync.LeaderTTL(3))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("primary:", id)

	go func() {
		rid := uuid.NewString()
		role2, err := s.Leader(context.Background(), "leader", sync.LeaderNS("default"), sync.LeaderId(rid), sync.LeaderTTL(3))
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("primary:", rid)

		time.Sleep(time.Second * 2)

		if err := role2.Resign(); err != nil {
			log.Fatalln(err)
		}

		fmt.Println(rid + " resign")
		done <- struct{}{}
	}()

	time.Sleep(time.Second * 1)

	ms, _ := s.ListMembers(context.Background(), sync.MemberNS("default"))
	for _, m := range ms {
		fmt.Println(m)
	}

	if err := role1.Resign(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(id + " resign")

	<-done
}
