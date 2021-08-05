package roc

import (
	plugin "github.com/hashicorp/go-plugin"
	"github.com/treethought/roc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// EndpointGRPC is an implementation of KV that talks over RPC.
type TransportGRPC struct {
	EndpointGRPC
	broker *plugin.GRPCBroker
	client proto.TransportClient
}

func (m *TransportGRPC) setDispatchServer(msg *proto.InitTransport, dispatcher Dispatcher) (stop func()) {
	dispatchServer := &DispatcherGRPCServer{Impl: dispatcher}

	var s *grpc.Server

	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterDispatcherServer(s, dispatchServer)
		return s
	}

	log.Debug("starting ephemeral dispatch server for endpoiny grpc call")
	brokerID := m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	msg.DispatcherServer = brokerID
	return s.Stop
}

func (m *TransportGRPC) Init(msg *InitTransport) error {

	protoScope := &proto.RequestScope{}
	for _, s := range msg.Scope.Spaces {
		protoScope.Spaces = append(protoScope.Spaces, newProtoSpace(s))
	}

	protoMsg := &proto.InitTransport{Scope: protoScope}
	m.setDispatchServer(protoMsg, msg.Dispatcher)

	_, err := m.client.Init(context.Background(), protoMsg)
	if err != nil {
		return err
	}

	return nil
}

// // Here is the gRPC server that EndpointGRPC talks to.
type TransportGRPCServer struct {
	EndpointGRPCServer
	proto.UnimplementedTransportServer

	// This is the real implementation
	Impl Transport

	broker *plugin.GRPCBroker
}

func (m *TransportGRPCServer) setDispatchClient(msg *InitTransport, brokerID uint32) (conn *grpc.ClientConn, err error) {
	log.Debug("setting dispatch client to handle grpc server request", "brokerID", brokerID)
	conn, err = m.broker.Dial(brokerID)
	if err != nil {
		log.Error("failed to create dispatcher client conn", "error", err)
		return nil, err
	}

	d := &DispatcherGRPC{
		// TODO: dont use same broker?
		// broker: &plugin.GRPCBroker{},
		client: proto.NewDispatcherClient(conn),
	}
	msg.Dispatcher = d
	log.Debug("set context dispatcher", "dispatcher", msg.Dispatcher)
	return conn, nil

}

func (m *TransportGRPCServer) Init(ctx context.Context, req *proto.InitTransport) (*proto.Empty, error) {
	msg := &InitTransport{}
	for _, s := range req.Scope.Spaces {
		msg.Scope.Spaces = append(msg.Scope.Spaces, protoToSpace(s))
	}

	conn, err := m.setDispatchClient(msg, req.DispatcherServer)
	if err != nil {
		log.Error("error setting dispatch client", "err", err)
		return nil, err
	}
	defer conn.Close()

	err = m.Impl.Init(msg)

	return &proto.Empty{}, err
}

func (m *TransportGRPCServer) Evaluate(ctx context.Context, req *proto.RequestContext) (*proto.Representation, error) {
	return m.EndpointGRPCServer.Evaluate(ctx, req)
}

func (m *TransportGRPCServer) Source(ctx context.Context, req *proto.RequestContext) (*proto.Representation, error) {
	return m.EndpointGRPCServer.Source(ctx, req)
}

func (m *TransportGRPCServer) Sink(ctx context.Context, req *proto.RequestContext) (*proto.Empty, error) {
	return m.EndpointGRPCServer.Sink(ctx, req)
}

func (m *TransportGRPCServer) New(ctx context.Context, req *proto.RequestContext) (*proto.IdentifierResponse, error) {
	return m.EndpointGRPCServer.New(ctx, req)
}

func (m *TransportGRPCServer) Delete(ctx context.Context, req *proto.RequestContext) (*proto.BoolResponse, error) {
	return m.EndpointGRPCServer.Delete(ctx, req)
}

func (m *TransportGRPCServer) Exists(ctx context.Context, req *proto.RequestContext) (*proto.BoolResponse, error) {
	return m.EndpointGRPCServer.Exists(ctx, req)
}
