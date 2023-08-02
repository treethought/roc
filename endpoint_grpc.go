package roc

import (
	"fmt"

	plugin "github.com/hashicorp/go-plugin"
	proto "github.com/treethought/roc/proto/v1"
	"golang.org/x/net/context"
)

// EndpointGRPC is an implementation of KV that talks over RPC.
type EndpointGRPC struct {
	broker *plugin.GRPCBroker
	client proto.EndpointClient
}

func (m *EndpointGRPC) Source(ctx *RequestContext) interface{} {
	log.Debug("making endpoint Source grpc call", "identifier", ctx.Request().Identifier())
	resp, err := m.client.Source(context.Background(), ctx.m)
	if err != nil {
		log.Error("error making grpc call", "error", err)
		panic(err)
	}
	return resp
}

func (m *EndpointGRPC) Sink(ctx *RequestContext) {
	_, err := m.client.Sink(context.Background(), ctx.m)
	if err != nil {
		panic(err)
	}

	return
}

func (m *EndpointGRPC) New(ctx *RequestContext) Identifier {
	resp, err := m.client.New(context.Background(), ctx.m)
	if err != nil {
		panic(err)
	}

	return NewIdentifier(resp.Value)

}

func (m *EndpointGRPC) Delete(ctx *RequestContext) bool {
	resp, err := m.client.Delete(context.Background(), ctx.m)
	if err != nil {
		panic(err)
	}

	return resp.Value
}

func (m *EndpointGRPC) Exists(ctx *RequestContext) bool {
	resp, err := m.client.Exists(context.Background(), ctx.m)
	if err != nil {
		panic(err)
	}

	return resp.Value
}

// // Here is the gRPC server that EndpointGRPC talks to.
type EndpointGRPCServer struct {
	proto.UnimplementedEndpointServer
	// This is the real implementation
	Impl Endpoint

	broker *plugin.GRPCBroker
}

func (m *EndpointGRPCServer) Source(ctx context.Context, req *proto.RequestContext) (*proto.Representation, error) {
	log.Trace("begining endpoint grpc source server implementation")

	rocCtx := &RequestContext{req}
	rep := m.Impl.Source(rocCtx)

	log.Debug("returning source implementation as grpc response", "rep", rep)
	return NewRepresentation(rep).message(), nil
}

func (m *EndpointGRPCServer) Sink(ctx context.Context, req *proto.RequestContext) (*proto.Empty, error) {
	rocCtx := &RequestContext{req}
	m.Impl.Sink(rocCtx)
	return &proto.Empty{}, nil
}

func (m *EndpointGRPCServer) New(ctx context.Context, req *proto.RequestContext) (*proto.IdentifierResponse, error) {
	rocCtx := &RequestContext{req}

	ident := m.Impl.New(rocCtx)
	resp := &proto.IdentifierResponse{
		Value: fmt.Sprint(ident),
	}
	return resp, nil
}

func (m *EndpointGRPCServer) Delete(ctx context.Context, req *proto.RequestContext) (*proto.BoolResponse, error) {
	rocCtx := &RequestContext{req}

	result := m.Impl.Delete(rocCtx)
	resp := &proto.BoolResponse{
		Value: result,
	}
	return resp, nil
}

func (m *EndpointGRPCServer) Exists(ctx context.Context, req *proto.RequestContext) (*proto.BoolResponse, error) {
	rocCtx := &RequestContext{req}

	result := m.Impl.Exists(rocCtx)
	resp := &proto.BoolResponse{
		Value: result,
	}
	return resp, nil
}
