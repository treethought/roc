package roc

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Endpoint represents the gateway between a logical resource and the computation
type Endpoint interface {
	Resource

	// CanResolve responds affirmatively if the endpoint can handle the request based on the identifier
	CanResolve(request *Request) bool

	// Grammer returns the defined set of identifiers that bind an endpoint to a Space
	// Grammar() Grammar

	// Evaluate processes a request to create or return a Representation of the requested resource
	Evaluate(request *Request) Representation

	// Type() string
	// Meta(ctx RequestArgument) map[string][]string
}

// Here is an implementation that talks over RPC
type EndpointRPC struct {
	client *rpc.Client
}

// CanResolve responds affirmatively if the endpoint can handle the request based on the identifier
func (e *EndpointRPC) CanResolve(request *Request) bool {
	var resp bool
	err := e.client.Call("Plugin.CanResolve", request, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Grammer returns the defined set of identifiers that bind an endpoint to a Space
// Grammar() Grammar
// Evaluate psses a request to create or return a Representation of the requested resource
func (e *EndpointRPC) Evaluate(request *Request) Representation {
	var resp Representation
	err := e.client.Call("Plugin.Evaluate", request, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Source retrieves representation of resource
func (e *EndpointRPC) Source(request *Request) Representation {
	var resp Representation
	err := e.client.Call("Plugin.Source", request, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp

}

// Sink updates resource to reflect representation
func (e *EndpointRPC) Sink(request *Request) {
	var resp interface{}
	err := e.client.Call("Plugin.Sink", request, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}
}

// New creates a resource and return identifier for created resource
// If primary representation is included, use it to initialize resource state
func (e *EndpointRPC) New(request *Request) Identifier {
	var resp Identifier
	err := e.client.Call("Plugin.New", request, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Delete remove the resource from the space that currently contains it
func (e *EndpointRPC) Delete(request *Request) bool {
	var resp bool
	err := e.client.Call("Plugin.Delete", request, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Exists tests to see if resource can be resolved and exists
func (e *EndpointRPC) Exists(request *Request) bool {
	var resp bool
	err := e.client.Call("Plugin.Exists", request, &resp)
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

func (s *EndpointRPCServer) CanResolve(request *Request, resp *bool) error {
	*resp = s.Impl.CanResolve(request)
	return nil
}
func (s *EndpointRPCServer) Evaluate(request *Request, resp *Representation) error {
	*resp = s.Impl.Evaluate(request)
	return nil
}

func (s *EndpointRPCServer) Source(request *Request, resp *Representation) error {
	*resp = s.Impl.Source(request)
	return nil
}
func (s *EndpointRPCServer) Sink(request *Request, resp *interface{}) error {
	s.Impl.Sink(request)
	*resp = nil
	return nil
}
func (s *EndpointRPCServer) New(request *Request, resp *Identifier) error {
	*resp = s.Impl.New(request)
	return nil
}
func (s *EndpointRPCServer) Delete(request *Request, resp *bool) error {
	*resp = s.Impl.Delete(request)
	return nil
}
func (s *EndpointRPCServer) Exists(request *Request, resp *bool) error {
	*resp = s.Impl.Exists(request)
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
