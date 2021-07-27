package roc

type ComparibleRepresentation interface {
	Equals(interface{}) bool
	HashCode() int
}

type Representation interface{}
