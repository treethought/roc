package roc

import (
	"fmt"

	"github.com/google/uuid"
	proto "github.com/treethought/roc/proto/v1"
)

type Dispatcher interface {
	Dispatch(ctx *RequestContext) (Representation, error)
}

type CoreDispatcher struct{}

func NewCoreDispatcher() *CoreDispatcher {
	return &CoreDispatcher{}
}

func (d CoreDispatcher) resolveEndpoint(ctx *RequestContext) (*proto.EndpointMeta, *proto.Space) {
	log.Debug("resolving request", "identifier", ctx.Request().Identifier().String())

	c := make(chan (*proto.EndpointMeta))
	sp := make(chan *proto.Space)

	go func() {
		for _, s := range ctx.m.Scope.Spaces {
			log.Trace("checking space: ", "space", s.Identifier)
			ed, ok := resolveToEndpoint(s, ctx)
			if ok {
				c <- ed
				sp <- s
			}
		}
	}()
	ed := <-c
	space := <-sp

	return ed, space

}

func newTransientSpace(endpoints ...*proto.EndpointMeta) *proto.Space {
	uid := uuid.New()
	spaceID := fmt.Sprintf("dynamic-space://%s", uid.String())

	id := NewIdentifier(spaceID)
	space := NewSpace(id, endpoints...)
	return space
}

func injectParsedArgs(ctx *RequestContext, e *proto.EndpointMeta) {
	args := parseGrammar(e.Grammar, ctx.Request().Identifier().String())
	log.Info("injecting parsed grammar args", "args", args)

	for k, v := range args {
		// TODO parse for identifier better
		// if strings.Contains(v[0], ":/") {
		// 	log.Warn("parsed arg is already identifier", "arg", k, "val", v[0])
		// 	ctx.Request().SetArgument(k, NewIdentifier(v[0]))
		// 	continue
		// }

		// TODO: not overwriting arguments already added to the request
		// might want to change this
		_, exists := ctx.Request().m.ArgumentValues[k]
		if !exists {
			log.Trace("injecting grammar argument ", "arg", k, "val", v[0])
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

func getMeta(ctx *RequestContext) *proto.EndpointMeta {
	log.Info("getting endpoint meta", "endpoint", ctx.Request().Identifier())
	id := ctx.Request().Identifier().String()
	for _, s := range ctx.Scope().Spaces {
		for _, e := range s.Endpoints {
			if e.Identifier == id {
				return e
			}
		}
	}
	return nil
}

func (d CoreDispatcher) Dispatch(ctx *RequestContext) (Representation, error) {
	log.Info("dispatching request",
		"identifier", ctx.Request().Identifier().String(),
		"scope_size", len(ctx.m.Scope.Spaces),
		"verb", ctx.Request().m.Verb,
		"arguments", ctx.Request().m.Arguments,
	)

	if ctx.Request().Verb() == proto.Verb_VERB_META {
		meta := getMeta(ctx)
		return NewRepresentation(meta), nil
	}

	ed, space := d.resolveEndpoint(ctx)
	log.Info("resolved to endpoint", "endpoint", ed.Identifier, "type", ed.Type)
	log.Trace(fmt.Sprintf("%+v", ed))

	injectParsedArgs(ctx, ed)

	ctx.injectValueSpace(ctx.Request())

	var endpoint Endpoint

	switch ed.Type {
	case EndpointTypeTransient:
		endpoint = NewTransientEndpoint(ed)

	case EndpointTypeAccessor:

		existing, ok := space.Clients[ed.Identifier]
		if !ok {
			log.Warn("no existing client for accessor, creating", "endpoint", ed.Identifier)
		}

		phys := NewPhysicalEndpoint(ed, existing)

		endpoint = phys

	case EndpointTypeFileset:
		endpoint = NewFilesetRegex(ed)

	case EndpointTypeTransparentOverlay:
		overlay := NewTransparentOverlay(ed)
		endpoint = overlay

	case OverlayTypeHTTPBridge:
		overlay := NewHTTPBridgeOverlay(ed)
		endpoint = overlay

	case EndpointTypeHTTPRequestAccessor:
		endpoint = NewHttpRequestEndpoint(ed)

	default:
		log.Error("Unknown endpoint type", "endpoint", ed)
		return NewRepresentation(nil), fmt.Errorf("unknown endpoint type")
	}

	log.Trace("evaluating request",
		"identifier", ctx.Request().Identifier(),
		"verb", ctx.Request().m.Verb,
		"ed_type", ed.Type,
	)

	rep := Evaluate(ctx, endpoint)

	repr := NewRepresentation(rep)

	// TODO route verbs to methods
	// rep := endpoint.Source(ctx)
	log.Info("dispatch received response",
		"identifier", ctx.Request().Identifier().String(),
		"representation", repr.Type(),
	)
	log.Info(repr.String())
	return repr, nil
}
