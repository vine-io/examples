package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/vine-io/plugins/registry/etcd"
	"github.com/vine-io/vine/core/registry"
	log "github.com/vine-io/vine/lib/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"os/signal"
	"time"
)

func main() {
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

	// 创建新的服务
	svc := &registry.Service{
		Name:     "go.vine.test",                                        // 服务名称，唯一值
		Version:  "v1.0.0",                                              // 服务版本
		Metadata: map[string]string{"Content-Type": "application/json"}, // 元数据
		Endpoints: []*registry.Endpoint{ // 服务接口
			{
				Name:     "",
				Request:  nil,
				Response: nil,
				Metadata: nil,
			},
		},
		Nodes: []*registry.Node{ // 该服务下的节点信息, 节点的 ID 必须是不同的，当一个 `Registry` 中多次注册相同服务时，服务的该字段就会合并
			{
				Id:      uuid.New().String(),
				Address: "127.0.0.1:11500",
				//Port:     11500,
				Metadata: map[string]string{"os": "windows"},
			},
		},
		Ttl: 30, // 服务过期时间，单位秒
	}

	// 注册服务
	ctx := context.TODO()
	if err := etcdRegistry.Register(ctx, svc); err != nil {
		log.Fatal(err)
	}

	// 查询所有服务
	list, err := etcdRegistry.ListServices(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("查询所有服务")
	for _, item := range list {
		fmt.Println(item.Name)
	}

	// 查询单个服务
	list, err = etcdRegistry.GetService(ctx, "go.vine.test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("查询单个服务")
	for _, item := range list {
		fmt.Println(item.Name, len(item.Nodes))
	}

	// 启动一个监听器
	watcher, err := etcdRegistry.Watch(ctx, registry.WatchService("go.vine.test"))
	if err != nil {
		log.Fatal(err)
	}
	// 停止监听器
	defer watcher.Stop()

	c := make(chan os.Signal)
	signal.Notify(c)
	go func() {
		for {
			e, err := watcher.Next() // 这里会阻塞，直到返回事假
			if err != nil {
				return
			}
			fmt.Printf("%d: %s, %s\n", e.Timestamp, e.Action, e.Service.Name)
		}
	}()
	<-c

	// 注销服务
	if err := etcdRegistry.Deregister(ctx, svc); err != nil {
		log.Fatal(err)
	}
}
