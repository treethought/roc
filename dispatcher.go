package roc

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Dispatcher interface {
	Dispatch(ctx *RequestContext) (Representation, error)
	// Space() []Space
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
type DispatcherPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl Dispatcher
}

func (p *DispatcherPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &DispatcherRPCServer{Impl: p.Impl}, nil
}

func (DispatcherPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &DispatcherRPC{client: c}, nil
}
