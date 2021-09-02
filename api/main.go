package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vine-io/cli"
	"github.com/vine-io/examples/api/handler"
	pb "github.com/vine-io/examples/api/proto"

	"github.com/vine-io/vine"
	ahandler "github.com/vine-io/vine/lib/api/handler"
	"github.com/vine-io/vine/lib/api/handler/openapi"
	arpc "github.com/vine-io/vine/lib/api/handler/rpc"
	"github.com/vine-io/vine/lib/api/resolver"
	"github.com/vine-io/vine/lib/api/resolver/grpc"
	"github.com/vine-io/vine/lib/api/router"
	regRouter "github.com/vine-io/vine/lib/api/router/registry"
	"github.com/vine-io/vine/lib/api/server"
	httpapi "github.com/vine-io/vine/lib/api/server/http"
	log "github.com/vine-io/vine/lib/logger"
	"github.com/vine-io/vine/util/helper"
	"github.com/vine-io/vine/util/namespace"
)

var (
	Address       = ":8080"
	Handler       = "rpc"
	Type          = "api"
	APIPath       = "/"
	enableOpenAPI = false

	flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "api-address",
			Usage:       "The specify for api address",
			EnvVars:     []string{"VINE_API_ADDRESS"},
			Value:       Address,
			Destination: &Address,
		},
		&cli.BoolFlag{
			Name:    "enable-openapi",
			Usage:   "Enable OpenAPI3",
			EnvVars: []string{"VINE_ENABLE_OPENAPI"},
			Value:   true,
		},
		&cli.BoolFlag{
			Name:    "enable-cors",
			Usage:   "Enable CORS, allowing the API to be called by frontend applications",
			EnvVars: []string{"VINE_API_ENABLE_CORS"},
			Value:   true,
		},
	}
)

func main() {
	// Init API
	var opts []server.Option

	// initialise service
	svc := vine.NewService(
		vine.Name("go.vine.api"),
		vine.Id(uuid.New().String()),
		vine.Version("v1.0.0"),
		vine.Metadata(map[string]string{
			"api-address": Address,
		}),
		vine.Flags(flags...),
		vine.Action(func(ctx *cli.Context) error {
			enableOpenAPI = ctx.Bool("enable-openapi")

			if ctx.Bool("enable-tls") {
				config, err := helper.TLSConfig(ctx)
				if err != nil {
					log.Errorf(err.Error())
					return err
				}

				opts = append(opts, server.EnableTLS(true))
				opts = append(opts, server.TLSConfig(config))
			}
			return nil
		}),
	)

	svc.Init()

	opts = append(opts, server.EnableCORS(true))

	// create the router
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	if enableOpenAPI {
		openapi.RegisterOpenAPI(app)
	}

	// create the namespace resolver
	nsResolver := namespace.NewResolver(Type, "go.vine")
	// resolver options
	ropts := []resolver.Option{
		resolver.WithNamespace(nsResolver.ResolveWithType),
		resolver.WithHandler(Handler),
	}

	log.Infof("Registering API RPC Handler at %s", APIPath)
	rr := grpc.NewResolver(ropts...)
	rt := regRouter.NewRouter(
		router.WithHandler(arpc.Handler),
		router.WithResolver(rr),
		router.WithRegistry(svc.Options().Registry),
	)
	rp := arpc.NewHandler(
		ahandler.WithNamespace("go.vine"),
		ahandler.WithRouter(rt),
		ahandler.WithClient(svc.Client()),
	)
	app.Group(APIPath, rp.Handle)

	api := httpapi.NewServer(Address)

	if err := api.Init(opts...); err != nil {
		log.Fatal(err)
	}
	api.Handle("/", app)

	// Start API
	if err := api.Start(); err != nil {
		log.Fatal(err)
	}

	_ = pb.RegisterHelloHandler(svc.Server(), &handler.HelloWorld{})

	// Run server
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}

	// Stop API
	if err := api.Stop(); err != nil {
		log.Fatal(err)
	}
}
