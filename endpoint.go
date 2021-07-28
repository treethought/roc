package roc

// Endpoint represents the gateway between a logical resource and the computation
type Endpoint interface {
	Resource

	// CanResolve responds affirmatively if the endpoint can handle the request based on the identifier
	CanResolve(ctx RequestContext) bool

	// Grammer returns the defined set of identifiers that bind an endpoint to a Space
	Grammar() Grammar

	// Evaluate processes a request to create or return a Representation of the requested resource
	Evaluate(ctx RequestContext) Representation

	Type() string
	// Meta(ctx RequestArgument) map[string][]string
}

// func (e Endpoint) Meta(ctx RequestContext) MetaRepresentation {
// 	return nil
// }
