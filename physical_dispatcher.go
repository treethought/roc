package roc

import (
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

func DispatchRequest(ctx *RequestContext) (Representation, error) {
	log.Info("dispatching physical request",
		"identifier", ctx.Request.Identifier,
		"scope_size", len(ctx.Scope.Spaces),
	)

	// We're a host! Start by launching the plugin pss.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         PluginMap,
		Cmd:             exec.Command("./dispatcher/dispatcher"),
		Logger: hclog.New(&hclog.LoggerOptions{
			Level:       hclog.Debug,
			Output:      os.Stderr,
			JSONFormat:  false,
			Name:        "dispatcher",
			Color:       hclog.ForceColor,
			DisableTime: true,
		}),
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Error("failed to connect to dispatch via rpc", "error", err)
		os.Exit(1)
	}

	// RequestContext the plugin
	log.Info("dispensing dispatcher")
	raw, err := rpcClient.Dispense("dispatcher")
	if err != nil {
		log.Error("failed to dispense dispatcher", "error", err)
		os.Exit(1)
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	dispatcher := raw.(Dispatcher)
	return dispatcher.Dispatch(ctx)
}
