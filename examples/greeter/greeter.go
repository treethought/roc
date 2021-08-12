package main

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/treethought/roc"
	"github.com/treethought/roc/proto"
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

// Source retrieves representation of resource
func (e *MyEndpoint) Source(ctx *roc.RequestContext) interface{} {
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
		log.Error("failed to get name argument", "err", err)
		// return roc.NewRepresentation(&proto.ErrorMessage{Message: err.Error)})
	}

	m := new(proto.String)
	err = name.MarshalTo(m)
	if err != nil {
		log.Error("failed to marshal argument to string", "err", err)
	}

	if m.GetValue() == "" {
		m.Value = "bonne"
	}

	req := ctx.CreateRequest(roc.NewIdentifier("res://name"))

	req.SetArgumentByValue("nameArg", roc.NewRepresentation(m))

	nameResp, err := ctx.IssueRequest(req)
	if err != nil {
		log.Error("failed to dispatch subrequest request", "error", err)
	}

	r := new(proto.String)
	err = nameResp.MarshalTo(r)
	if err != nil {
		r.Value = "MISSING NAME"
	}

	r.Value = fmt.Sprintf("hello world: %s", r.Value)

	return r
}

func main() {
	log.Error("GREETER MAIN")
	endpoint := New()
	roc.Serve(endpoint)

}
