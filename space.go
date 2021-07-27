package roc

import "log"

type Space struct {
	Identifier Identifier
	Endpoints  []EndpointInteface
	Imports    []Space
}

func (s Space) MatchEndpoint(ctx RequestContext, c chan (EndpointInteface)) {
	log.Printf("interrogating endpoints of space: %s", s.Identifier)
	for _, e := range s.Endpoints {
		if e.Resolve(ctx) {
            log.Printf("endpoint affirmed to resolve!: %s", e)
			c <- e
		}
	}
}
