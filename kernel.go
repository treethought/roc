package roc

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/treethought/roc/proto"
	"google.golang.org/grpc"
)

type Kernel struct {
	Spaces     map[Identifier]Space
	receiver   chan (*RequestContext)
	Dispatcher Dispatcher
	logger     hclog.Logger
	plugins    map[string]PhysicalEndpoint

	// TODO maybe remove this. ephermeral server created
	// for endpoint grpc client works
	// this is just trying to get transport working
	broker         *plugin.GRPCBroker
	dispatchServer uint32
}

func NewKernel() *Kernel {
	k := &Kernel{
		Spaces:     make(map[Identifier]Space),
		receiver:   make(chan *RequestContext),
		broker:     &plugin.GRPCBroker{},
		Dispatcher: &CoreDispatcher{},
		logger: hclog.New(&hclog.LoggerOptions{
			Level:       hclog.Info,
			Output:      os.Stderr,
			JSONFormat:  false,
			Name:        "kernel",
			Color:       hclog.ForceColor,
			DisableTime: true,
		}),
	}

	return k
}

type CoreDispatcher struct {
}

func (d CoreDispatcher) resolveEndpoint(ctx *RequestContext) Endpoint {
	log.Info("resolving request", "identifier", ctx.Request.Identifier)

	c := make(chan (Endpoint))
	for _, s := range ctx.Scope.Spaces {
		log.Info("checking space: ", "space", s.Identifier)
		go s.Resolve(ctx, c)
	}

	return <-c
}

func (d CoreDispatcher) Dispatch(ctx *RequestContext) (Representation, error) {
	log.Warn("receivied disptach call",
		"identifier", ctx.Request.Identifier,
		"scope_size", len(ctx.Scope.Spaces),
	)

	endpoint := d.resolveEndpoint(ctx)
	log.Info("resolved to endpoint")
	phys, ok := endpoint.(PhysicalEndpoint)
	if !ok {
		return nil, fmt.Errorf("resolved to non-physical endpoint")
	}

	defer phys.Client.Kill()

	log.Info("evaluating request",
		"identifier", ctx.Request.Identifier,
	)
	// TODO route verbs to methods
	rep := endpoint.Source(ctx)
	log.Warn("returning response from dispatcher",
		"identifier", ctx.Request.Identifier,
		"representation", rep,
	)
	return rep, nil
}

func (k *Kernel) startDisptcher() {
	dispatchServer := &DispatcherGRPCServer{Impl: k.Dispatcher}

	var s *grpc.Server

	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterDispatcherServer(s, dispatchServer)
		return s
	}

	brokerID := k.broker.NextId()
	log.Info("starting kernel core dispatcher", "broker", brokerID)
	go k.broker.AcceptAndServe(brokerID, serverFunc)
	k.dispatchServer = brokerID

}

func (k Kernel) startTransport() (PhysicalTransport, error) {
	k.logger.Info("creating http transport")
	httpt := NewPhysicalTransport("./bin/std/transport")

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

	initMsg := &InitTransport{Scope: scope, Dispatcher: k.Dispatcher}

	err := phys.Init(initMsg)
	if err != nil {
		log.Error("failed to initialize transport scope", "transport", httpt)
		return phys, err
	}
	log.Info("initialized transport")
	return phys, nil
}

func (k *Kernel) Start() error {
	// k.startDisptcher()
	transport, err := k.startTransport()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	defer transport.Client.Kill()
	defer plugin.CleanupClients()
	<-sig

	return nil
}

func (k *Kernel) Dispatch(ctx *RequestContext) (Representation, error) {
	if k.Dispatcher == nil {
		k.Dispatcher = &CoreDispatcher{}
	}
	for _, s := range k.Spaces {
		k.logger.Debug("adding to scope", "space", s.Identifier)
		ctx.Scope.Spaces = append(ctx.Scope.Spaces, s)
	}

	k.logger.Info("dispatching request from kernel",
		"num_spaces", len(ctx.Scope.Spaces),
	)

	ctx.Dispatcher = k.Dispatcher
	return k.Dispatcher.Dispatch(ctx)

	// return DispatchRequest(ctx)
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
