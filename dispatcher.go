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
	log.Info("resolving request", "identifier", ctx.Request.Identifier)

	c := make(chan (EndpointDefinition))
	for _, s := range ctx.Scope.Spaces {
		log.Debug("checking space: ", "space", s.Identifier)
		go s.Resolve(ctx, c)
	}

	return <-c
}

func injectArguments(ctx *RequestContext, e EndpointDefinition) {
	log.Debug("injecting arguments into request context")
	args := e.Grammar.Parse(ctx.Request.Identifier)

	uid := uuid.New()
	spaceID := fmt.Sprintf("dynamic-space://%s", uid.String())

	refArgs := make(map[string][]string)

	transientDefs := []EndpointDefinition{}
	for k, v := range args {
		refArgs[k] = []string{}

		// TODO better way?
		for _, val := range v {
			log.Info("creating transient argument endpoint", "arg", k, "val", val)
			endpoint := NewTransientEndpoint(val)
			transientDefs = append(transientDefs, endpoint.Definition())

			refArgs[k] = append(refArgs[k], endpoint.Identifier().String())
			log.Debug("set argument refernece", "name", k, "ref", endpoint.Identifier().String())
		}
	}

	dynamicSpace := NewSpace(Identifier(spaceID), transientDefs...)
	ctx.InjectSpace(dynamicSpace)
	ctx.Request.Arguments = refArgs

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
	log.Warn("dispatching request",
		"identifier", ctx.Request.Identifier,
		"scope_size", len(ctx.Scope.Spaces),
		"verb", ctx.Request.Verb,
	)

	ed := d.resolveEndpoint(ctx)
	log.Info("resolved to endpoint", "endpoint", ed.Name, "type", ed.Type())

	injectArguments(ctx, ed)

	var endpoint Endpoint

	switch ed.Type() {
	case EndpointTypeTransient:
		endpoint = NewTransientEndpoint(ed.Literal)

	case EndpointTypeAccessor:
		endpoint = NewPhysicalEndpoint(ed.Cmd)

	case EndpointTypeFileset:
		endpoint = NewFilesetRegex(ed.Regex)

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
		"verb", ctx.Request.Verb,
	)
	rep := Evaluate(ctx, endpoint)

	// TODO route verbs to methods
	// rep := endpoint.Source(ctx)
	log.Debug("returning response from dispatcher",
		"identifier", ctx.Request.Identifier,
		"representation", rep,
	)
	return rep, nil
}
