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
}

func NewRequest(i Identifier, verb Verb, class RepresentationClass) *Request {
	return &Request{
		Identifier:          i,
		Verb:                verb,
		RepresentationClass: class,
		Arguments:           make(map[string][]string),
	}
}

// SetRepresentationClass sets the desired format of the representation response
func (r *Request) SetRepresentationClass(class RepresentationClass) {
	r.RepresentationClass = class
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
