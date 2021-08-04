package roc

import (
	"fmt"

	"github.com/hashicorp/go-plugin"
	"github.com/treethought/roc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type DispatcherGRPC struct {
	// broker *plugin.GRPCBroker
	client proto.DispatcherClient
}

func (m *DispatcherGRPC) Dispatch(ctx *RequestContext) (Representation, error) {
    log.Debug("making dispatch grpc call")
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

    // dont need this if not doing plugin
	broker *plugin.GRPCBroker
}

func (m *DispatcherGRPCServer) Dispatch(ctx context.Context, req *proto.RequestContext) (*proto.Representation, error) {
    log.Debug("peforming dispatch in grpc server")
	rocCtx := protoToContext(req)
	rep, err := m.Impl.Dispatch(rocCtx)
    if err != nil {
        log.Error("error calling dispatch implementation", "err", err)
        return nil, err
    }
	return &proto.Representation{Value: fmt.Sprint(rep)}, err
}

func (e *DispatcherPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
    proto.RegisterDispatcherServer(s, &DispatcherGRPCServer{
        Impl: e.Impl,
        // broker: broker,
    })
	return nil
}

func (p *DispatcherPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &DispatcherGRPC{
		client: proto.NewDispatcherClient(c),
		// broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &DispatcherPlugin{}
