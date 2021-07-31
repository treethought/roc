// Package shared contains shared data between the host and plugins.
package roc

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Endpoint represents the gateway between a logical resource and the computation
type Endpoint interface {
	Resource

	// CanResolve responds affirmatively if the endpoint can handle the request based on the identifier
	CanResolve(ctx *RequestContext) bool

	// Grammer returns the defined set of identifiers that bind an endpoint to a Space
	// Grammar() Grammar

	// Evaluate processes a request to create or return a Representation of the requested resource
	Evaluate(ctx *RequestContext) Representation

	// Type() string
	// Meta(ctx RequestContextArgument) map[string][]string
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a EndpointRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return EndpointRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type EndpointPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl Endpoint
}

func (p *EndpointPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &EndpointRPCServer{Impl: p.Impl}, nil
}

func (EndpointPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &EndpointRPC{client: c}, nil
}

// func (p *CounterPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
// 	proto.RegisterCounterServer(s, &GRPCServer{
// 		Impl:   p.Impl,
// 		broker: broker,
// 	})
// 	return nil
// }

// func (p *CounterPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
// 	return &GRPCClient{
// 		client: proto.NewCounterClient(c),
// 		broker: broker,
// 	}, nil
// }

// var _ plugin.GRPCPlugin = &CounterPlugin{}

// type AddHelper interface {
// 	Sum(int64, int64) (int64, error)
// }

// // KV is the interface that we're exposing as a plugin.
// type Counter interface {
// 	Put(key string, value int64, a AddHelper) error
// 	Get(key string) (int64, error)
// }

// // This is the implementation of plugin.Plugin so we can serve/consume this.
// // We also implement GRPCPlugin so that this plugin can be served over
// // gRPC.
// type CounterPlugin struct {
// 	plugin.NetRPCUnsupportedPlugin
// 	// Concrete implementation, written in Go. This is only used for plugins
// 	// that are written in Go.
// 	Impl Counter
// }
