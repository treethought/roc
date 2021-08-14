package roc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	proto "github.com/treethought/roc/proto/v1"
)

var EndpointTypeHTTPRequestAccessor = "accessor:std"
var OverlayTypeHTTPBridge = "httpBridge"

type HttpBridgeOverlay struct {
	BaseEndpoint
	Grammar Grammar
	Space   Space
}

type HttpRequestEndpoint struct {
	BaseEndpoint
	Grammar Grammar `yaml:"grammar,omitempty"`
	// request *http.Request
	request *proto.HttpRequest
}

func NewHttpRequestMessage(req *http.Request) *proto.HttpRequest {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	repr := &proto.HttpRequest{
		RequestMethod: req.Method,
		RequestUrl:    req.URL.String(),
		UserAgent:     req.UserAgent(),
		RemoteIp:      req.RemoteAddr,
		Protocol:      req.Proto,
		RequestBody:   body,
	}
	return repr
}

func NewHttpRequestDefinition(req *http.Request) EndpointDefinition {
	log.Debug("creating new httpRequest literal definition")

	typeGroup := GroupElement{
		GroupElement: &proto.GroupElement{
			Name:  "type",
			Regex: "([^/].*)/",
		},
	}
	repr := NewHttpRequestMessage(req)

	// subGroup := GroupElement{
	// 	Name:  "sub",
	// 	Regex: ".*",
	// }
	grammar, err := NewGrammar("httpRequest://params", typeGroup)
	if err != nil {
		log.Error("failed to create httpRequest definition grammar", "err", err)
		panic(err)
	}

	// grammar.m.Groups = append(grammar.m.Groups, typeGroup)

	litRep := NewRepresentation(repr)
	ed := EndpointDefinition{
		EndpointDefinition: &proto.EndpointDefinition{
			Name:    "httpRequest",
			Grammar: grammar.m,
			Literal: litRep.m,
			Type:    EndpointTypeHTTPRequestAccessor,
		},
	}
	return ed

}

func NewHttpRequestEndpoint(ed EndpointDefinition) HttpRequestEndpoint {
	rep := Representation{ed.Literal}

	log.Info("creating httpRequest accessor",
		"type", rep.Type(),
	)

	if !rep.Is(&proto.HttpRequest{}) {
		log.Warn("http accessor definition literal is not httpRequest, trying to convert", "url", rep.URL())
	}

	m := new(proto.HttpRequest)
	err := rep.To(m)
	if err != nil {
		log.Error("failed to marshal representation to httpRequest", "err", err)
		panic(err)
	}

	return HttpRequestEndpoint{
		BaseEndpoint: BaseEndpoint{},
		Grammar:      ed.grammar(),
		request:      m,
	}
}

func (e HttpRequestEndpoint) Definition() EndpointDefinition {
	repr := NewRepresentation(e.request)

	return EndpointDefinition{
		EndpointDefinition: &proto.EndpointDefinition{
			Name:    e.Grammar.String(),
			Grammar: e.Grammar.m,
			Type:    EndpointTypeHTTPRequestAccessor,
			Literal: repr.message(),
		},
	}
}

func (e *HttpRequestEndpoint) Identifier() Identifier {
	return NewIdentifier(e.Grammar.String())
}

func (e HttpRequestEndpoint) Type() string {
	return EndpointTypeHTTPRequestAccessor
}

func (e HttpRequestEndpoint) Source(ctx *RequestContext) interface{} {
	log.Debug("sourcing httpRequest accessor")
	return e.request

	// part, err := ctx.GetArgumentValue("type")
	// if err != nil {
	// 	return err
	// }

	// switch part {
	// case "params":
	// 	return e.request.URL.Query()

	// case "url":
	// 	return e.request.URL

	// default:
	// 	return e.request
	// }

}

func NewHTTPBridgeOverlay(ed EndpointDefinition) HttpBridgeOverlay {
	return HttpBridgeOverlay{
		BaseEndpoint: BaseEndpoint{},
		Grammar:      ed.grammar(),
		Space:        Space{ed.Space},
	}
}

func (o HttpBridgeOverlay) Evaluate(ctx *RequestContext) interface{} {
	// transparent hook, cannot modify response
	log.Warn("http bridge evaluating request", "identifier", ctx.Request().Identifier())

	req, err := ctx.GetArgumentValue("httpRequest")
	if err != nil {
		log.Error("failed to get httpRequest representation", "err", err)
		return err
	}

	log.Debug("converting arg value to httpRequest", "url", req.URL())

	if !req.Is(&proto.HttpRequest{}) {
		log.Warn("httprequest arg is not httpRequest, will try to convert", "type", req.Type())
	}

	m := new(proto.HttpRequest)
	err = req.To(m)
	if err != nil {
		log.Error("failed to marshal representation to httpRequest", "err", err)
		return err
	}

	// construct dynamic space for request and response

	// make the golang http request
	// TODO: should make a better http request message
	httpReq, err := http.NewRequest(m.RequestMethod, m.RequestUrl, bytes.NewReader(m.RequestBody))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("creating dynamic http request space")

	httpDef := NewHttpRequestDefinition(httpReq)
	reqAccessor := NewHttpRequestEndpoint(httpDef).Definition()

	spaceID := NewIdentifier("space://http_bridge")
	space := NewSpace(spaceID, reqAccessor)
	ctx.InjectSpace(space)

	identifier := NewIdentifier(fmt.Sprintf("res:/%s", httpReq.URL.Path))

	log.Info("mapped http request to identifier", "identifier", identifier)

	rocReq := ctx.CreateRequest(identifier)

	ctx.InjectSpace(o.Space)

	repr, err := ctx.IssueRequest(rocReq)
	if err != nil {
		log.Error("failed to issue request into wrapped space", "err", err)
		return NewRepresentation(&proto.ErrorMessage{Message: err.Error()})
	}
	return repr
}
