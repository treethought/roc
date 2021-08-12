package roc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/treethought/roc/proto"
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
	log.Info("creating new httpRequest literal definition")

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
			Literal: litRep.Representation,
			Type:    EndpointTypeHTTPRequestAccessor,
		},
	}
	log.Warn("created with httpRequest def with literal", "type", litRep.Name().Name())
	return ed

}

func NewHttpRequestEndpoint(ed EndpointDefinition) HttpRequestEndpoint {
	rep := Representation{ed.Literal}

	// any := ed.Literal.GetValue()
	log.Warn("creating httpRequest literal endpoint for",
		"any_url", rep.Type(),
		"type_name", rep.Name().Name(),
	)

	if !rep.Is(&proto.HttpRequest{}) {
		log.Error("http accessor definition literal is not httpRequest, trying to convert", "type", rep.Type())
	}

	m := new(proto.HttpRequest)
	err := rep.MarshalTo(m)
	if err != nil {
		log.Error("failed to marshal representation to httpRequest")
		panic(err)
	}

	log.Warn("CREATING HTTP REQUEST ENDPOINT")

	return HttpRequestEndpoint{
		BaseEndpoint: BaseEndpoint{},
		Grammar:      ed.grammar(),
		request:      m,
	}
}

func (e HttpRequestEndpoint) Definition() EndpointDefinition {

	repr := NewRepresentation(e.request)
	log.Error("creating http request endpoint DEFI",
		"type", repr.Type(),
	)

	return EndpointDefinition{
		EndpointDefinition: &proto.EndpointDefinition{
			Name:    e.Grammar.String(),
			Grammar: e.Grammar.m,
			Type:    EndpointTypeHTTPRequestAccessor,
			Literal: repr.Representation,
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
	log.Error("sourcing httpRequest accessor")
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
	log.Warn("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	log.Info("getting httpRequest arg")

	// reqId := fmt.Sprintf("httpRequest:%s", url)
	req, err := ctx.GetArgumentValue("httpRequest")
	if err != nil {
		log.Error("failed to get requst representation")
		return NewRepresentation(&proto.ErrorMessage{Message: err.Error()})
	}
	log.Warn("received httpRequest arg representation",
		"type", req.Name().Name(),
		"type_url", req.Type(),
	)

	any := req.Representation.Value
	log.Warn("converting arg value to httpRequest", "any_url", any.TypeUrl)

	if !req.Is(&proto.HttpRequest{}) {
		log.Error("httprequest arg is not httpRequest, trying to convert", "type", req.Name().Name())
	}

	m := new(proto.HttpRequest)
	err = req.MarshalTo(m)
	if err != nil {
		log.Error("failed to marshal ARG representation to httpRequest", "err")
		return err
	}

	// httpReq, ok := req.(*http.Request)
	// if !ok {
	// 	log.Error("sourced value is not an http request")
	// 	return fmt.Errorf("sourced value is not an http request")
	// }

	log.Info("got httpReqeuest arg", "class", any.GetTypeUrl())

	// respId := fmt.Sprintf("httpResponse:%s", url)
	// resp, err := ctx.GetArgumentValue("httpResponse")
	// if err != nil {
	// 	log.Error("failed to get response representation")
	// 	return err
	// }
	// httpResp, ok = req.(http.ResponseWriter)
	// if !ok {
	// 	log.Error("sourced value is not an http response")
	// 	return fmt.Errorf("sourced value is not an http response")
	// }

	// construct dynamic space for request and response

	// make the golang http request
	// TODO: should make a better http request message
	httpReq, err := http.NewRequest(m.RequestMethod, m.RequestUrl, bytes.NewReader(m.RequestBody))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("creating dynamic http request space")
	reqAccessor := NewHttpRequestDefinition(httpReq)

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
	log.Warn("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")

	return repr
}
