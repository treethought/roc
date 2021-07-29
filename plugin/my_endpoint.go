package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/treethought/roc"
)

// Here is a real implementation of Greeter
type MyEndpoint struct {
	logger     hclog.Logger
	Identifier roc.Identifier
}

// Source retrieves representation of resource
func (e *MyEndpoint) Source(request *roc.Request) roc.Representation {
	e.logger.Debug("Sourcing")
	return "hello world"
}

// Sink updates resource to reflect representation
func (e *MyEndpoint) Sink(request *roc.Request) {
}

// New creates a resource and return identifier for created resource
// If primary representation is included, use it to initialize resource state
func (e *MyEndpoint) New(request *roc.Request) roc.Identifier {
	return ""
}

// Delete remove the resource from the space that currently contains it
func (e *MyEndpoint) Delete(request *roc.Request) bool {
	return false
}

// Exists tests to see if resource can be resolved and exists
func (e *MyEndpoint) Exists(request *roc.Request) bool {
	return true
}

// CanResolve responds affirmatively if the endpoint can handle the request based on the identifier
func (e *MyEndpoint) CanResolve(request *roc.Request) bool {
	e.logger.Debug("checking if can resolve")
	return request.Identifier == e.Identifier
}

// Evaluate processes a request to create or return a Representation of the requested resource
func (e *MyEndpoint) Evaluate(request *roc.Request) roc.Representation {
	return e.Source(request)
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
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	endpoint := &MyEndpoint{
		logger:     logger,
		Identifier: "res://hello-world",
	}

	logger.Debug("starting plugin", "identifier", endpoint.Identifier)

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"endpoint": &roc.EndpointPlugin{Impl: endpoint},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
