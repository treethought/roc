package roc

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	proto "github.com/treethought/roc/proto/v1"
)

type Kernel struct {
	Spaces     map[Identifier]*proto.Space
	receiver   chan (*RequestContext)
	Dispatcher Dispatcher
	logger     hclog.Logger
	plugins    map[string]PhysicalEndpoint
}

func NewKernel() *Kernel {
	k := &Kernel{
		Spaces:     make(map[Identifier]*proto.Space),
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

func (k Kernel) startTransport(ed *proto.EndpointMeta) (PhysicalTransport, error) {
	k.logger.Info("creating http transport")
	httpt := NewPhysicalTransport(ed.Cmd)

	log.Debug("initializing transport scope")

	phys, ok := httpt.(PhysicalTransport)
	if !ok {
		k.logger.Error("transport is not physical transport")
		os.Exit(1)
	}

	scope := &proto.RequestScope{}
	for _, s := range k.Spaces {
		k.logger.Debug("adding to scope", "space", s.GetIdentifier())
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
	log.Info("starting kernel")


	for _, s := range k.Spaces {
		startAccessors(s)
		for _, ed := range s.GetEndpoints() {
			if ed.Type == "transport" {
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
		k.logger.Debug("adding to scope", "space", s.GetIdentifier())
		ctx.m.Scope.Spaces = append(ctx.m.Scope.Spaces, s)
	}

	k.logger.Info("dispatching request from kernel",
		"num_spaces", len(ctx.m.Scope.Spaces),
	)

	return k.Dispatcher.Dispatch(ctx)
}

func (k *Kernel) Register(spaces ...*proto.Space) {
	for _, s := range spaces {
		k.logger.Info("registering space", "space", s.GetIdentifier())
		id := NewIdentifier(s.GetIdentifier())
		k.Spaces[id] = s
		k.logger.Info("registered spaces", "size", len(spaces))

		for _, e := range s.Endpoints {
			log.Debug("loaded endpoint", "identifier", e.GetIdentifier(), "space", s.GetIdentifier())
		}
	}
}
