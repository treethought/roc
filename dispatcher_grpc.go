package roc

import (
	"fmt"

	"github.com/treethought/roc/proto"
	"golang.org/x/net/context"
)

type DispatcherGRPC struct {
	client proto.DispatcherClient
}

func (m *DispatcherGRPC) Dispatch(ctx *RequestContext) (Representation, error) {
	log.Debug("making dispatch grpc call")

	protoCtx := newProtoContext(ctx)

	resp, err := m.client.Dispatch(context.Background(), protoCtx)
	if resp == nil {
        log.Error("failed to make dispatch grpc call", "err", err)
		return nil, fmt.Errorf("response from dispatch was nil")
	}

	return resp.Value, err
}

type DispatcherGRPCServer struct {
	proto.UnimplementedDispatcherServer
	// This is the real implementation
	Impl Dispatcher
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
