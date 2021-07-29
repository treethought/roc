package roc

import (
	"log"
)

type Space struct {
	Identifier Identifier         `yaml:"identifier,omitempty"`
	Endpoints  []PhysicalEndpoint `yaml:"endpoints,omitempty"`
	Imports    []Space            `yaml:"imports,omitempty"`
	channel    chan (*Request)
}

func NewSpace(identifier Identifier) Space {
	return Space{
		Identifier: identifier,
		Endpoints:  []PhysicalEndpoint{},
		Imports:    []Space{},
		channel:    make(chan *Request),
	}
}

func (s Space) Resolve(request *Request, c chan (PhysicalEndpoint)) {
	log.Printf("interrogating endpoints of space: %s", s.Identifier)
	for _, e := range s.Endpoints {
		if e.Impl.CanResolve(request) {
			log.Print("endpoint affirmed to resolve!: ")
			c <- e
		}
	}
}

// Bind binds an endpoint to to the space using it's grammar
func (s Space) Bind(endpoint PhysicalEndpoint) {
	// TODO map of identifiers -> endpoint?
	s.Endpoints = append(s.Endpoints, endpoint)

}
