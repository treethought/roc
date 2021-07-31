package main

import (
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
	logger   hclog.Logger
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
	d.logger.Info("dispatching request", "identifier", ctx.Request.Identifier, "spaces", ctx.Scope.Spaces)


    

	// ctx.Dispatcher = d
	endpoint := d.resolveEndpoint(ctx)

	// ctx.Dispatcher = k.DispatchClient

	// phys, ok := endpoint.(PhysicalEndpoint)
	// if !ok {
	//     log.Println("resolved endpoint is not a plugin")
	//     return nil
	// }

	// log.Printf("resolved to endpoint: %s", phys.Impl.New)

	// TODO route verbs to methods
	rep := endpoint.Source(ctx)
	return rep, nil

}



func main() {
    // gob.Register(roc.EndpointRPC{})
	d := DefaultDispatcher{
		logger: hclog.New(&hclog.LoggerOptions{
			Level:      hclog.Trace,
			Output:     os.Stderr,
			JSONFormat: true,
		}),
    }

	d.logger.Debug("starting dispatcher")

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"dispatcher": &roc.DispatchPlugin{Impl: d},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
