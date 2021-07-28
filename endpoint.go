package roc

import (
	"fmt"
)



type EndpointInteface interface {
	// TODO return Resolution Response
	Resolve(ctx RequestContext) bool
	Source(ctx RequestContext) Representation
	Sink(ctx RequestContext)
	New(ctx RequestContext) Identifier
	Delete(ctx RequestContext) bool
	Exists(ctx RequestContext) bool
	Transrept(ctx RequestContext) Representation
	Type() string
	// Meta(ctx RequestArgument) MetaRepresentation
}

const EndpointTypeAccessor string = "assesor"

type Endpoint struct {
	Grammar Grammar `yaml:"grammar,omitempty"`
}

func (e Endpoint) Type() string {
	return EndpointTypeAccessor
}

func (e Endpoint) String() string {
	return fmt.Sprintf("endpoint://%s", e.Grammar.String())
}

func (e Endpoint) Resolve(ctx RequestContext) bool {
	return e.Grammar.Match(ctx.Request.Identifier())
}

func (e Endpoint) Source(ctx RequestContext) Representation {
	return nil
}

func (e Endpoint) Sink(ctx RequestContext) {}

func (e Endpoint) New(ctx RequestContext) Identifier {
	return ""
}
func (e Endpoint) Delete(ctx RequestContext) bool {
	return false
}
func (e Endpoint) Exists(ctx RequestContext) bool {
	return false
}
func (e Endpoint) Transrept(ctx RequestContext) Representation {
	return nil
}

// func (e Endpoint) Meta(ctx RequestContext) MetaRepresentation {
// 	return nil
// }
