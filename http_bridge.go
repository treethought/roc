package roc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/treethought/roc/proto"
	"google.golang.org/protobuf/types/known/anypb"
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
	log.Warn("creating new httpRequest definition")

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
		panic(err)
	}

	// grammar.m.Groups = append(grammar.m.Groups, typeGroup)

	lit, err := repToProto(NewRepresentation(repr))
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	ed := EndpointDefinition{
		EndpointDefinition: &proto.EndpointDefinition{
			Name:    "httpRequest",
			Grammar: grammar.m,
			Literal: lit,
			Type:    EndpointTypeHTTPRequestAccessor,
		},
	}
	log.Warn("created with literal", "class", reflect.TypeOf(ed.Literal))
	return ed

}

func NewHttpRequestEndpoint(ed EndpointDefinition) HttpRequestEndpoint {
	// any, err := anypb.New(ed.Literal.GetValue())
	// if err != nil {
	// 	log.Error(err.Error())
	// 	panic(err)
	// }

	any := ed.Literal.GetValue()
	log.Warn("creating httpRequest literal endpoint",
		"any_url", any.TypeUrl,
		"lit_type", ed.Literal.ProtoReflect().Descriptor().Name(),
	)

	m := new(proto.HttpRequest)
	err := any.UnmarshalTo(m)
	if err != nil {
		log.Error("httpRequest acessor literal is not a request", "literal", ed.Literal.GetValue().TypeUrl)
		log.Error(fmt.Sprintf("%+v", ed))
		panic("httpRequest acessor literal is not a request")

	}

	return HttpRequestEndpoint{
		BaseEndpoint: BaseEndpoint{},
		Grammar:      ed.grammar(),
		request:      m,
	}
}

func (e HttpRequestEndpoint) Definition() EndpointDefinition {

	lit, err := repToProto(NewRepresentation(e.request))
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	return EndpointDefinition{
		EndpointDefinition: &proto.EndpointDefinition{
			Name:    e.Grammar.String(),
			Grammar: e.Grammar.m,
			Type:    EndpointTypeHTTPRequestAccessor,
			Literal: lit,
		},
	}
}

func (e *HttpRequestEndpoint) Identifier() Identifier {
	return NewIdentifier(e.Grammar.String())
}

func (e HttpRequestEndpoint) Type() string {
	return EndpointTypeHTTPRequestAccessor
}

func (e HttpRequestEndpoint) Source(ctx *RequestContext) Representation {
	log.Error("sourcing httpRequest accessor")
	return NewRepresentation(e.request)

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

func (o HttpBridgeOverlay) Evaluate(ctx *RequestContext) Representation {
	// transparent hook, cannot modify response
	log.Warn("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	log.Info("getting httpRequest arg")

	// reqId := fmt.Sprintf("httpRequest:%s", url)
	req, err := ctx.GetArgumentValue("httpRequest")
	if err != nil {
		log.Error("failed to get requst representation")
		return NewRepresentation(&proto.ErrorMessage{Message: err.Error()})
	}

	any, err := anypb.New(req)
	if err != nil {
		log.Error(err.Error())
		return NewRepresentation(&proto.ErrorMessage{Message: err.Error()})
	}

	m := new(proto.HttpRequest)
	err = any.UnmarshalTo(m)
	if err != nil {
		log.Error("httpRequest acessor literal is not a request", "literal", any.GetTypeUrl())
		return NewRepresentation(&proto.ErrorMessage{Message: err.Error()})

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
		return NewRepresentation(&proto.ErrorMessage{Message: err.Error()})
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
