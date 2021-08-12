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

	repProto, err := repToProto(NewRepresentation(rep))
	if err != nil {
		log.Error("failed to convert transient rep to proto", "err", err)
	}

	uid := uuid.New()
	uri := fmt.Sprintf("transient://%s", uid.String())

	any := repProto.GetValue()
	log.Warn("creating transient endpoint",
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
		Representation: NewRepresentation(rep),
	}
}

func (e *TransientEndpoint) Definition() EndpointDefinition {
	repProto, err := repToProto(e.Representation)
	if err != nil {
		log.Error("failed to set transient endpoint definition literal", "err", err)
		panic(err)
	}

	any := repProto.GetValue()
	log.Warn("creating transient definition",
		"any_url", any.TypeUrl,
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
			Literal: repProto,
		},
	}
}

func (e *TransientEndpoint) Identifier() Identifier {
	return NewIdentifier(e.Grammar.String())
}

func (e TransientEndpoint) Type() string {
	return EndpointTypeTransient
}

func (e TransientEndpoint) Source(ctx *RequestContext) Representation {
	return e.Representation
}
