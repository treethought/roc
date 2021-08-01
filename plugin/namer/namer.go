package main

import (
	"github.com/treethought/roc"
)

type MyEndpoint struct {
	*roc.Accessor
}

func New(grammar roc.Grammar) *MyEndpoint {
	return &MyEndpoint{
		Accessor: roc.NewAccessor("namer", grammar),
	}
}

// Source retrieves representation of resource
func (e *MyEndpoint) Source(ctx *roc.RequestContext) roc.Representation {
	e.Logger.Debug("Sourcing", ctx.Request)
	ctx.Request.SetRepresentationClass(nil)
	return "BOBO"
}

func main() {
	grammar := roc.NewGrammar("res://namer")
	endpoint := New(grammar)
	roc.Serve(endpoint)

}
