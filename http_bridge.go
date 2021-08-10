package roc

import (
	"fmt"
	"net/http"
	"reflect"
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
	request *http.Request
}

func NewHttpRequestDefinition(req *http.Request) EndpointDefinition {
	log.Warn("creating new httpRequest definition")
	grammar, err := NewGrammar("httpRequest://params")
	if err != nil {
		panic(err)
	}

	typeGroup := GroupElement{
		Name:  "type",
		Regex: "([^/].*)/",
	}
	// subGroup := GroupElement{
	// 	Name:  "sub",
	// 	Regex: ".*",
	// }

	grammar.Groups = append(grammar.Groups, typeGroup)
	ed := EndpointDefinition{
		Name:         "httpRequest",
		Grammar:      grammar,
		Literal:      req,
		EndpointType: EndpointTypeHTTPRequestAccessor,
	}
	log.Warn("created with literal", "class", reflect.TypeOf(ed.Literal))
	return ed

}

func NewHttpRequestEndpoint(ed EndpointDefinition) HttpRequestEndpoint {
	req, ok := ed.Literal.(*http.Request)
	if !ok {
		log.Error("httpRequest acessor literal is not a request", "literal", reflect.TypeOf(ed.Literal))
		log.Error(fmt.Sprintf("%+v", ed))
		panic("httpRequest acessor literal is not a request")
	}
	return HttpRequestEndpoint{
		BaseEndpoint: BaseEndpoint{},
		Grammar:      ed.Grammar,
		request:      req,
	}
}

func (e HttpRequestEndpoint) Definition() EndpointDefinition {
	return EndpointDefinition{
		Name:         e.Grammar.String(),
		Grammar:      e.Grammar,
		EndpointType: EndpointTypeHTTPRequestAccessor,
		Literal:      e.request,
	}
}

func (e *HttpRequestEndpoint) Identifier() Identifier {
	return Identifier(e.Grammar.String())
}

func (e HttpRequestEndpoint) Type() string {
	return EndpointTypeHTTPRequestAccessor
}

func (e HttpRequestEndpoint) Source(ctx *RequestContext) Representation {
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
		Grammar:      ed.Grammar,
		Space:        ed.Space,
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
		return err
	}

	httpReq, ok := req.(*http.Request)
	if !ok {
		log.Error("sourced value is not an http request")
		return fmt.Errorf("sourced value is not an http request")
	}

	log.Info("got httpReqeuest arg", "class", reflect.TypeOf(httpReq))

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

	log.Info("creating dynamic http request space")
	reqAccessor := NewHttpRequestDefinition(httpReq)

	space := NewSpace("space://http_bridge", reqAccessor)
	ctx.InjectSpace(space)

	identifier := Identifier(fmt.Sprintf("res:/%s", httpReq.URL.Path))

	log.Info("mapped http request to identifier", "identifier", identifier)

	rocReq := ctx.CreateRequest(identifier)

	ctx.InjectSpace(o.Space)

	repr, err := ctx.IssueRequest(rocReq)
	if err != nil {
		log.Error("failed to issue request into wrapped space", "err", err)
		return err
	}
	log.Warn("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")

	return repr
}
