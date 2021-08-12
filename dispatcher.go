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
	log.Debug("resolving request", "identifier", ctx.Request().Identifier().String())

	c := make(chan (EndpointDefinition))
	for _, s := range ctx.m.Scope.Spaces {
		log.Debug("checking space: ", "space", s.Identifier)
		wrap := Space{s}
		go wrap.Resolve(ctx, c)
	}
	ed := <-c
	close(c)

	return ed

}

func newTransientSpace(endpoints ...EndpointDefinition) Space {
	uid := uuid.New()
	spaceID := fmt.Sprintf("dynamic-space://%s", uid.String())

	id := NewIdentifier(spaceID)
	space := NewSpace(id, endpoints...)
	return space
}

// https: //stackoverflow.com/questions/23030884/is-there-a-way-to-create-an-instance-of-a-struct-from-a-string

func injectParsedArgs(ctx *RequestContext, e EndpointDefinition) {
	log.Debug("injecting parsed arguments into request context")
	args := e.grammar().Parse(ctx.Request().Identifier())

	for k, v := range args {
		// TODO: not overwriting arguments already added to the request
		// might want to change this
		_, exists := ctx.Request().m.ArgumentValues[k]
		if !exists {
			log.Warn("setting grammar arg", "arg", k, "val", v[0])
			rep := NewRepresentation(v[0])
			ctx.Request().SetArgumentByValue(k, rep)
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
		"identifier", ctx.Request().Identifier().String(),
		"scope_size", len(ctx.m.Scope.Spaces),
		"verb", ctx.Request().m.Verb,
	)

	ed := d.resolveEndpoint(ctx)
	log.Info("resolved to endpoint", "endpoint", ed.Name, "type", ed.Type)
	log.Trace(fmt.Sprintf("%+v", ed))

	if ed.Literal != nil {
		log.Warn("ed literal",
			"val", ed.Literal,
			"type_url", ed.Literal.GetValue().TypeUrl,
			"type", ed.Literal.ProtoReflect().Descriptor().Name(),
			"id", ctx.Request().Identifier)

	}

	injectParsedArgs(ctx, ed)

	ctx.injectValueSpace(ctx.Request())

	var endpoint Endpoint

	switch ed.Type {
	case EndpointTypeTransient:
		endpoint = NewTransientEndpoint(ed.Literal)

	case EndpointTypeAccessor:
		endpoint = NewPhysicalEndpoint(ed.Cmd)

	case EndpointTypeFileset:
		endpoint = NewFilesetRegex(ed.Regex)

	case EndpointTypeTransparentOverlay:
		overlay := NewTransparentOverlay(ed)
		endpoint = overlay

	case OverlayTypeHTTPBridge:
		overlay := NewHTTPBridgeOverlay(ed)
		endpoint = overlay

	case EndpointTypeHTTPRequestAccessor:
		log.Error("CREATING HTTPREQUEST ACCESSOR")
		overlay := NewHttpRequestEndpoint(ed)
		endpoint = overlay

	default:
		log.Error("Unknown endpoint type", "endpoint", ed)
		return NewRepresentation(nil), fmt.Errorf("unknown endpoint type")
	}

	phys, ok := endpoint.(PhysicalEndpoint)
	if ok {
		defer phys.Client.Kill()
	}

	log.Debug("evaluating request",
		"identifier", ctx.Request().Identifier,
		"verb", ctx.Request().m.Verb,
	)
	rep := Evaluate(ctx, endpoint)

	// TODO route verbs to methods
	// rep := endpoint.Source(ctx)
	log.Info("dispatch received response",
		"identifier", ctx.Request().Identifier().String(),
		"representation", rep,
		"type_name", rep.ProtoReflect().Descriptor().Name(),
	)
	return rep, nil
}
