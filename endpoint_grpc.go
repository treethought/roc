package roc

import (
	"fmt"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/treethought/roc/proto"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/anypb"
)

// EndpointGRPC is an implementation of KV that talks over RPC.
type EndpointGRPC struct {
	broker *plugin.GRPCBroker
	client proto.EndpointClient
}

// func newProtoGrammar(g Grammar) *proto.Grammar {
// 	pg := &proto.Grammar{
// 		Base: g.Base,
// 	}
// 	for _, group := range g.Groups {
// 		pgroup := &proto.GroupElement{
// 			Name:     group.Name,
// 			Min:      group.Min,
// 			Max:      group.Max,
// 			Encoding: group.Encoding,
// 			Regex:    group.Regex,
// 		}
// 		pg.Groups = append(pg.Groups, pgroup)
// 	}

// 	return pg
// }

// func protoToGrammar(p *proto.Grammar) Grammar {
// 	g, err := NewGrammar(p.Base)
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, group := range p.Groups {
// 		gel := GroupElement{
// 			Name:     group.Name,
// 			Min:      group.Min,
// 			Max:      group.Max,
// 			Encoding: group.Encoding,
// 			Regex:    group.Regex,
// 		}
// 		g.Groups = append(g.Groups, gel)
// 	}
// 	return g

// }

// func newProtoSpace(space Space) *proto.Space {
// 	protoSpace := &proto.Space{Identifier: fmt.Sprint(space.Identifier)}
// 	for _, ed := range space.m.EndpointDefinitions {
// 		protoSpace.EndpointDefinitions = append(protoSpace.EndpointDefinitions, &proto.EndpointDefinition{
// 			Name:         ed.Name,
// 			Cmd:          ed.Cmd,
// 			Grammar:      ed.Grammar,
// 			EndpointType: ed.EndpointType,
//             Literal:
// 			Literal:      &proto.Representation{Value: fmt.Sprint(ed.Literal)},
// 			Space:        newProtoSpace(ed.Space),
// 		})
// 	}
// 	for _, s := range space.Imports {
// 		protoSpace.Imports = append(protoSpace.Imports, newProtoSpace(s))
// 	}

// 	return protoSpace
// }

// func protoToSpace(p *proto.Space) Space {
// 	space := NewSpace(Identifier(p.Identifier))
// 	for _, ed := range p.EndpointDefinitions {
// 		space.EndpointDefinitions = append(space.EndpointDefinitions, EndpointDefinition{
// 			Name:         ed.Name,
// 			Cmd:          ed.Cmd,
// 			Grammar:      protoToGrammar(ed.Grammar),
// 			EndpointType: ed.EndpointType,
// 			Literal:      ed.Literal.Value,
// 			Space:        protoToSpace(ed.Space),
// 		})
// 	}

// 	for _, s := range p.Imports {
// 		space.Imports = append(space.Imports, protoToSpace(s))
// 	}
// 	return space

// }
//func newProtoMap(args map[string][]string) []*proto.MapField {

//	fields := []*proto.MapField{}
//	for k, v := range args {
//		p := &proto.MapField{Key: k, Value: v}
//		fields = append(fields, p)
//	}
//	return fields
//}
//func protoToMap(p []*proto.MapField) map[string][]string {
//	res := make(map[string][]string)
//	for _, f := range p {
//		res[f.Key] = f.Value
//	}
//	return res

//}

//func newProtoContext(ctx *RequestContext) *proto.RequestContext {
//	protoCtx := &proto.RequestContext{
//		Request: &proto.Request{
//			Identifier: fmt.Sprint(ctx.Request.Identifier),
//			Verb:       proto.Verb(ctx.Request.Verb),
//			//TODO
//			Arguments: newProtoMap(ctx.Request.Arguments),
//		},
//		Scope: &proto.RequestScope{
//			Spaces: []*proto.Space{},
//		},
//	}
//	for _, s := range ctx.Scope.Spaces {
//		protoCtx.Scope.Spaces = append(protoCtx.Scope.Spaces, newProtoSpace(s))

//	}
//	return protoCtx
//}

//func protoToContext(p *proto.RequestContext) *RequestContext {
//	verb, ok := proto.Verb_value[p.Request.Verb.String()]
//	if !ok {
//		panic("unsopported  verb")
//	}
//	ctx := NewRequestContext(Identifier(p.Request.Identifier), Verb(verb))
//	for _, s := range p.Scope.Spaces {
//		ctx.Scope.Spaces = append(ctx.Scope.Spaces, protoToSpace(s))
//	}

//	args := protoToMap(p.Request.Arguments)
//	ctx.Request.Arguments = args

//	return ctx

//}

func (m *EndpointGRPC) Source(ctx *RequestContext) Representation {
	log.Debug("making endpoint Source grpc call", "identifier", ctx.Request().Identifier)
	resp, err := m.client.Source(context.Background(), ctx.m)
	if err != nil {
		log.Error("error making grpc call", "error", err)
		panic(err)
	}

	log.Debug("received rpc call response", "resp", resp)
	return NewRepresentation(resp)
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

func repToProto(rep Representation) (*proto.Representation, error) {
	any, err := anypb.New(rep)
	if err != nil {
		log.Error("failed to construct any")
		return nil, err
	}

	log.Warn("translating representation to proto msg",
		"desc_name", rep.ProtoReflect().Descriptor().Name(),
		"any_url", any.TypeUrl,
	)

	// m := new(proto.Representation)
	// err = any.UnmarshalTo(m)
	// if err != nil {
	// 	log.Error("failed to unmarshal to rep")
	// 	return nil, err
	// }
	// return m, nil
	return &proto.Representation{Value: any}, nil
}

func (m *EndpointGRPCServer) Source(ctx context.Context, req *proto.RequestContext) (*proto.Representation, error) {
	log.Debug("begining endpoint grpc source server implementation")

	rocCtx := &RequestContext{req}
	rep := m.Impl.Source(rocCtx)

	log.Debug("returning source implementation as grpc response", "rep", rep)
	return repToProto(rep)
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
