package roc

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

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
			Level:      hclog.Info,
			Output:     os.Stderr,
			JSONFormat: false,
			Name:       "kernel",
			Color:      hclog.ForceColor,
		}),
	}

	return k
}

func (k Kernel) Dispatch(ctx *RequestContext) (Representation, error) {
	for _, s := range k.Spaces {
		k.logger.Debug("adding to scope", "space", s.Identifier())
		ctx.Scope.Spaces = append(ctx.Scope.Spaces, s)
	}
	k.logger.Debug("dispatching request from kernel",
		"num_spaces", len(ctx.Scope.Spaces),
	)

	return DispatchRequest(ctx)
}

func (k *Kernel) Register(space Space) {
	k.logger.Info("registering space",
		"space", space.Identifier(),
	)
	k.Spaces[space.Identifier()] = space

}

func (k *Kernel) Receiver() chan (*RequestContext) {
	return k.receiver
}

func (k Kernel) buildResolveRequestContext(request *Request) *RequestContext {
	return NewRequestContext(request.Identifier, Resolve)

}
