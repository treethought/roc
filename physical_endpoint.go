package roc

import (
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	proto "github.com/treethought/roc/proto/v1"
)

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

// PhysicalEndpoint represents an interface to a running  endpoint instance.
// The endpoint's binary is executed and may be killed via the Client
type PhysicalEndpoint struct {
	Endpoint
	Client *plugin.Client
	path   string
}

type netAddr struct {
	network   string
	stringVal string
}

func (a netAddr) Network() string {
	return a.network
}

func (a netAddr) String() string {
	return a.stringVal
}

func attachExisting(config *proto.ClientConfig) *plugin.Client {
	log.Debug("re-attaching client to existing plugin", "config", config)

	addr := netAddr{network: config.AddressNetwork, stringVal: config.AddressString}

	reattach := &plugin.ReattachConfig{
		Protocol:        plugin.Protocol(config.Protocol),
		ProtocolVersion: int(config.ProtocolVersion),
		Addr:            addr,
		Pid:             int(config.Pid),
	}

	return plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         PluginMap,
		Reattach:        reattach,
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

}

func createPlugin(ed *proto.EndpointMeta) *plugin.Client {
	log.Debug("creating new plugin instance")
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         PluginMap,
		Cmd:             exec.Command(ed.GetCmd()),
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

	return client
}

// NewPhysicalEndpoint starts or connects to a physical process and returns an RPC client implementation
func NewPhysicalEndpoint(ed *proto.EndpointMeta, config *proto.ClientConfig) *PhysicalEndpoint {
	var client *plugin.Client
	log.Debug("creating plugin client", "cmd", ed.GetCmd(), "config", config)
	if config != nil {
		client = attachExisting(config)
	} else {
		client = createPlugin(ed)

	}

	// Connect via RPC
	log.Debug("connecting via rpc", "cmd", ed.GetCmd())
	rpcClient, err := client.Client()
	if err != nil {
		log.Error("failed to connect via rpc", "endpoint", ed.GetCmd(), "error", err)
		os.Exit(1)
	}

	// RequestContext the plugin
	log.Info("dispensing plugin", "cmd", ed.GetCmd())
	raw, err := rpcClient.Dispense("endpoint")
	if err != nil {
		log.Error("failed to dispense endpoint", "endpoint", ed.GetCmd(), "error", err)
		panic(err)
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	endpoint := raw.(Endpoint)
	return &PhysicalEndpoint{
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
	log.Debug("creating transport plugin client", "cmd", path)
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
			Color:       hclog.ForceColor,
			DisableTime: true,
		}),
	})

	// Connect via RPC
	log.Debug("connecting via rpc", "cmd", path)
	rpcClient, err := client.Client()
	if err != nil {
		log.Error("failed to connect via rpc", "endpoint", path, "error", err)
		os.Exit(1)
	}

	// RequestContext the plugin
	log.Info("dispensing transport plugin", "cmd", path)
	raw, err := rpcClient.Dispense("transport")
	if err != nil {
		log.Error("failed to dispense endpoint", "endpoint", path, "error", err)
		panic(err)
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	transport := raw.(Transport)
	return PhysicalTransport{
		Transport: transport,
		Client:    client,
	}
}
