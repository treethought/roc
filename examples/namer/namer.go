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
		Accessor: roc.NewAccessor("namer"),
	}
}

// Source retrieves representation of resource
func (e *MyEndpoint) Source(ctx *roc.RequestContext) roc.Representation {
	log.Debug("Sourcing", "request", ctx.Request)
	ctx.Request.SetRepresentationClass(nil)
	names, ok := ctx.Request.Arguments["nameArg"]
	if !ok {
		return "BOBO"
	}

	response := ""
	for _, nameRef := range names {
		resp, err := ctx.Source(roc.Identifier(nameRef), nil)
		if err != nil {
			return "oh no"
		}
        response = fmt.Sprintf("%s %s", response, resp)
	}

	return response
}

func main() {
	log.Error("STARTING NAMER")

	endpoint := New()
	roc.Serve(endpoint)

}
