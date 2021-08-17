package roc

import (
	"fmt"
	"net/rpc"
)

// EndpointRPC represents the RPC client implementation of an Endpoint
// All Endpoint methods are executed by making an RPC call to the plugin server
type EndpointRPC struct {
	client *rpc.Client
}

// CanResolve responds affirmatively if the endpoint can handle the request based on the identifier
func (e *EndpointRPC) CanResolve(ctx *RequestContext) bool {
	// ctx.Dispatcher = NewPhysicalDispatcher()
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

// Source retrieves representation of resource
func (e *EndpointRPC) Source(ctx *RequestContext) Representation {
	// ctx.Dispatcher = NewPhysicalDispatcher()
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
	// ctx.Dispatcher = NewPhysicalDispatcher()
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
	// ctx.Dispatcher = NewPhysicalDispatcher()
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
	// ctx.Dispatcher = NewPhysicalDispatcher()
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
	// ctx.Dispatcher = NewPhysicalDispatcher()
	var resp bool
	err := e.client.Call("Plugin.Exists", ctx, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// EndpointRPCServer is the server side implementation of Endpoint
// This server receives calls from EndpointRPC and provides the result
// by calling the implemtnation defined in the custom plugin code
type EndpointRPCServer struct {
	// This is the real implementation
	Impl Endpoint
}

func (s *EndpointRPCServer) Source(ctx *RequestContext, resp *Representation) error {
	res := s.Impl.Source(ctx)
	*resp = NewRepresentation(res)
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
