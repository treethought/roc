package roc

import (
	"fmt"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/treethought/roc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type DispatcherGRPC struct {
	broker *plugin.GRPCBroker
	client proto.DispatcherClient
}

func (m *DispatcherGRPC) Dispatch(ctx *RequestContext) (Representation, error) {
	protoCtx := newProtoContext(ctx)
	resp, err := m.client.Dispatch(context.Background(), protoCtx)
    if resp == nil {
        return nil, fmt.Errorf("NIL RESPONSE FROM DISPATCH")
    }

	return resp.Value, err
}

type DispatcherGRPCServer struct {
	proto.UnimplementedDispatcherServer
	// This is the real implementation
	Impl Dispatcher

	broker *plugin.GRPCBroker
}

func (m *DispatcherGRPCServer) Dispatch(ctx context.Context, req *proto.RequestContext) (*proto.Representation, error) {
	rocCtx := protoToContext(req)
	rep, err := m.Impl.Dispatch(rocCtx)
    fmt.Println("GOT RESPONSE FROM DISPATCHED REQUEST")
    fmt.Println("representation", rep)
	return &proto.Representation{Value: fmt.Sprint(rep)}, err
}

func (e *DispatcherPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
    proto.RegisterDispatcherServer(s, &DispatcherGRPCServer{
        Impl: e.Impl,
        broker: broker,
    })
	return nil
}

func (p *DispatcherPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &DispatcherGRPC{
		client: proto.NewDispatcherClient(c),
		broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &DispatcherPlugin{}
