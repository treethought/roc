package roc

import (
	plugin "github.com/hashicorp/go-plugin"
	"github.com/treethought/roc/proto"
	"golang.org/x/net/context"
)

// EndpointGRPC is an implementation of KV that talks over RPC.
type TransportGRPC struct {
	EndpointGRPC
	broker *plugin.GRPCBroker
	client proto.TransportClient
}

func (m *TransportGRPC) Init(scope *RequestScope) error {

	protoScope := &proto.RequestScope{}
	for _, s := range scope.Spaces {
		protoScope.Spaces = append(protoScope.Spaces, newProtoSpace(s))
	}

	_, err := m.client.Init(context.Background(), protoScope)
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

func (m *TransportGRPCServer) Init(ctx context.Context, req *proto.RequestScope) (*proto.Empty, error) {
	scope := RequestScope{}
	for _, s := range req.Spaces {
		scope.Spaces = append(scope.Spaces, protoToSpace(s))
	}

	err := m.Impl.Init(scope)

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
