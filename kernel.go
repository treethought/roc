package roc

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

type Kernel struct {
	Spaces     map[Identifier]Space
	receiver   chan (*RequestContext)
	Dispatcher Dispatcher
	logger     hclog.Logger
	plugins    map[string]PhysicalEndpoint
}

func NewKernel() *Kernel {
	k := &Kernel{
		Spaces:     make(map[Identifier]Space),
		receiver:   make(chan *RequestContext),
		Dispatcher: NewCoreDispatcher(),
		logger: hclog.New(&hclog.LoggerOptions{
			Level:       LogLevel,
			Output:      os.Stderr,
			JSONFormat:  false,
			Name:        "kernel",
			Color:       hclog.ForceColor,
			DisableTime: true,
		}),
	}

	return k
}

func (k Kernel) startTransport(ed EndpointDefinition) (PhysicalTransport, error) {
	k.logger.Info("creating http transport")
	httpt := NewPhysicalTransport(ed.Cmd)

	log.Debug("initializing transport scope")

	phys, ok := httpt.(PhysicalTransport)
	if !ok {
		k.logger.Error("transport is not physical transport")
		os.Exit(1)
	}

	scope := RequestScope{}
	for _, s := range k.Spaces {
		k.logger.Debug("adding to scope", "space", s.Identifier)
		scope.Spaces = append(scope.Spaces, s)
	}

	initMsg := &InitTransport{Scope: scope}

	err := phys.Init(initMsg)
	if err != nil {
		log.Error("failed to initialize transport scope", "transport", httpt)
		return phys, err
	}
	log.Info("initialized transport")
	return phys, nil
}

func (k *Kernel) Start() error {
	for _, s := range k.Spaces {
		for _, ed := range s.EndpointDefinitions {
			if ed.EndpointType == "transport" {
				client, err := k.startTransport(ed)
				if err != nil {
					log.Error("error starting transport:", "err", err)
					os.Exit(1)
				}
				defer client.Client.Kill()

			}
		}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// defer transport.Client.Kill()
	defer plugin.CleanupClients()
	<-sig

	return nil
}

func (k *Kernel) Dispatch(ctx *RequestContext) (Representation, error) {
	if k.Dispatcher == nil {
		k.Dispatcher = NewCoreDispatcher()
	}
	for _, s := range k.Spaces {
		k.logger.Debug("adding to scope", "space", s.Identifier)
		ctx.Scope.Spaces = append(ctx.Scope.Spaces, s)
	}

	k.logger.Info("dispatching request from kernel",
		"num_spaces", len(ctx.Scope.Spaces),
	)

	return k.Dispatcher.Dispatch(ctx)
}

func (k *Kernel) Register(spaces ...Space) {
	for _, space := range spaces {
		k.logger.Info("registering space",
			"space", space.Identifier,
		)
		k.Spaces[space.Identifier] = space
	}
}

func (k *Kernel) Receiver() chan (*RequestContext) {
	return k.receiver
}

func (k Kernel) buildResolveRequestContext(request *Request) *RequestContext {
	return NewRequestContext(request.Identifier, Resolve)

}
