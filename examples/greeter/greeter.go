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
	e.Logger.Info("Sourcing in greeter", ctx.Request.Identifier)
	e.Logger.Warn("Making subrequest", "target", ctx.Request.Identifier)

	name, err := ctx.Source("res://name", nil)
	if err != nil {
		e.Logger.Error("failed to dispatch request", "error", err)
	}
	return fmt.Sprintf("hello world: %s", name)
}

func main() {
	endpoint := New()
	roc.Serve(endpoint)

}
