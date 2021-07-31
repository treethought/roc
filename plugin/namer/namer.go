package main

import (
	"log"
	"net/url"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/treethought/roc"
)

// Here is a real implementation of Greeter
type MyEndpoint struct {
	*roc.Accessor
	logger     hclog.Logger
	Identifier roc.Identifier
}

func New(grammar roc.Grammar) *MyEndpoint {
	return &MyEndpoint{
		Accessor: roc.NewAccessor(grammar),
		logger: hclog.New(&hclog.LoggerOptions{
			Level:      hclog.Trace,
			Output:     os.Stderr,
			JSONFormat: true,
		}),
	}

}

// Source retrieves representation of resource
func (e *MyEndpoint) Source(ctx *roc.RequestContext) roc.Representation {
	e.logger.Debug("Sourcing", ctx.Request)
	ctx.Request.SetRepresentationClass(nil)
	return "BOBO"
}

func main() {

	base, err := url.Parse("res://namer")
	if err != nil {
		log.Fatal(err)
	}

	grammar := roc.Grammar{
		Base: base,
	}
	endpoint := New(grammar)

	endpoint.logger.Debug("starting plugin", "identifier", endpoint.Grammar().String())

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"endpoint": &roc.EndpointPlugin{Impl: endpoint},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: roc.Handshake,
		Plugins:         pluginMap,
	})
}
