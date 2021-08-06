package roc

import (
	"fmt"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/treethought/roc/proto"
	"golang.org/x/net/context"
)

// EndpointGRPC is an implementation of KV that talks over RPC.
type EndpointGRPC struct {
	broker *plugin.GRPCBroker
	client proto.EndpointClient
}

func newProtoGrammar(g Grammar) *proto.Grammar {
	pg := &proto.Grammar{
		Base: g.Base,
	}
	for _, group := range g.Groups {
		pgroup := &proto.GroupElement{
			Name:     group.Name,
			Min:      group.Min,
			Max:      group.Max,
			Encoding: group.Encoding,
			Regex:    group.Regex,
		}
		pg.Groups = append(pg.Groups, pgroup)
	}

	return pg
}

func protoToGrammar(p *proto.Grammar) Grammar {
	g, err := NewGrammar(p.Base)
	if err != nil {
		panic(err)
	}
	for _, group := range p.Groups {
		gel := GroupElement{
			Name:     group.Name,
			Min:      group.Min,
			Max:      group.Max,
			Encoding: group.Encoding,
			Regex:    group.Regex,
		}
		g.Groups = append(g.Groups, gel)
	}
	return g

}

func newProtoSpace(space Space) *proto.Space {
	protoSpace := &proto.Space{Identifier: fmt.Sprint(space.Identifier)}
	for _, ed := range space.EndpointDefinitions {
		protoSpace.EndpointDefinitions = append(protoSpace.EndpointDefinitions, &proto.EndpointDefinition{
			Name:         ed.Name,
			Cmd:          ed.Cmd,
			Grammar:      newProtoGrammar(ed.Grammar),
			EndpointType: ed.EndpointType,
			Literal:      &proto.Representation{Value: fmt.Sprint(ed.Literal)},
		})
	}
	for _, s := range space.Imports {
		protoSpace.Imports = append(protoSpace.Imports, newProtoSpace(s))
	}

	return protoSpace
}

func protoToSpace(p *proto.Space) Space {
	space := NewSpace(Identifier(p.Identifier))
	for _, ed := range p.EndpointDefinitions {
		space.EndpointDefinitions = append(space.EndpointDefinitions, EndpointDefinition{
			Name:         ed.Name,
			Cmd:          ed.Cmd,
			Grammar:      protoToGrammar(ed.Grammar),
			EndpointType: ed.EndpointType,
			Literal:      ed.Literal.Value,
		})
	}

	for _, s := range p.Imports {
		space.Imports = append(space.Imports, protoToSpace(s))
	}
	return space

}
func newProtoMap(args map[string][]string) []*proto.MapField {

	fields := []*proto.MapField{}
	for k, v := range args {
		p := &proto.MapField{Key: k, Value: v}
		fields = append(fields, p)
	}
	return fields
}
func protoToMap(p []*proto.MapField) map[string][]string {
	res := make(map[string][]string)
	for _, f := range p {
		res[f.Key] = f.Value
	}
	return res

}

func newProtoContext(ctx *RequestContext) *proto.RequestContext {
	protoCtx := &proto.RequestContext{
		Request: &proto.Request{
			Identifier: fmt.Sprint(ctx.Request.Identifier),
			Verb:       proto.Verb(ctx.Request.Verb),
			//TODO
			Arguments: newProtoMap(ctx.Request.Arguments),
		},
		Scope: &proto.RequestScope{
			Spaces: []*proto.Space{},
		},
	}
	for _, s := range ctx.Scope.Spaces {
		protoCtx.Scope.Spaces = append(protoCtx.Scope.Spaces, newProtoSpace(s))

	}
	return protoCtx
}

func protoToContext(p *proto.RequestContext) *RequestContext {
	verb, ok := proto.Verb_value[p.Request.Verb.String()]
	if !ok {
		panic("unsopported  verb")
	}
	ctx := NewRequestContext(Identifier(p.Request.Identifier), Verb(verb))
	for _, s := range p.Scope.Spaces {
		ctx.Scope.Spaces = append(ctx.Scope.Spaces, protoToSpace(s))
	}

	args := protoToMap(p.Request.Arguments)
	ctx.Request.Arguments = args

	return ctx

}

func (m *EndpointGRPC) Evaluate(ctx *RequestContext) Representation {
	protoCtx := newProtoContext(ctx)

	resp, err := m.client.Evaluate(context.Background(), protoCtx)
	if err != nil {
		panic(err)
	}

	return resp.Value
}

func (m *EndpointGRPC) Source(ctx *RequestContext) Representation {
	log.Debug("making endpoint Source grpc call", "identifier", ctx.Request.Identifier)
	protoCtx := newProtoContext(ctx)

	resp, err := m.client.Source(context.Background(), protoCtx)
	if err != nil {
		log.Error("error making grpc call", "error", err)
		panic(err)
	}

	log.Debug("received rpc call response", "resp", resp)
	return resp.Value
}

func (m *EndpointGRPC) Sink(ctx *RequestContext) {
	protoCtx := newProtoContext(ctx)

	_, err := m.client.Sink(context.Background(), protoCtx)
	if err != nil {
		panic(err)
	}

	return
}

func (m *EndpointGRPC) New(ctx *RequestContext) Identifier {
	protoCtx := newProtoContext(ctx)

	resp, err := m.client.New(context.Background(), protoCtx)
	if err != nil {
		panic(err)
	}

	return Identifier(resp.Value)
}

func (m *EndpointGRPC) Delete(ctx *RequestContext) bool {
	protoCtx := newProtoContext(ctx)

	resp, err := m.client.Delete(context.Background(), protoCtx)
	if err != nil {
		panic(err)
	}

	return resp.Value
}

func (m *EndpointGRPC) Exists(ctx *RequestContext) bool {
	protoCtx := newProtoContext(ctx)

	resp, err := m.client.Exists(context.Background(), protoCtx)
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

// func (m *EndpointGRPCServer) setDispatchClient(ctx *RequestContext, brokerID uint32) (conn *grpc.ClientConn, err error) {
// 	log.Debug("setting dispatch client to handle grpc server request", "brokerID", brokerID)
// 	conn, err = m.broker.Dial(brokerID)
// 	if err != nil {
// 		log.Error("failed to create dispatcher client conn", "error", err)
// 		return nil, err
// 	}
//     // defer conn.Close()

// 	d := &DispatcherGRPC{
// 		// TODO: dont use same broker?
// 		// broker: &plugin.GRPCBroker{},
// 		client: proto.NewDispatcherClient(conn),
// 	}
// 	ctx.Dispatcher = d
// 	log.Debug("set context dispatcher", "dispatcher", ctx.Dispatcher)
// 	return conn, nil

// }

func (m *EndpointGRPCServer) Evaluate(ctx context.Context, req *proto.RequestContext) (*proto.Representation, error) {

	rocCtx := protoToContext(req)
	rep := m.Impl.Evaluate(rocCtx)

	return &proto.Representation{Value: fmt.Sprint(rep)}, nil
}

func (m *EndpointGRPCServer) Source(ctx context.Context, req *proto.RequestContext) (*proto.Representation, error) {
	log.Debug("begining endpoint grpc source server implementation")

	rocCtx := protoToContext(req)
	// _, err := m.setDispatchClient(rocCtx, req.DispatcherServer)
	// defer conn.Close()

	rep := m.Impl.Source(rocCtx)

	log.Debug("returning source implementation as grpc response", "rep", rep)
	return &proto.Representation{Value: fmt.Sprint(rep)}, nil
}

func (m *EndpointGRPCServer) Sink(ctx context.Context, req *proto.RequestContext) (*proto.Empty, error) {
	rocCtx := protoToContext(req)

	m.Impl.Sink(rocCtx)
	return &proto.Empty{}, nil
}

func (m *EndpointGRPCServer) New(ctx context.Context, req *proto.RequestContext) (*proto.IdentifierResponse, error) {
	rocCtx := protoToContext(req)

	ident := m.Impl.New(rocCtx)
	resp := &proto.IdentifierResponse{
		Value: fmt.Sprint(ident),
	}
	return resp, nil
}

func (m *EndpointGRPCServer) Delete(ctx context.Context, req *proto.RequestContext) (*proto.BoolResponse, error) {
	rocCtx := protoToContext(req)

	result := m.Impl.Delete(rocCtx)
	resp := &proto.BoolResponse{
		Value: result,
	}
	return resp, nil
}

func (m *EndpointGRPCServer) Exists(ctx context.Context, req *proto.RequestContext) (*proto.BoolResponse, error) {
	rocCtx := protoToContext(req)

	result := m.Impl.Exists(rocCtx)
	resp := &proto.BoolResponse{
		Value: result,
	}
	return resp, nil
}
