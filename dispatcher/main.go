package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/treethought/roc"
)

type DefaultDispatcher struct {
	logger hclog.Logger
}

func (d DefaultDispatcher) resolveEndpoint(ctx *roc.RequestContext) roc.Endpoint {
	d.logger.Info("resolving endpoint")
	c := make(chan (roc.Endpoint))
	for _, s := range ctx.Scope.Spaces {
		d.logger.Debug("checking space: ", "space", s.Identifier())
		go s.Resolve(ctx, c)
	}

	return <-c
}

func (d DefaultDispatcher) Dispatch(ctx *roc.RequestContext) (roc.Representation, error) {
	d.logger.Debug("receivied disptach call",
		"identifier", ctx.Request.Identifier,
		"scope_size", len(ctx.Scope.Spaces),
	)

	endpoint := d.resolveEndpoint(ctx)
    phys, ok := endpoint.(roc.PhysicalEndpoint)
    if !ok {
        return nil, fmt.Errorf("resolved to non-physical endpoint")
    }

    defer phys.Client.Kill()

	d.logger.Info("dispatching request",
		"identifier", ctx.Request.Identifier,
	)
	// TODO route verbs to methods
	rep := endpoint.Source(ctx)
	return rep, nil

}

func main() {
	d := DefaultDispatcher{
		logger: hclog.New(&hclog.LoggerOptions{
			Level:      hclog.Info,
			Output:     os.Stderr,
			JSONFormat: false,
			Color:      hclog.ForceColor,
			Name:       "dispatcher",
		}),
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: roc.Handshake,
		Plugins: map[string]plugin.Plugin{
			"dispatcher": &roc.DispatcherPlugin{Impl: &d},
		},
	})

}
