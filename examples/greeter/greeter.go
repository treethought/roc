package main

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/treethought/roc"
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
func (e *MyEndpoint) Source(ctx *roc.RequestContext) roc.Representation {
	log.Error("Executing source in greeter", "identifier", ctx.Request.Identifier)

	log.Warn("Making subrequest", "target", "res://name")
	log.Warn("sourcing http request")

	httpReq, err := ctx.Source(roc.Identifier("httpRequest://params"), nil)
	if err != nil {
		log.Error("failed to get httpRequest://", "err", err)
		return err
	}
	return httpReq

	name, err := ctx.GetArgumentValue("name")
	if err != nil {
		return err
	}

	if name == "" {
		name = "bonnie"
	}

	req := ctx.CreateRequest("res://name")

	req.SetArgumentByValue("nameArg", name)

	nameResp, err := ctx.IssueRequest(req)
	if err != nil {
		log.Error("failed to dispatch subrequest request", "error", err)
	}
	return fmt.Sprintf("hello world: %s", nameResp)
}

func main() {
	log.Error("GREETER MAIN")
	endpoint := New()
	roc.Serve(endpoint)

}
