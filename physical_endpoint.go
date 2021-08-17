package roc

import (
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// PhysicalEndpoint represents an interface to a running  endpoint instance.
// The endpoint's binary is executed and may be killed via the Client
type PhysicalEndpoint struct {
	Endpoint
	Client *plugin.Client
	path   string
}

func Serve(e Endpoint) {
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"endpoint": &EndpointPlugin{Impl: e},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:        "phys-server",
			Level:       LogLevel,
			Color:       hclog.AutoColor,
			DisableTime: true,
		}),
	})

}

// NewPhysicalEndpoint starts a physical process and returns an RPC implementation
func NewPhysicalEndpoint(path string) Endpoint {
	// We're a host! Start by launching the plugin pss.
	log.Info("creating plugin client", "cmd", path)
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         PluginMap,
		Cmd:             exec.Command(path),
		Managed:         true,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:        "phys-client",
			Level:       LogLevel,
			Color:       hclog.AutoColor,
			DisableTime: true,
		}),
	})

	// Connect via RPC
	log.Info("connecting via rpc", "cmd", path)
	rpcClient, err := client.Client()
	if err != nil {
		log.Error("failed to connect via rpc", "endpoint", path, "error", err)
		os.Exit(1)
	}

	// RequestContext the plugin
	log.Info("dispensing plugin", "cmd", path)
	raw, err := rpcClient.Dispense("endpoint")
	if err != nil {
		log.Error("failed to dispense endpoint", "endpoint", path, "error", err)
		panic(err)
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

// PhysicalTransport represents an interface to a running  endpoint instance.
// The endpoint's binary is executed and may be killed via the Client
type PhysicalTransport struct {
	Transport
	Client *plugin.Client
}

// NewPhysicalTransport starts a physical process and returns an RPC implementation
func NewPhysicalTransport(path string) Transport {
	// We're a host! Start by launching the plugin pss.
	log.Info("creating transport plugin client", "cmd", path)
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         PluginMap,
		Cmd:             exec.Command(path),
		Managed:         true,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC},
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:        "transport-client",
			Level:       LogLevel,
			Color:       hclog.AutoColor,
			DisableTime: true,
		}),
	})

	// Connect via RPC
	log.Info("connecting via rpc", "cmd", path)
	rpcClient, err := client.Client()
	if err != nil {
		log.Error("failed to connect via rpc", "endpoint", path, "error", err)
		os.Exit(1)
	}

	// RequestContext the plugin
	log.Info("dispensing plugin", "cmd", path)
	raw, err := rpcClient.Dispense("transport")
	if err != nil {
		log.Error("failed to dispense endpoint", "endpoint", path, "error", err)
		panic(err)
		os.Exit(1)
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	transport := raw.(Transport)
	return PhysicalTransport{
		Transport: transport,
		Client:    client,
	}
}
