package main

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/treethought/roc"
	"github.com/treethought/roc/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var log = hclog.Default()

type MyEndpoint struct {
	*roc.Accessor
}

func New() *MyEndpoint {
	return &MyEndpoint{
		Accessor: roc.NewAccessor("greeter"),
	}
}

func getAsString(rep roc.Representation) (*proto.String, error) {
	// TODO: this needs to be handled with response vlass behnd the scenes
	any, err := anypb.New(rep)
	if err != nil {
		return nil, err
	}

	m := new(proto.String)
	err = any.UnmarshalTo(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Source retrieves representation of resource
func (e *MyEndpoint) Source(ctx *roc.RequestContext) roc.Representation {
	log.Error("Executing source in greeter", "identifier", ctx.Request().Identifier)

	log.Warn("Making subrequest", "target", "res://name")
	log.Warn("sourcing http request")

	// httpReq, err := ctx.Source(roc.NewIdentifier("httpRequest://params"), nil)
	// if err != nil {
	// 	log.Error("failed to get httpRequest://", "err", err)
	// 	return &proto.ErrorMessage{Message: err.Error()}
	// }

	name, err := ctx.GetArgumentValue("name")
	if err != nil {
		return roc.NewRepresentation(&proto.ErrorMessage{Message: err.Error()})
	}

	m, err := getAsString(name)
	if err != nil {
		return roc.NewRepresentation(&proto.ErrorMessage{Message: err.Error()})
	}

	if m.Value == "" {
		m.Value = "bonne"
	}

	req := ctx.CreateRequest(roc.NewIdentifier("res://name"))

	req.SetArgumentByValue("nameArg", roc.NewRepresentation(m))

	nameResp, err := ctx.IssueRequest(req)
	if err != nil {
		log.Error("failed to dispatch subrequest request", "error", err)
	}

	m, err = getAsString(nameResp)
	if err != nil {
		return roc.NewRepresentation(&proto.ErrorMessage{Message: err.Error()})
	}

	m.Value = fmt.Sprintf("hello world: %s", m.Value)

	return roc.NewRepresentation(m)
}

func main() {
	log.Error("GREETER MAIN")
	endpoint := New()
	roc.Serve(endpoint)

}
