package roc

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/treethought/roc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const EndpointTypeTransport string = "transport"

// EndpointTransport is an endpoint that issues external events into the roc system
type Transport interface {
	Endpoint
	// Init is used to seed the transport with it's spatial scope
	Init(scope RequestScope) error
}

// Transport is a struct implementing the default behavior for an empty EndpointTransport
// This type is useful for embedding custom implementations of EndpointTransport
// and automatically handles scope initialization
type TransportImpl struct {
	*Accessor
	Scope  RequestScope
	OnInit func() error
}

func NewTransport(name string) *TransportImpl {
	return &TransportImpl{
		Accessor: NewAccessor(name),
		Scope:    RequestScope{},
		OnInit:   func() error { return nil },
	}
}

func (t *TransportImpl) Init(scope RequestScope) error {
	log.Debug("initializing transport scope")
	t.Scope = scope
	return t.OnInit()
}

type TransportPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl Transport
}

func (p *TransportPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &TransportRPCServer{Impl: p.Impl}, nil
}

func (TransportPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &TransportRPC{client: c}, nil
}

func (e *TransportPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
    proto.RegisterTransportServer(s, &TransportGRPCServer{
        Impl: e.Impl,
        broker: broker,
    })
	return nil
}

func (p *TransportPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &TransportGRPC{
		client: proto.NewTransportClient(c),
		broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &DispatcherPlugin{}

// ServeTransport starts the plugin's RPC server
// Because Transports typically will not implement the Resource methods,
// this can simply be called in a transport so that initial request scope can be initialized
func ServeTransport(e Transport) {
	// log.Debug("starting transport",
	// 	"name", e.Name,
	// 	"identifier", e.Identifier(),
	// )

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"transport": &TransportPlugin{Impl: e},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins:         pluginMap,
	})

}
