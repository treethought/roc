package roc

// Verb specifies the acton to be taken by endpoint when resolving request
type Verb int

const (
	// Source retrieves representation of resource
	Source Verb = iota + 1

	// Sink updates resource to reflect representation
	Sink

	// Exists tests to see if resource can be resolved and exists
	Exists

	// Delete remove the resource from the space that currently contains it
	Delete

	// New creates a resource and return identifier for created resource
	// If primary representation is included, use it to initialize resource state
	New

	// Transrept converts primary representation into an alternate representation
	// specified by required representation field in the request
	Transrept

	// Resolve performs resolution on the request passed as the primary representation
	Resolve

	// Meta retrieves a meta data representation for the identified space or space element
	Meta
)

// String - Creating common behavior - give the type a String function
func (v Verb) String() string {
	return [...]string{"SOURCE", "SINK", "EXISTS", "DELETE", "NEW", "TRANSREPT", "RESOLVE", "META"}[v-1]
}

func (v Verb) EnumIndex() int {
	return int(v)
}
