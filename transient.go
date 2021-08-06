package roc

import (
	"fmt"

	"github.com/google/uuid"
)

const EndpointTypeTransient string = "transient"

// TransientEndpoint is dynamically generated in-memory endpoint
// these are typically used for internal temporary resources.
type TransientEndpoint struct {
	BaseEndpoint
	Grammar        Grammar `yaml:"grammar,omitempty"`
	Representation Representation
}

func NewTransientEndpoint(rep Representation) TransientEndpoint {
	uid := uuid.New()
	uri := fmt.Sprintf("transient://%s", uid.String())

	grammar, err := NewGrammar(uri)
	if err != nil {
		panic(err)
	}

	return TransientEndpoint{
		BaseEndpoint:   BaseEndpoint{},
		Grammar:        grammar,
		Representation: rep,
	}
}

func (e *TransientEndpoint) Definition() EndpointDefinition {
	return EndpointDefinition{
		Name:         e.Grammar.String(),
		EndpointType: EndpointTypeTransient,
		Grammar:      e.Grammar,
		Literal:      e.Representation,
	}
}

func (e *TransientEndpoint) Identifier() Identifier {
	return Identifier(e.Grammar.String())
}

func (e TransientEndpoint) Type() string {
	return EndpointTypeTransient
}

func (e TransientEndpoint) Source(ctx *RequestContext) Representation {
	return e.Representation
}
