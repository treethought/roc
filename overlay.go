package roc

import (
	"fmt"
)

var EndpointTypeTransparentOverlay = "transparentOverlay"

type Overlay interface {
	Endpoint
}

type TransparentOverlay struct {
	BaseEndpoint
	Space      Space
	onRequest  func(ctx *RequestContext)
	onResponse func(ctx *RequestContext, resp Representation)
	Dispatcher Dispatcher
}

func NewTransparentOverlay(ed EndpointDefinition) TransparentOverlay {
	return TransparentOverlay{
		BaseEndpoint: BaseEndpoint{},
		Space:        ed.Space,
		onRequest:    func(ctx *RequestContext) {},
		onResponse:   func(ctx *RequestContext, resp Representation) {},
		Dispatcher:   NewCoreDispatcher(),
	}
}

func (e TransparentOverlay) Type() string {
	return EndpointTypeTransparentOverlay
}

func (o TransparentOverlay) sourceURI(ctx *RequestContext) (string, error) {
	uriRefs, ok := ctx.Request.Arguments["uri"]
	if !ok {
		return "", fmt.Errorf("uri argument not in context")
	}

	refIdentifier := Identifier(uriRefs[0])
	log.Debug("obtained reference arg", "ref", refIdentifier)

	uri, err := ctx.Source(refIdentifier, nil)
	if err != nil {
		log.Error("failed to source uri argument")
		return "", err
	}
	return fmt.Sprint(uri), nil
}

func (o TransparentOverlay) Evaluate(ctx *RequestContext) Representation {
	// transparent hook, cannot modify response
	log.Warn("-------------------------------------")

	o.onRequest(ctx)

	log.Info("overlay handling request", "identifier", ctx.Request.Identifier)

	uri, err := o.sourceURI(ctx)
	if err != nil {
		return err
	}

	// reformat the identifier for context of wrapped space
	// build new res:// scheme from overlay prefix's root
	// i.e. res://my-app/helloworld -> uri=/helloworld -> res://helloworld
	id := Identifier(fmt.Sprintf("res:/%s", uri))

	log.Info("issuing request to wrapped space", "identifier", id)

	log.Debug("initial scope", "size", len(ctx.Scope.Spaces))

	// inject the wrappes space into the request scope
	// we don't create a new request, because this is transparent.
	// we just issue requests into our wrapped space which is otherwise
	// unavailable to outside of the overlay

	// we also replace the existing scope completely, to prevent resolving to this overlay in a loop
	log.Debug("injecting wrapped space", "space", o.Space.Identifier, "size", len(o.Space.EndpointDefinitions))
	ctx.InjectSpace(o.Space)

	ctx.Request = NewRequest(id, ctx.Request.Verb, ctx.Request.RepresentationClass)

	log.Debug("new scope", "spaces", len(ctx.Scope.Spaces), "size", len(ctx.Scope.Spaces[0].EndpointDefinitions))

	// forward the request into our wrapped space
	resp, err := o.Dispatcher.Dispatch(ctx)
	if err != nil {
		panic(err)
	}

	// transparent hook, cannot modify response
	o.onResponse(ctx, resp)
	log.Warn("-------------------------------------")

	return resp
}
