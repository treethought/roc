package roc

import (
	"log"
)

type Space struct {
	identifier Identifier `yaml:"identifier,omitempty"`
	Endpoints  []Endpoint `yaml:"endpoints,omitempty"`
	Imports    []Space    `yaml:"imports,omitempty"`
	channel    chan (*Request)
}

func NewSpace(identifier Identifier, endpoints ...*PhysicalEndpoint) Space {
	s := Space{
		identifier: identifier,
		Imports:    []Space{},
		channel:    make(chan *Request),
	}
	for _, e := range endpoints {
		s.Endpoints = append(s.Endpoints, e.Impl)
	}
	return s
}

func (s Space) Identifier() Identifier {
	return s.identifier
}

func (s Space) Resolve(ctx *RequestContext, c chan (Endpoint)) {
	log.Printf("interrogating endpoints of space: %s", s.Identifier())
	for _, e := range s.Endpoints {

		log.Print("calling plugin")
		log.Printf("%+v", ctx.Dispatcher)

		if e.CanResolve(ctx) {
			log.Print("endpoint affirmed to resolve!: ")
			c <- e
		}
	}
}

// // Bind binds an endpoint to to the space using it's grammar
// func (s *Space) Bind(endpoint PhysicalEndpoint) {
// 	// TODO map of identifiers -> endpoint?
// 	s.Endpoints = append(s.Endpoints, endpoint)

// }
