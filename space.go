package roc

import (
	"log"
)

type Space struct {
	Identifier Identifier `yaml:"identifier,omitempty"`
	Endpoints  []Endpoint `yaml:"endpoints,omitempty"`
	Imports    []Space    `yaml:"imports,omitempty"`
}

func (s Space) MatchEndpoint(ctx RequestContext, c chan (Endpoint)) {
	log.Printf("interrogating endpoints of space: %s", s.Identifier)
	for _, e := range s.Endpoints {
		if e.CanResolve(ctx) {
			log.Printf("endpoint affirmed to resolve!: %s", e)
			c <- e
		}
	}
}
