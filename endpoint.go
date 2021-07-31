package roc

import (
	"fmt"
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

// Here is an implementation that talks over RPC
type EndpointRPC struct {
	client *rpc.Client
    Dummy string

}

// CanResolve responds affirmatively if the endpoint can handle the request based on the identifier
func (e *EndpointRPC) CanResolve(ctx *RequestContext) bool {
	var resp bool
	err := e.client.Call("Plugin.CanResolve", ctx, &resp)
	if err != nil {
        fmt.Println(err)
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Grammer returns the defined set of identifiers that bind an endpoint to a Space
// Grammar() Grammar
// Evaluate psses a request to create or return a Representation of the requested resource
func (e *EndpointRPC) Evaluate(ctx *RequestContext) Representation {
	var resp Representation
	err := e.client.Call("Plugin.Evaluate", ctx, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Source retrieves representation of resource
func (e *EndpointRPC) Source(ctx *RequestContext) Representation {
	var resp Representation
	err := e.client.Call("Plugin.Source", ctx, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp

}

// Sink updates resource to reflect representation
func (e *EndpointRPC) Sink(ctx *RequestContext) {
	var resp interface{}
	err := e.client.Call("Plugin.Sink", ctx, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}
}

// New creates a resource and return identifier for created resource
// If primary representation is included, use it to initialize resource state
func (e *EndpointRPC) New(ctx *RequestContext) Identifier {
	var resp Identifier
	err := e.client.Call("Plugin.New", ctx, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Delete remove the resource from the space that currently contains it
func (e *EndpointRPC) Delete(ctx *RequestContext) bool {
	var resp bool
	err := e.client.Call("Plugin.Delete", ctx, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Exists tests to see if resource can be resolved and exists
func (e *EndpointRPC) Exists(ctx *RequestContext) bool {
	var resp bool
	err := e.client.Call("Plugin.Exists", ctx, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// func (e *EndpointRPC) Type() string {
// 	var resp string
// 	err := e.client.Call("Plugin.Type", new(interface{}), &resp)
// 	if err != nil {
// 		// You usually want your interfaces to return errors. If they don't,
// 		// there isn't much other choice here.
// 		panic(err)
// 	}

// 	return resp
// }

// Here is the RPC server that EndpointRPC talks to, conforming to
// the requirements of net/rpc
type EndpointRPCServer struct {
	// This is the real implementation
	Impl Endpoint
}

func (s *EndpointRPCServer) CanResolve(ctx *RequestContext, resp *bool) error {
	*resp = s.Impl.CanResolve(ctx)
	return nil
}
func (s *EndpointRPCServer) Evaluate(ctx *RequestContext, resp *Representation) error {
	*resp = s.Impl.Evaluate(ctx)
	return nil
}

func (s *EndpointRPCServer) Source(ctx *RequestContext, resp *Representation) error {
	*resp = s.Impl.Source(ctx)
	return nil
}
func (s *EndpointRPCServer) Sink(ctx *RequestContext, resp *interface{}) error {
	s.Impl.Sink(ctx)
	*resp = nil
	return nil
}
func (s *EndpointRPCServer) New(ctx *RequestContext, resp *Identifier) error {
	*resp = s.Impl.New(ctx)
	return nil
}
func (s *EndpointRPCServer) Delete(ctx *RequestContext, resp *bool) error {
	*resp = s.Impl.Delete(ctx)
	return nil
}
func (s *EndpointRPCServer) Exists(ctx *RequestContext, resp *bool) error {
	*resp = s.Impl.Exists(ctx)
	return nil
}
// func (s *EndpointRPCServer) Type(resp *string) error {
// 	*resp = s.Impl.Type()
// 	return nil
// }

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
type EndpointPlugin struct {
	// Impl Injection
	Impl Endpoint
}

func (p *EndpointPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &EndpointRPCServer{Impl: p.Impl}, nil
}

func (EndpointPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &EndpointRPC{client: c}, nil
}
