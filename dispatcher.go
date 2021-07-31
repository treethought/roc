package roc

import (
	"fmt"
	"log"
	"net/rpc"
	"os/exec"

	"github.com/hashicorp/go-plugin"
)

type Dispatcher interface {
	Dispatch(ctx *RequestContext) (Representation, error)
	// Space() []Space
}

// Here is an implementation that talks over RPC
type DispatcherClient struct {
	client *rpc.Client
    EndpointClients []Endpoint
}

// CanResolve responds affirmatively if the endpoint can handle the request based on the identifier
func (e *DispatcherClient) Dispatch(ctx *RequestContext) (Representation, error) {
	var resp Representation
	err := e.client.Call("Plugin.Dispatch", ctx, &resp)
	if err != nil {
		fmt.Println(err)
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp, err
}

// Here is the RPC server that DispatchRPC talks to, conforming to
// the requirements of net/rpc
type DispatcherRPCServer struct {
	// This is the real implementation
	Impl Dispatcher
}

func (s *DispatcherRPCServer) Dispatch(ctx *RequestContext, resp *Representation) error {
	rep, err := s.Impl.Dispatch(ctx)
	*resp = rep
	return err

}

// func (s *DispatcherRPCServer) Dispatch(ctx *RequestContext, resp *Representation)  error {
//     *resp, err := s.Impl.Dispatch(ctx)
// 	return resp, err
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
type DispatchPlugin struct {
	// Impl Injection
	Impl Dispatcher
}

func (p *DispatchPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &DispatcherRPCServer{Impl: p.Impl}, nil
}

func (p DispatchPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &DispatcherClient{client: c}, nil
}

type PhysicalDispatcher struct {
	client *plugin.Client
	rpc    plugin.ClientProtocol
	Impl   Dispatcher
}

func NewPhysicalDispatcher() *PhysicalDispatcher {
	dispatcher := &PhysicalDispatcher{}
	// We're a host! Start by launching the plugin pss.
	dispatcher.client = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./dispatcher/dispatcher"),
		// Logger:          logger,
	})

	// Connect via RPC
	rpcClient, err := dispatcher.client.Client()
	if err != nil {
		log.Fatal(err)
	}
	dispatcher.rpc = rpcClient

	// RequestContext the plugin
	raw, err := rpcClient.Dispense("dispatcher")
	if err != nil {
		log.Fatal(err)
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	dispatcher.Impl = raw.(Dispatcher)
	return dispatcher

}
