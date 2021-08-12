package roc

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/treethought/roc/proto"
)

const EndpointTypeTransient string = "transient"

// TransientEndpoint is dynamically generated in-memory endpoint
// these are typically used for internal temporary resources.
type TransientEndpoint struct {
	BaseEndpoint
	Grammar        Grammar `yaml:"grammar,omitempty"`
	Representation Representation
}

func NewTransientEndpoint(rep *proto.Representation) TransientEndpoint {

	// repProto, err := repToProto(NewRepresentation(rep))
	// if err != nil {
	// 	log.Error("failed to convert transient rep to proto", "err", err)
	// }

	uid := uuid.New()
	uri := fmt.Sprintf("transient://%s", uid.String())

	repr := Representation{rep}
	any := repr.Any()
	log.Debug("creating transient endpoint from representation",
		"any_url", any.TypeUrl,
		"uri", uri,
		// "lit_type", ed.Literal.ProtoReflect().Descriptor().Name(),
	)

	grammar, err := NewGrammar(uri)
	if err != nil {
		panic(err)
	}

	return TransientEndpoint{
		BaseEndpoint:   BaseEndpoint{},
		Grammar:        grammar,
		Representation: repr,
	}
}

func (e *TransientEndpoint) Definition() EndpointDefinition {
	// repProto, err := repToProto(e.Representation)
	// if err != nil {
	// 	log.Error("failed to set transient endpoint definition literal", "err", err)
	// 	panic(err)
	// }

	any := e.Representation.Any()
	log.Debug("creating transient definition",
		"any_url", any.TypeUrl,
		"grammar", e.Grammar.String(),
		"literal", e.Representation.Any().String(),

		// "lit_type", ed.Literal.ProtoReflect().Descriptor().Name(),
	)

	// desc := e.Representation.ProtoReflect().Descriptor()
	// m := dynamicpb.NewMessage(desc)

	return EndpointDefinition{
		EndpointDefinition: &proto.EndpointDefinition{
			Name:    e.Grammar.String(),
			Type:    EndpointTypeTransient,
			Grammar: e.Grammar.m,
			// TODO: repProto?
			Literal: e.Representation.Representation,
		},
	}
}

func (e *TransientEndpoint) Identifier() Identifier {
	return NewIdentifier(e.Grammar.String())
}

func (e TransientEndpoint) Type() string {
	return EndpointTypeTransient
}

func (e TransientEndpoint) Source(ctx *RequestContext) interface{} {
	log.Debug("sourcing transient endpoint",
		"identifier", ctx.Request().Identifier(),
	)
	m, err := e.Representation.ToMessage()
	if err != nil {
		log.Error("failed to contruct concreate transient represent")
		return err
	}

	log.Info("returning transient representation",
		"type", m.ProtoReflect().Descriptor().FullName().Name(),
		"identifier", e.Identifier().String(),
		"gramar", e.Grammar.String(),
	)

	return m
}
