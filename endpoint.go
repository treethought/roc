// Package shared contains shared data between the host and plugins.
package roc

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/treethought/roc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func Serve(e Endpoint) {
	a, ok := e.(Accessor)
	if ok {
		a.Logger.Debug("starting accessor",
			"name", a.Name,
			"identifier", a.Identifier(),
		)
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"endpoint": &EndpointPlugin{Impl: e},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins:         pluginMap,
        GRPCServer: plugin.DefaultGRPCServer,
	})

}

// Endpoint represents the gateway between a logical resource and the computation
type Endpoint interface {
	Resource

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

func (e *EndpointPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
    proto.RegisterEndpointServer(s, &EndpointGRPCServer{
        Impl: e.Impl,
        broker: broker,
    })
	return nil
}

func (p *EndpointPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &EndpointGRPC{
		client: proto.NewEndpointClient(c),
		broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &EndpointPlugin{}
