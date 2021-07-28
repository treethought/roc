package roc

// Identifier is an opaque token that identifies a single resource
// a resource may have one or more identifiers
type Identifier string

type RepresentationClass interface {
	String() string
	Identifier() Identifier
}

type ComparibleRepresentation interface {
	Equals(interface{}) bool
	HashCode() int
}

type Representation interface{}
