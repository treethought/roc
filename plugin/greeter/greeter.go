package main

import (
	"fmt"

	"github.com/treethought/roc"
)

type MyEndpoint struct {
	*roc.Accessor
}

func New(grammar roc.Grammar) *MyEndpoint {
	return &MyEndpoint{
		Accessor: roc.NewAccessor("greeter", grammar),
	}
}

// Source retrieves representation of resource
func (e *MyEndpoint) Source(ctx *roc.RequestContext) roc.Representation {
	e.Logger.Debug("Sourcing", ctx.Request.Identifier)
	e.Logger.Error("Making subrequest", "target", ctx.Request.Identifier)

	name, err := ctx.Source("res://namer", nil)
	if err != nil {
		e.Logger.Error("failed to dispatch request", "error", err)
	}
	return fmt.Sprintf("hello world: %s", name)
}

func main() {
	g := roc.NewGrammar("res://hello-world")
	endpoint := New(g)
	roc.Serve(endpoint)

}
