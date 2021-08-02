package main

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/treethought/roc"
)

type DefaultDispatcher struct {
}

var log = hclog.New(&hclog.LoggerOptions{
	DisableTime: true,
})

func (d DefaultDispatcher) resolveEndpoint(ctx *roc.RequestContext) roc.Endpoint {
	log.Info("resolving request", "identifier", ctx.Request.Identifier)

	c := make(chan (roc.Endpoint))
	for _, s := range ctx.Scope.Spaces {
		log.Info("checking space: ", "space", s.Identifier)
		go s.Resolve(ctx, c)
	}

	return <-c
}

func (d DefaultDispatcher) Dispatch(ctx *roc.RequestContext) (roc.Representation, error) {
	log.Warn("receivied disptach call",
		"identifier", ctx.Request.Identifier,
		"scope_size", len(ctx.Scope.Spaces),
	)

	endpoint := d.resolveEndpoint(ctx)
	log.Info("resolved to endpoint")
	phys, ok := endpoint.(roc.PhysicalEndpoint)
	if !ok {
		return nil, fmt.Errorf("resolved to non-physical endpoint")
	}

	defer phys.Client.Kill()

	log.Info("evaluating request",
		"identifier", ctx.Request.Identifier,
	)
	// TODO route verbs to methods
	rep := endpoint.Source(ctx)
	log.Warn("returning response",
		"identifier", ctx.Request.Identifier,
		"representation", rep,
	)
	return rep, nil

}

func main() {
	d := DefaultDispatcher{
		// logger: hclog.New(&hclog.LoggerOptions{
		// 	Level:      hclog.Info,
		// 	Output:     os.Stderr,
		// 	JSONFormat: false,
		// 	Color:      hclog.ForceColor,
		// 	Name:       "dispatcher",
		// }),
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: roc.Handshake,
		Plugins: map[string]plugin.Plugin{
			"dispatcher": &roc.DispatcherPlugin{Impl: &d},
		},
	})

}
