package roc

import "fmt"

type Dispatcher interface {
	Dispatch(ctx *RequestContext) (Representation, error)
}

type CoreDispatcher struct{}

func NewCoreDispatcher() *CoreDispatcher {
	return &CoreDispatcher{}
}

func (d CoreDispatcher) resolveEndpoint(ctx *RequestContext) EndpointDefinition {
	log.Info("resolving request", "identifier", ctx.Request.Identifier)

	c := make(chan (EndpointDefinition))
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

	ed := d.resolveEndpoint(ctx)
	log.Info("resolved to endpoint", "endpoint", ed, "type", ed.Type())
	var endpoint Endpoint

	switch ed.Type() {
	case EndpointTypeAccessor:
		endpoint = NewPhysicalEndpoint(ed.Cmd)

	default:
		log.Error("Unknown endpoint type", "endpoint", ed)
		return nil, fmt.Errorf("unknown endpoint type")
	}

	phys, ok := endpoint.(PhysicalEndpoint)
	if ok {
		defer phys.Client.Kill()
	}


	log.Info("evaluating request",
		"identifier", ctx.Request.Identifier,
	)
	rep := endpoint.Source(ctx)

	// TODO route verbs to methods
	// rep := endpoint.Source(ctx)
	log.Debug("returning response from dispatcher",
		"identifier", ctx.Request.Identifier,
		"representation", rep,
	)
	return rep, nil
}
