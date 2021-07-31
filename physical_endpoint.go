package roc

import (
	"log"
	"os/exec"

	"github.com/hashicorp/go-plugin"
)

type PhysicalEndpoint struct {
	client *plugin.Client
	rpc    plugin.ClientProtocol
	Impl   Endpoint
}

func NewPhysicalEndpoint(path string) *PhysicalEndpoint {
	endpoint := &PhysicalEndpoint{}
	// We're a host! Start by launching the plugin pss.
	endpoint.client = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(path),
		// Logger:          logger,
	})

	// Connect via RPC
	rpcClient, err := endpoint.client.Client()
	if err != nil {
		log.Fatal(err)
	}
	endpoint.rpc = rpcClient

	// RequestContext the plugin
	raw, err := rpcClient.Dispense("endpoint")
	if err != nil {
		log.Fatal(err)
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	endpoint.Impl = raw.(Endpoint)
	return endpoint

}

func (e *PhysicalEndpoint) Kill() {
	e.client.Kill()
}
