package main

import (
	"examples/api/handler"
	pb "examples/api/proto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/vine-io/plugins/registry/etcd"
	"github.com/vine-io/vine"
	"github.com/vine-io/vine/lib/api/handler/openapi"
	"github.com/vine-io/vine/lib/api/server"
	httpapi "github.com/vine-io/vine/lib/api/server/http"
	log "github.com/vine-io/vine/lib/logger"
	uapi "github.com/vine-io/vine/util/api"
	"github.com/vine-io/vine/util/helper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var (
	Address       = "127.0.0.1:8994"
	APIPath       = "/"
	enableOpenAPI = false
	enableCors    = false
)

func init() {
	pflag.BoolVar(&enableOpenAPI, "enable-openapi", false, "Enable OpenAPI3")
	pflag.StringVar(&Address, "api-address", "127.0.0.1:8994", "The specify for api address")
	pflag.BoolVar(&enableCors, "enable-cors", false, "Enable OpenAPI3")
}

func main() {
	pflag.Parse()
	var opts []server.Option

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
	// 构建新的服务
	app := vine.NewService(
		vine.Name("go.vine.api"),
		vine.ID(uuid.New().String()),
		vine.Version("v1.0.0"),
		vine.Metadata(map[string]string{
			"api-address": Address,
		}),
		vine.Action(func(command *cobra.Command, strings []string) error {
			ctx := command.PersistentFlags()

			enableTLS, _ := ctx.GetBool("enable-tls")
			if enableTLS {
				config, err := helper.TLSConfig(command)
				if err != nil {
					log.Errorf(err.Error())
					return err
				}

				opts = append(opts, server.EnableTLS(true))
				opts = append(opts, server.TLSConfig(config))
			}
			return nil
		}),
		vine.Registry(etcdRegistry),
	)

	// 服务初始化
	if err := app.Init(); err != nil {
		log.Fatalf("init api: %v", err)
	}

	opts = append(opts, server.EnableCORS(true))

	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()

	if enableOpenAPI {
		openapi.RegisterOpenAPI(app.Client(), ginEngine)
	}

	uapi.PrimpHandler(ginEngine, app.Client(), "go.vine")

	api := httpapi.NewServer(Address)

	if err := api.Init(opts...); err != nil {
		log.Fatal(err)
	}
	api.Handle(APIPath, ginEngine)

	// Start API
	if err := api.Start(); err != nil {
		log.Fatal(err)
	}

	// 注册服务
	if err := pb.RegisterHelloHandler(app.Server(), &handler.HelloWorld{}); err != nil {
		log.Fatalf("register HelloWorld api: %v", err)
	}

	// 服务启动
	if err := app.Run(); err != nil {
		log.Fatalf("start api server: %v", err)
	}

	if err := api.Stop(); err != nil {
		log.Fatal(err)
	}
}
