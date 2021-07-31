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
	e.logger.Debug("Sourcing", ctx.Request.Identifier)
    e.logger.Error("Making subrequest", "target", ctx.Request.Identifier)

    name := "heh"
    // name, err := ctx.Source("res://namer", nil)
    // if err != nil {
    //     e.logger.Error("failed to dispatch request", "error", err)
    // }
    return fmt.Sprintf("hello world: %s", name)
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
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

	endpoint.logger.Debug("starting plugin", "identifier", endpoint.Grammar().String())

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"endpoint": &roc.EndpointPlugin{Impl: endpoint},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
