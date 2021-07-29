package roc

import (
	"log"
)

type Space struct {
	Identifier Identifier `yaml:"identifier,omitempty"`
	Endpoints  []Endpoint `yaml:"endpoints,omitempty"`
	Imports    []Space    `yaml:"imports,omitempty"`
	channel    chan (*Request)
}

func NewSpace(identifier Identifier) *Space {
	return &Space{
		Identifier: identifier,
		Endpoints:  []Endpoint{},
		Imports:    []Space{},
	}
}

func (s Space) Resolve(request *Request, c chan (Endpoint)) {
	log.Printf("interrogating endpoints of space: %s", s.Identifier)
	for _, e := range s.Endpoints {
		if e.CanResolve(request) {
			log.Printf("endpoint affirmed to resolve!: %s", e)
			c <- e
		}
	}
}

// Bind binds an endpoint to to the space using it's grammar
func (s *Space) Bind(endpoint Endpoint) {

	// TODO map of identifiers -> endpoint?
	s.Endpoints = append(s.Endpoints, endpoint)

}
