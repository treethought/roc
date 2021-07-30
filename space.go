package roc

import (
	"log"
)

type Space struct {
	identifier Identifier         `yaml:"identifier,omitempty"`
	Endpoints  []*PhysicalEndpoint `yaml:"endpoints,omitempty"`
	Imports    []Space            `yaml:"imports,omitempty"`
	channel    chan (*Request)
}

func NewSpace(identifier Identifier, endpoints ...*PhysicalEndpoint) Space {
	return Space{
		identifier: identifier,
		Endpoints:  endpoints,
		Imports:    []Space{},
		channel:    make(chan *Request),
	}
}

func (s Space) Identifier() Identifier {
    return s.identifier
}


func (s Space) Resolve(ctx *RequestContext, c chan (Endpoint)) {
	log.Printf("interrogating endpoints of space: %s", s.Identifier())
	for _, e := range s.Endpoints {

        log.Print("calling plugin")

		if e.Impl.CanResolve(ctx) {
			log.Print("endpoint affirmed to resolve!: ")
			c <- e.Impl
		}
	}
}

// // Bind binds an endpoint to to the space using it's grammar
// func (s *Space) Bind(endpoint PhysicalEndpoint) {
// 	// TODO map of identifiers -> endpoint?
// 	s.Endpoints = append(s.Endpoints, endpoint)

// }
