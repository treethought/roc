package roc

// Resource is an abstract model of information that is identified by one or more identifiers.
type Resource interface {
	// Source retrieves representation of resource
	Source(ctx RequestContext) Representation

	// Sink updates resource to reflect representation
	Sink(ctx RequestContext)

	// New creates a resource and return identifier for created resource
	// If primary representation is included, use it to initialize resource state
	New(ctx RequestContext) Identifier

	// Delete remove the resource from the space that currently contains it
	Delete(ctx RequestContext) bool

	// Exists tests to see if resource can be resolved and exists
	Exists(ctx RequestContext) bool
}

type Transreptor interface {
	// Transrept converts primary representation into an alternate representation
	// specified by required representation field in the request
	Transrept(ctx RequestContext) Representation
}
