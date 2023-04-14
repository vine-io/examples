// Code generated by proto-gen-vine. DO NOT EDIT.
// source: examples/api/proto/hello.proto

package hello

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	client "github.com/vine-io/vine/core/client"
	server "github.com/vine-io/vine/core/server"
	api "github.com/vine-io/vine/lib/api"
	openapi "github.com/vine-io/vine/lib/api/handler/openapi"
	openapipb "github.com/vine-io/vine/lib/api/handler/openapi/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// API Endpoints for Hello service
func NewHelloEndpoints() []api.Endpoint {
	return []api.Endpoint{
		api.Endpoint{
			Name:        "Hello.Echo",
			Description: "Hello.Echo",
			Path:        []string{"/api/v1/echo/{name}"},
			Method:      []string{"GET"},
			Body:        "*",
			Handler:     "rpc",
		},
	}
}

// Swagger OpenAPI 3.0 for Hello service
func NewHelloOpenAPI() *openapipb.OpenAPI {
	return &openapipb.OpenAPI{
		Openapi: "3.0.1",
		Info: &openapipb.OpenAPIInfo{
			Title:       "HelloService",
			Description: "OpenAPI3.0 for Hello",
			Version:     "v1.0.0",
		},
		Servers: []*openapipb.OpenAPIServer{},
		Tags: []*openapipb.OpenAPITag{
			&openapipb.OpenAPITag{
				Name:        "Hello",
				Description: "OpenAPI3.0 for Hello",
			},
		},
		Paths: map[string]*openapipb.OpenAPIPath{
			"/api/v1/echo/{name}": &openapipb.OpenAPIPath{
				Get: &openapipb.OpenAPIPathDocs{
					Tags:        []string{"Hello"},
					Description: "Hello Echo",
					OperationId: "HelloEcho",
					Parameters: []*openapipb.PathParameters{
						&openapipb.PathParameters{
							Name:        "name",
							In:          "path",
							Description: "EchoReq field name",
							Required:    true,
							Explode:     true,
							Schema: &openapipb.Schema{
								Type: "string",
							},
						},
					},
					Responses: map[string]*openapipb.PathResponse{
						"200": &openapipb.PathResponse{
							Description: "successful response (stream response)",
							Content: &openapipb.PathRequestBodyContent{
								ApplicationJson: &openapipb.ApplicationContent{
									Schema: &openapipb.Schema{Ref: "#/components/schemas/examples.api.proto.EchoRsp"},
								},
							},
						},
					},
					Security: []*openapipb.PathSecurity{},
				},
			},
		},
		Components: &openapipb.OpenAPIComponents{
			SecuritySchemes: &openapipb.SecuritySchemes{},
			Schemas: map[string]*openapipb.Model{
				"examples.api.proto.EchoReq": &openapipb.Model{
					Type: "object",
					Properties: map[string]*openapipb.Schema{
						"name": &openapipb.Schema{
							Type: "string",
						},
					},
				},
				"examples.api.proto.EchoRsp": &openapipb.Model{
					Type: "object",
					Properties: map[string]*openapipb.Schema{
						"reply": &openapipb.Schema{
							Type: "string",
						},
					},
				},
			},
		},
	}
}

// Client API for Hello service
// +gen:openapi
type HelloService interface {
	// +gen:get=/api/v1/echo/{name}
	Echo(ctx context.Context, in *EchoReq, opts ...client.CallOption) (*EchoRsp, error)
}

type helloService struct {
	c    client.Client
	name string
}

func NewHelloService(name string, c client.Client) HelloService {
	return &helloService{
		c:    c,
		name: name,
	}
}

func (c *helloService) Echo(ctx context.Context, in *EchoReq, opts ...client.CallOption) (*EchoRsp, error) {
	req := c.c.NewRequest(c.name, "Hello.Echo", in)
	out := new(EchoRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Hello service
// +gen:openapi
type HelloHandler interface {
	// +gen:get=/api/v1/echo/{name}
	Echo(context.Context, *EchoReq, *EchoRsp) error
}

func RegisterHelloHandler(s server.Server, hdlr HelloHandler, opts ...server.HandlerOption) error {
	type helloImpl interface {
		Echo(ctx context.Context, in *EchoReq, out *EchoRsp) error
	}
	type Hello struct {
		helloImpl
	}
	h := &helloHandler{hdlr}
	endpoints := NewHelloEndpoints()
	for _, ep := range endpoints {
		opts = append(opts, api.WithEndpoint(&ep))
	}
	openapi.RegisterOpenAPIDoc(NewHelloOpenAPI())
	openapi.InjectEndpoints(endpoints...)
	return s.Handle(s.NewHandler(&Hello{h}, opts...))
}

type helloHandler struct {
	HelloHandler
}

func (h *helloHandler) Echo(ctx context.Context, in *EchoReq, out *EchoRsp) error {
	return h.HelloHandler.Echo(ctx, in, out)
}
