package main

import (
	"fmt"
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
	logger hclog.Logger
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
	e.logger.Debug("Sourcing", ctx.Request.Identifier)
	e.logger.Error("Making subrequest", "target", ctx.Request.Identifier)

	name, err := ctx.Source("res://namer", nil)
	if err != nil {
		e.logger.Error("failed to dispatch request", "error", err)
	}
	return fmt.Sprintf("hello world: %s", name)
}


func main() {

	base, err := url.Parse("res://hello-world")
	if err != nil {
		log.Fatal(err)
	}

	grammar := roc.Grammar{
		Base: base,
	}
	endpoint := New(grammar)

	endpoint.logger.Debug("starting plugin",
		"grammar", endpoint.Grammar().String(),
		"identifier", endpoint.Identifier(),
	)

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"endpoint": &roc.EndpointPlugin{Impl: endpoint},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: roc.Handshake,
		Plugins:         pluginMap,
	})
}
