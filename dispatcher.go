package roc

import "fmt"

type Dispatcher interface {
	Dispatch(ctx *RequestContext) (Representation, error)
}

type CoreDispatcher struct{}

func NewCoreDispatcher() *CoreDispatcher {
	return &CoreDispatcher{}
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
