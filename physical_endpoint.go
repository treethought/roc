package roc

import (
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
)

// PhysicalEndpoint represents an interface to a running  endpoint instance.
// The endpoint's binary is executed and may be killed via the Client
type PhysicalEndpoint struct {
	Endpoint
	Client *plugin.Client
}

// NewPhysicalEndpoint starts a physical process and returns an RPC implementation
func NewPhysicalEndpoint(path string) Endpoint {
	// We're a host! Start by launching the plugin pss.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         PluginMap,
		Cmd:             exec.Command(path),
		// Logger:          logger,
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Error("failed to connect via rpc", "endpoint", path, "error", err)
		os.Exit(1)
	}
	// endpoint.rpc = rpcClient

	// RequestContext the plugin
	raw, err := rpcClient.Dispense("endpoint")
	if err != nil {
		log.Error("failed to dispense endpoint", "endpoint", path, "error", err)
		os.Exit(1)
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	endpoint := raw.(Endpoint)
	return PhysicalEndpoint{
		Endpoint: endpoint,
		Client:   client,
	}

}

// func (e *PhysicalEndpoint) Kill() {
// 	e.client.Kill()
// }
