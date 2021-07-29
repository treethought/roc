package roc

import "fmt"

const EndpointTypeAccessor string = "accessor"

// EndpointAccessor is an endpoint that provides access
// to a set of resources or services within a space
type EndpointAccessor interface {
	Resource
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

func (e Accessor) CanResolve(request *Request) bool {
	return e.Grammar().Match(request.Identifier)
}

func (e Accessor) Evaluate(request *Request) Representation {

	switch request.Verb {
	case Source:
		return e.Source(request)
	case Sink:
		e.Sink(request)
		return nil
	case New:
		return e.New(request)
	case Delete:
		return e.Delete(request)
	case Exists:
		return e.Exists(request)

	default:
		return e.Source(request)

	}
}

func (e Accessor) String() string {
	return fmt.Sprintf("endpoint://%s", e.Grammar().String())
}

func (e Accessor) Source(request *Request) Representation {
	return nil
}

func (e Accessor) Sink(request *Request) {}

func (e Accessor) New(request *Request) Identifier {
	return ""
}
func (e Accessor) Delete(request *Request) bool {
	return false
}
func (e Accessor) Exists(request *Request) bool {
	return false
}
func (e Accessor) Transrept(request *Request) Representation {
	return nil
}
