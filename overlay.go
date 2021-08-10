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
}

func NewTransparentOverlay(ed EndpointDefinition) TransparentOverlay {
	return TransparentOverlay{
		BaseEndpoint: BaseEndpoint{},
		Space:        ed.Space,
		onRequest:    func(ctx *RequestContext) {},
		onResponse:   func(ctx *RequestContext, resp Representation) {},
	}
}

func (e TransparentOverlay) Type() string {
	return EndpointTypeTransparentOverlay
}

func (o TransparentOverlay) Evaluate(ctx *RequestContext) Representation {
	// transparent hook, cannot modify response
	log.Warn("-------------------------------------")

	o.onRequest(ctx)

	log.Info("overlay handling request", "identifier", ctx.Request.Identifier)

	uri, err := ctx.GetArgumentValue("uri")
	if err != nil {
		return err
	}

	// reformat the identifier for context of wrapped space
	// build new res:// scheme from overlay prefix's root
	// i.e. res://app/helloworld -> uri=/helloworld -> res://helloworld
	id := Identifier(fmt.Sprintf("res://%s", uri))

	// inject the wrapped space into the request scope and
	// issue request into our wrapped space which is otherwise
	// unavailable to outside of the overlay

	ctx.InjectSpace(o.Space)

	req := ctx.CreateRequest(id)
	req.Verb = ctx.Request.Verb
	req.SetRepresentationClass(ctx.Request.RepresentationClass)

	log.Info("issuing request to wrapped space", "identifier", id)
	resp, err := ctx.IssueRequest(req)
	if err != nil {
		log.Error("failed to issue request into wrapped space", "err", err)
		return err
	}

	// transparent hook, cannot modify response
	o.onResponse(ctx, resp)
	log.Warn("-------------------------------------")

	return resp
}
