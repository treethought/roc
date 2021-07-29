package roc

import (
	"log"
	"net/http"
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
	Resolve(request *Request, ch chan (Evaluator))
}

type Evaluator interface {
	Evaluate(request *Request) Representation
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

	// Request the plugin
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

type Dispatcher struct {
	resolvers   []Resolver
	evalutators []Evaluator
}

type Kernel struct {
	Spaces   map[Identifier]Space
	receiver chan (*Request)
	server   http.Server
	client   *plugin.Client
}

func NewKernel() *Kernel {
	k := &Kernel{
		Spaces:   make(map[Identifier]Space),
		receiver: make(chan *Request),
	}

	return k
}

func (k *Kernel) Register(space Space, endpointPath string) {
	log.Printf("registering endpoint to space: %s", space.Identifier)
	endpoint := NewPhysicalEndpoint(endpointPath)

	if k.Spaces == nil {
		k.Spaces = make(map[Identifier]Space)
	}
	space, ok := k.Spaces[space.Identifier]
	if !ok {
		k.Spaces[space.Identifier] = space
	}

	space.Bind(*endpoint)
}

func (k *Kernel) Receiver() chan (*Request) {
	return k.receiver
}

func (k Kernel) startReceiver() {
	for {
		incoming := <-k.receiver
		k.Dispatch(incoming)
	}
}

func (k Kernel) buildResolveRequest(request *Request) *Request {
	return NewRequest(request.Identifier(), Resolve, nil)

}

func (k Kernel) resolveEndpoint(request *Request) PhysicalEndpoint {
	c := make(chan (PhysicalEndpoint))
	for _, s := range k.Spaces {
		go s.Resolve(request, c)
	}

	return <-c
}

func (k Kernel) Dispatch(request *Request) Representation {
	log.Printf("dispatching request for identifer: %s", request.Identifier())

	endpoint := k.resolveEndpoint(request)
	log.Printf("resolved to endpoint: %s", endpoint.Impl.New)

	// TODO route verbs to methods
	rep := endpoint.Impl.Source(request)
	return rep

}
