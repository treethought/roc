package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/treethought/roc"
)

var log = hclog.Default()

type MyEndpoint struct {
	*roc.Accessor
}

func New() *MyEndpoint {
	return &MyEndpoint{
		Accessor: roc.NewAccessor("namer"),
	}
}

// Source retrieves representation of resource
func (e *MyEndpoint) Source(ctx *roc.RequestContext) roc.Representation {
	e.Logger.Debug("Sourcing", ctx.Request)
	ctx.Request.SetRepresentationClass(nil)
	return "BOBO"
}

func main() {
	log.Error("STARTING NAMER")

	endpoint := New()
	roc.Serve(endpoint)

}
