package roc

import (
	"github.com/hashicorp/go-hclog"
)

var log = hclog.Default()

type Space struct {
	identifier Identifier `yaml:"identifier,omitempty"`
	Imports    []Space    `yaml:"imports,omitempty"`
	channel    chan (*Request)

	// use identifier instead of string, should reference
	// plugin binaries as a res:// or file://
	Endpoints []string
}

func NewSpace(identifier Identifier, endpointPaths ...string) Space {
	s := Space{
		identifier: identifier,
		Imports:    []Space{},
		Endpoints:  endpointPaths,
		channel:    make(chan *Request),
	}

	log.Debug("created space", "identifier", s.Identifier(), "endpoints", s.Endpoints)
	return s
}

func (s Space) Identifier() Identifier {
	return s.identifier
}

func (s Space) Resolve(ctx *RequestContext, c chan (Endpoint)) {
	log.Debug("interrogating endpoints",
		"space", s.Identifier(),
	)
	for _, ePath := range s.Endpoints {
		e := NewPhysicalEndpoint(ePath)

		if e.CanResolve(ctx) {
			log.Info("resolve affirmed", "endpoint", ePath)
			c <- e
		}
	}
}

// // Bind binds an endpoint to to the space using it's grammar
// func (s *Space) Bind(endpoint PhysicalEndpoint) {
// 	// TODO map of identifiers -> endpoint?
// 	s.Endpoints = append(s.Endpoints, endpoint)

// }
