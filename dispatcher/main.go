package main

import (
	"log"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/treethought/roc"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

type DefaultDispatcher struct {
	logger hclog.Logger
}

func (d DefaultDispatcher) resolveEndpoint(ctx *roc.RequestContext) roc.Endpoint {
	d.logger.Info("resolving endpoint")
	c := make(chan (roc.Endpoint))
	for _, s := range ctx.Scope.Spaces {
		d.logger.Info("checking space: ", s.Identifier())
		go s.Resolve(ctx, c)
	}

	return <-c
}

func (d DefaultDispatcher) Dispatch(ctx *roc.RequestContext) (roc.Representation, error) {
	d.logger.Error("#################")
	d.logger.Error("DISPATCH PLUGIN RECEIVED CALL")
	log.Printf("WITH SPACES: %d", len(ctx.Scope.Spaces))
	log.Printf("WITH SPACES ENDPOINTS: %d", len(ctx.Scope.Spaces[0].Endpoints))

	d.logger.Error("dispatching request", "identifier", ctx.Request.Identifier, "spaces", ctx.Scope.Spaces)

	endpoint := d.resolveEndpoint(ctx)

	// TODO route verbs to methods
	rep := endpoint.Source(ctx)
	return rep, nil

}

func main() {
	// fmt.Println("STARTING DISPATCHER")
	// gob.Register(roc.EndpointRPC{})
	d := DefaultDispatcher{
		logger: hclog.New(&hclog.LoggerOptions{
			Level:      hclog.Trace,
			Output:     os.Stderr,
			JSONFormat: false,
		}),
	}

	hclog.SetDefault(d.logger)

	// var pluginMap = map[string]plugin.Plugin{
	// 	"endpoint":   &roc.EndpointPlugin{},
	// 	"dispatcher": &d,
	// }

	// d.logger.Debug("BEGINGINS SERVER")
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: roc.Handshake,
		// Plugins:         roc.PluginMap,
		Plugins: map[string]plugin.Plugin{
			"dispatcher": &roc.DispatcherPlugin{Impl: &d},
		},
	})
	// d.logger.Error("SERVE STOPPED")

}
