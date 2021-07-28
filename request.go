package roc

import (
	"net/http"
	"net/url"
)

// Request represents request for a resouce
type Request struct {
	identifier          Identifier
	verb                Verb
	representationClass RepresentationClass
	argmuments          url.Values
	headers             http.Header
}

func NewRequest(i Identifier, verb Verb, class RepresentationClass) *Request {
	return &Request{
		identifier:          i,
		verb:                verb,
		representationClass: class,
	}
}

// SetRepresentationClass sets the desired format of the representation response
func (r *Request) SetRepresentationClass(class RepresentationClass) {
	r.representationClass = class
}

// Identifier returns the identifier of the requested resource
func (r Request) Identifier() Identifier {
	return r.identifier
}

// Verb returns the specified action to be taken when evaluating the request
func (r Request) Verb() Verb {
	return r.verb
}

func (r *Request) Headers() http.Header {
	return r.headers
}

func (r Request) Arguments() url.Values {
	return r.argmuments
}
