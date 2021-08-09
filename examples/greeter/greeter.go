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
	log.Info("Executing source in greeter", "identifier", ctx.Request.Identifier)

	log.Warn("Making subrequest", "target", "res://name")

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
	endpoint := New()
	roc.Serve(endpoint)

}
