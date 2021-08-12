package roc

import (
	"fmt"

	"github.com/treethought/roc/proto"
)

var EndpointTypeTransparentOverlay = "transparentOverlay"

type Overlay interface {
	Endpoint
}

type TransparentOverlay struct {
	BaseEndpoint
	Space      *proto.Space
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

func (o TransparentOverlay) Evaluate(ctx *RequestContext) interface{} {
	// transparent hook, cannot modify response
	log.Warn("-------------------------------------")

	o.onRequest(ctx)

	log.Info("overlay handling request", "identifier", ctx.Request().Identifier().String())

	uri, err := ctx.GetArgumentValue("uri")
	if err != nil {
		log.Error("failed to source uri argument representation")
		return NewRepresentation(&proto.ErrorMessage{Message: err.Error()})
	}

	log.Warn("OVERLAY GOT URI", "uri", uri)

	m := new(proto.String)
	err = uri.MarshalTo(m)
	if err != nil {
		log.Error("fialed to convert uri to string ")
		return err
	}

	// if !uri.Is(&proto.String{}) {
	// 	log.Error("uri was is not a string, convertying", "type", uri.Name().Name())
	// }

	// m := new(proto.HttpRequest)
	// err = req.MarshalTo(m)
	// if err != nil {
	// 	log.Error("failed to marshal ARG representation to httpRequest", "err")
	// 	return err
	// }

	// reformat the identifier for context of wrapped space
	// build new res:// scheme from overlay prefix's root
	// i.e. res://app/helloworld -> uri=/helloworld -> res://helloworld
	id := NewIdentifier(fmt.Sprintf("res://%s", m.GetValue()))

	// inject the wrapped space into the request scope and
	// issue request into our wrapped space which is otherwise
	// unavailable to outside of the overlay

	ctx.InjectSpace(Space{o.Space})

	req := ctx.CreateRequest(id)
	req.m.Verb = ctx.Request().m.Verb
	req.SetRepresentationClass(ctx.Request().m.RepresentationClass)

	log.Info("issuing request to wrapped space", "identifier", id)
	resp, err := ctx.IssueRequest(req)
	if err != nil {
		log.Error("failed to issue request into wrapped space", "err", err)
		return NewRepresentation(&proto.ErrorMessage{Message: err.Error()})
	}

	// transparent hook, cannot modify response
	o.onResponse(ctx, resp)
	log.Warn("-------------------------------------")

	return resp
}
