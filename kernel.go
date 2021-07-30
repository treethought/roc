package roc

import (
	"log"
	"os/exec"

	"github.com/hashicorp/go-plugin"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"endpoint": &EndpointPlugin{},
}

type Resolver interface {
	Resolve(request *RequestContext, ch chan (Endpoint))
	Identifier() Identifier
	// Bind(PhysicalEndpoint)
}

type Evaluator interface {
	Evaluate(request *RequestContext) Representation
	Identifier() Identifier
}

type PhysicalEndpoint struct {
	client *plugin.Client
	rpc    plugin.ClientProtocol
	Impl   Endpoint
}

func NewPhysicalEndpoint(path string) *PhysicalEndpoint {
	endpoint := &PhysicalEndpoint{}
	// We're a host! Start by launching the plugin process.
	endpoint.client = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(path),
		// Logger:          logger,
	})

	// Connect via RPC
	rpcClient, err := endpoint.client.Client()
	if err != nil {
		log.Fatal(err)
	}
	endpoint.rpc = rpcClient

	// RequestContext the plugin
	raw, err := rpcClient.Dispense("endpoint")
	if err != nil {
		log.Fatal(err)
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	endpoint.Impl = raw.(Endpoint)
	return endpoint

}

func (e *PhysicalEndpoint) Kill() {
	e.client.Kill()
}

// type Dispatcher struct {
// 	resolvers   []Resolver
// 	evalutators []Evaluator
// }

type Kernel struct {
	Spaces   map[Identifier]Resolver
	receiver chan (*RequestContext)
	client   *plugin.Client
}

func NewKernel() *Kernel {
	k := &Kernel{
		Spaces:   make(map[Identifier]Resolver),
		receiver: make(chan *RequestContext),
	}

	return k
}

func (k *Kernel) Register(space Resolver) {
	log.Printf("registering endpoint to space: %s", space.Identifier())
	k.Spaces[space.Identifier()] = space

	// space.Bind(*endpoint)
}

func (k *Kernel) Receiver() chan (*RequestContext) {
	return k.receiver
}

func (k Kernel) startReceiver() {
	for {
		incoming := <-k.receiver
		k.Dispatch(incoming)
	}
}

func (k Kernel) buildResolveRequestContext(request *Request) *RequestContext {
	return NewRequestContext(request.Identifier, Resolve)

}

func (k Kernel) resolveEndpoint(ctx *RequestContext) Endpoint {
	c := make(chan (Endpoint))
	for _, s := range k.Spaces {
		go s.Resolve(ctx, c)
	}

	return <-c
}

func (k Kernel) Dispatch(ctx *RequestContext) (Representation, error) {
	log.Printf("dispatching request for identifer: %s", ctx.Request.Identifier)

	endpoint := k.resolveEndpoint(ctx)

	// phys, ok := endpoint.(PhysicalEndpoint)
	// if !ok {
	//     log.Println("resolved endpoint is not a plugin")
	//     return nil
	// }

	// log.Printf("resolved to endpoint: %s", phys.Impl.New)

	// TODO route verbs to methods
	rep := endpoint.Source(ctx)
	return rep, nil

}
