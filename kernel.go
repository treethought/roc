package roc

import (
	"os"

	"github.com/hashicorp/go-hclog"
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
	"endpoint":   &EndpointPlugin{},
	"dispatcher": &DispatchPlugin{},
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

type Kernel struct {
	Spaces     map[Identifier]Space
	receiver   chan (*RequestContext)
	Dispatcher Dispatcher
	logger     hclog.Logger
}

func NewKernel() *Kernel {
	k := &Kernel{
		Spaces:   make(map[Identifier]Space),
		receiver: make(chan *RequestContext),
		logger: hclog.New(&hclog.LoggerOptions{
			Level:      hclog.Trace,
			Output:     os.Stderr,
			JSONFormat: false,
			Name:       "kernel",
		}),
	}

	return k
}

func (k Kernel) Dispatch(ctx *RequestContext) (Representation, error) {
	k.logger.Info("Dispatching request")
	for _, s := range k.Spaces {
		k.logger.Debug("adding to scope", "space", s.Identifier())
		ctx.Scope.Spaces = append(ctx.Scope.Spaces, s)
		ctx.Scope.EndpointClients = append(ctx.Scope.EndpointClients, s.Endpoints...)
	}
	// if k.Dispatcher {
	// 	return nil, fmt.Errorf("dispatcher not set")
	// }
	// ctx.Dispatcher = k.DispatcherClient
	return k.Dispatcher.Dispatch(ctx)
	// k.Dispatcher.= MainDispatcher{
	// 	Spaces: k.Spaces,
	// }
	// return k.Dispatcher.Dispatch(ctx)
}

func (k *Kernel) StartDispatcher() {
	k.logger.Info("starting dispatcher")

	k.Dispatcher = NewPhysicalDispatcher().Impl

	// 	var pluginMap = map[string]plugin.Plugin{
	// 		"dispatcher": &DispatchPlugin{},
	//         "endpoint": &EndpointPlugin{},
	// 	}

	// 	// plugin.Serve(&plugin.ServeConfig{
	// 	// 	HandshakeConfig: handshakeConfig,
	// 	// 	Plugins:         pluginMap,
	// 	// 	// Cmd:             exec.Command("./dispatcher/dispatcher"),
	// 	// 	Logger: hclog.New(&hclog.LoggerOptions{
	// 	// 		Level:      hclog.Trace,
	// 	// 		Output:     os.Stderr,
	// 	// 		JSONFormat: false,
	// 	// 		Name:       "dispatch-server",
	// 	// 	}),
	// 	// })

	// 	// We're a host! Start by launching the plugin process.
	// 	client := plugin.NewClient(&plugin.ClientConfig{
	// 		HandshakeConfig: handshakeConfig,
	// 		Plugins:         pluginMap,
	// 		Cmd:             exec.Command("./dispatcher/dispatcher"),
	// 		// Logger:          logger,
	// 	})
	// 	// // defer client.Kill()

	// 	// // Connect via RPC
	// 	rpcClient, err := client.Client()
	// 	if err != nil {
	// 		log.Fatal(err)

	// 	}

	// 	// // Request the plugin
	// 	raw, err := rpcClient.Dispense("dispatcher")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	// // // We should have a Greeter now! This feels like a normal interface
	// 	// // // implementation but is in fact over an RPC connection.
	// 	k.Dispatcher = raw.(Dispatcher)
}

func (k *Kernel) Register(space Space) {
	k.logger.Info("registering space",
		"space", space.Identifier(),
	)
	k.Spaces[space.Identifier()] = space

	// space.Bind(*endpoint)
}

func (k *Kernel) Receiver() chan (*RequestContext) {
	return k.receiver
}

// func (k Kernel) startReceiver() {
// 	for {
// 		incoming := <-k.receiver
// 		k.Dispatch(incoming)
// 	}
// }

func (k Kernel) buildResolveRequestContext(request *Request) *RequestContext {
	return NewRequestContext(request.Identifier, Resolve)

}
