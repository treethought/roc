package roc

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/treethought/roc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const EndpointTypeTransport string = "transport"

type InitTransport struct {
	Scope      RequestScope
	Dispatcher Dispatcher
}

// EndpointTransport is an endpoint that issues external events into the roc system
type Transport interface {
	Endpoint
	// Init is used to seed the transport with it's spatial scope
	Init(ctx *InitTransport) error
}


// Transport is a struct implementing the default behavior for an empty EndpointTransport
// This type is useful for embedding custom implementations of EndpointTransport
// and automatically handles scope initialization
type TransportImpl struct {
	*Accessor
	Scope          RequestScope
	OnInit         func() error
	Dispatcher     Dispatcher
	broker         *plugin.GRPCBroker
	dispatchServer uint32
}

func (t TransportImpl) startDispatcher() {
	log.Info("starting transport's dispatcher")
	dispatcher := &CoreDispatcher{}
	dispatchServer := &DispatcherGRPCServer{Impl: dispatcher}

	var s *grpc.Server

	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterDispatcherServer(s, dispatchServer)
		return s
	}

	brokerID := t.broker.NextId()
	go t.broker.AcceptAndServe(brokerID, serverFunc)
	t.dispatchServer = brokerID
}

func NewTransport(name string) *TransportImpl {
	// this is done inside the transport plugin
	return &TransportImpl{
		Accessor:   NewAccessor(name),
		Scope:      RequestScope{},
		OnInit:     func() error { return nil },
		broker:     &plugin.GRPCBroker{},
		Dispatcher: &CoreDispatcher{},
	}
}

func (t *TransportImpl) Init(msg *InitTransport) error {
	log.Debug("initializing transport scope")
	t.Scope = msg.Scope
	t.Dispatcher = msg.Dispatcher
	return t.OnInit()
}

func (t *TransportImpl) Dispatch(ctx *RequestContext) (Representation, error) {
	if t.Dispatcher == nil {
		t.Dispatcher = &CoreDispatcher{}
	}
	for _, s := range t.Scope.Spaces {
		log.Debug("adding to scope", "space", s.Identifier)
		ctx.Scope.Spaces = append(ctx.Scope.Spaces, s)
	}

	// ctx.Scope = t.Scope

	log.Info("dispatching request from transport",
		"num_spaces", len(ctx.Scope.Spaces),
	)

	ctx.Dispatcher = t.Dispatcher
	return t.Dispatcher.Dispatch(ctx)
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
		Impl:   e.Impl,
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


// ServeTransport starts the plugin's RPC server
// Because Transports typically will not implement the Resource methods,
// this can simply be called in a transport so that initial request scope can be initialized
func ServeTransport(e Transport) {
	log.Info("serving transport")
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"transport": &TransportPlugin{Impl: e},
	}

	t, ok := e.(*TransportImpl)
	if ok {
		t.startDispatcher()
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
	})

}
