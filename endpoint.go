// Package shared contains shared data between the host and plugins.
package roc

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/treethought/roc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Endpoint represents the gateway between a logical resource and the computation
type Endpoint interface {
	Resource

	// Grammer returns the defined set of identifiers that bind an endpoint to a Space
	// Grammar() Grammar

	// Type() string
	// Meta(ctx RequestContextArgument) map[string][]string
}

// Evaluator can be implemented by an endpoint to overide the default request evaluation switch
type Evaluator interface {
	// Evaluate processes a request to create or return a Representation of the requested resource
	Evaluate(ctx *RequestContext) Representation
}

func Evaluate(ctx *RequestContext, e Endpoint) Representation {

	// defer to endpoint's custom implementation if defined
	defined, ok := e.(Evaluator)
	if ok {
		return defined.Evaluate(ctx)
	}

	log.Debug("using default evaluate handler")

	// use default verb routing
	switch ctx.Request.Verb {
	case Source:
		return e.Source(ctx)
	case Sink:
		e.Sink(ctx)
		return nil
	case New:
		return e.New(ctx)
	case Delete:
		return e.Delete(ctx)
	case Exists:
		return e.Exists(ctx)

	default:
		return e.Source(ctx)

	}

}

type BaseEndpoint struct{}

func (e BaseEndpoint) Source(ctx *RequestContext) Representation {
	log.Info("using default source handler")
	return nil
}

func (e BaseEndpoint) Sink(ctx *RequestContext) {}

func (e BaseEndpoint) New(ctx *RequestContext) Identifier {
	return ""
}
func (e BaseEndpoint) Delete(ctx *RequestContext) bool {
	return false
}
func (e BaseEndpoint) Exists(ctx *RequestContext) bool {
	return false
}
func (e BaseEndpoint) Transrept(ctx *RequestContext) Representation {
	return nil
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
		Impl:   e.Impl,
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
