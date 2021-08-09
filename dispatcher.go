package roc

import (
	"fmt"

	"github.com/google/uuid"
)

type Dispatcher interface {
	Dispatch(ctx *RequestContext) (Representation, error)
}

type CoreDispatcher struct{}

func NewCoreDispatcher() *CoreDispatcher {
	return &CoreDispatcher{}
}

func (d CoreDispatcher) resolveEndpoint(ctx *RequestContext) EndpointDefinition {
	log.Debug("resolving request", "identifier", ctx.Request.Identifier)

	c := make(chan (EndpointDefinition))
	for _, s := range ctx.Scope.Spaces {
		log.Debug("checking space: ", "space", s.Identifier)
		go s.Resolve(ctx, c)
	}

	return <-c

}

func newTransientSpace(endpoints ...EndpointDefinition) Space {
	uid := uuid.New()
	spaceID := fmt.Sprintf("dynamic-space://%s", uid.String())
	space := NewSpace(Identifier(spaceID), endpoints...)
	return space

}

func injectParsedArgs(ctx *RequestContext, e EndpointDefinition) {
	log.Debug("injecting parsed arguments into request context")
	args := e.Grammar.Parse(ctx.Request.Identifier)

	for k, v := range args {
		// TODO: not overwriting arguments already added to the request
		// might want to change this
		_, exists := ctx.Request.argumentValues[k]
		if !exists {
			ctx.Request.SetArgumentByValue(k, v[0])
		}
	}

	// TODO

	// pass by reference

	// pass by value (literal)
	// place representation into dynamic generated space
	// with gnerated identifier. inject into scope
	// then give the argument the value of the identifier

	// pass by request (lazy load)
	// nstead of putting representation into dynamic space
	// a dynamic generated request is placed into the space
	// the request is executed if the endpoint sources the argument
	// otherwise, not executed

}

func (d CoreDispatcher) Dispatch(ctx *RequestContext) (Representation, error) {
	log.Info("dispatching request",
		"identifier", ctx.Request.Identifier,
		"scope_size", len(ctx.Scope.Spaces),
		"verb", ctx.Request.Verb,
	)

	ed := d.resolveEndpoint(ctx)
	log.Info("resolved to endpoint", "endpoint", ed.Name, "type", ed.Type())
	log.Trace(fmt.Sprintf("%+v", ed))

	injectParsedArgs(ctx, ed)

	ctx.injectValueSpace(ctx.Request)

	var endpoint Endpoint

	switch ed.Type() {
	case EndpointTypeTransient:
		endpoint = NewTransientEndpoint(ed.Literal)

	case EndpointTypeAccessor:
		endpoint = NewPhysicalEndpoint(ed.Cmd)

	case EndpointTypeFileset:
		endpoint = NewFilesetRegex(ed.Regex)

	case EndpointTypeTransparentOverlay:
		overlay := NewTransparentOverlay(ed)
		endpoint = overlay

	default:
		log.Error("Unknown endpoint type", "endpoint", ed)
		return nil, fmt.Errorf("unknown endpoint type")
	}

	phys, ok := endpoint.(PhysicalEndpoint)
	if ok {
		defer phys.Client.Kill()
	}

	log.Debug("evaluating request",
		"identifier", ctx.Request.Identifier,
		"verb", ctx.Request.Verb,
	)
	rep := Evaluate(ctx, endpoint)

	// TODO route verbs to methods
	// rep := endpoint.Source(ctx)
	log.Info("dispatch received response",
		"identifier", ctx.Request.Identifier,
		"representation", rep,
	)
	return rep, nil
}
