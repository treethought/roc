package roc

import (
	"net/http"
)

// Request represents request for a resouce
type Request struct {
	Identifier          Identifier
	Verb                Verb
	RepresentationClass RepresentationClass
	Arguments           map[string][]string
	Headers             http.Header
	argumentValues      map[string][]Representation
}

func NewRequest(i Identifier, verb Verb, class RepresentationClass) *Request {
	return &Request{
		Identifier:          i,
		Verb:                verb,
		RepresentationClass: class,
		Arguments:           make(map[string][]string),
		argumentValues:      make(map[string][]Representation),
	}
}

// SetRepresentationClass sets the desired format of the representation response
func (r *Request) SetRepresentationClass(class RepresentationClass) {
	r.RepresentationClass = class
}

// SetArgument sets the value of an argument to an identifier
// The argument's representation can then be sources during evalutation
// This replaces any existing values, to append an identifier, use AddArgument
func (r *Request) SetArgument(name string, i Identifier) {
	r.Arguments[name] = []string{i.String()}
}

// AddArgument appends an identifier to any existing ones for the named arguement
func (r *Request) AddArgument(name string, i Identifier) {
	_, exists := r.Arguments[name]
	if !exists {
		r.Arguments[name] = []string{}
	}
	r.Arguments[name] = append(r.Arguments[name], i.String())
}

func (r *Request) SetArgumentByValue(name string, val Representation) {
	r.argumentValues[name] = []Representation{val}
}

// // Identifier returns the identifier of the requested resource
// func (r Request) Identifier() Identifier {
// 	return r.identifier
// }

// // Verb returns the specified action to be taken when evaluating the request
// func (r Request) Verb() Verb {
// 	return r.verb
// }

// func (r *Request) Headers() http.Header {
// 	return r.headers
// }

// func (r Request) Arguments() url.Values {
// 	return r.argmuments
// }
