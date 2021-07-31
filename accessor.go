package roc

import (
	"fmt"
)

const EndpointTypeAccessor string = "accessor"

// EndpointAccessor is an endpoint that provides access
// to a set of resources or services within a space
type EndpointAccessor interface {
	// shared.Resource
	Endpoint
}

// Accessor is a struct implementing the default behavior for an empty EndpointAccessor
// This type is useful for embedding custom implementations of EndpointAccessor
type Accessor struct {
	grammar Grammar `yaml:"grammar,omitempty"`
}

func NewAccessor(grammar Grammar) *Accessor {
	return &Accessor{
		grammar: grammar,
	}
}

func (e Accessor) Grammar() Grammar {
	return e.grammar
}
func (e *Accessor) SetGrammar(grammar Grammar) {
	e.grammar = grammar
}

func (e Accessor) Type() string {
	return EndpointTypeAccessor
}

func (e Accessor) CanResolve(ctx *RequestContext) bool {
	return e.Grammar().Match(ctx.Request.Identifier)
}

func (e Accessor) Evaluate(ctx *RequestContext) Representation {

	switch ctx.Request.Verb {
	case Source:
		return e.Source(ctx)
	case Sink:
		e.Sink(ctx)
		return nil
	case New:
		return e.New(ctx)
	case Delete:
		return e.Delete(ctx)
	case Exists:
		return e.Exists(ctx)

	default:
		return e.Source(ctx)

	}
}

func (e Accessor) String() string {
	return fmt.Sprintf("endpoint://%s", e.Grammar().String())
}

func (e Accessor) Source(ctx *RequestContext) Representation {
	return nil
}

func (e Accessor) Sink(ctx *RequestContext) {}

func (e Accessor) New(ctx *RequestContext) Identifier {
	return ""
}
func (e Accessor) Delete(ctx *RequestContext) bool {
	return false
}
func (e Accessor) Exists(ctx *RequestContext) bool {
	return false
}
func (e Accessor) Transrept(ctx *RequestContext) Representation {
	return nil
}
