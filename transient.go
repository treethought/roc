package roc

import (
	"fmt"

	"github.com/google/uuid"
	proto "github.com/treethought/roc/proto/v1"
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
	uid := uuid.New()
	uri := fmt.Sprintf("transient://%s", uid.String())

	repr := Representation{rep}
	log.Debug("creating transient endpoint",
		"type", repr.Type(),
		"uri", uri,
	)
	log.Trace(repr.String())

	grammar, err := NewGrammar(uri)
	if err != nil {
		log.Error("failed to create grammar", "err", err)
		panic(err)
	}

	return TransientEndpoint{
		BaseEndpoint:   BaseEndpoint{},
		Grammar:        grammar,
		Representation: repr,
	}
}

func (e *TransientEndpoint) Definition() EndpointDefinition {
	return EndpointDefinition{
		EndpointDefinition: &proto.EndpointDefinition{
			Name:    e.Grammar.String(),
			Type:    EndpointTypeTransient,
			Grammar: e.Grammar.m,
			Literal: e.Representation.message(),
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
		"type", e.Representation.Type(),
	)
	log.Trace(e.Representation.String())

	m, err := e.Representation.ToMessage()
	if err != nil {
		log.Error("failed to construct transient message", "err", err)
		return err
	}

	return m
}
