package roc

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v3"
)

var log = hclog.New(&hclog.LoggerOptions{
	DisableTime: true,
	Level:       hclog.Debug,
})

type SpaceDefinition struct {
	Spaces []Space `json:"spaces" yaml:"spaces"`
}

type EndpointDefinition struct {
	Name         string  `json:"name,omitempty" yaml:"name,omitempty"`
	Grammar      Grammar `json:"grammar,omitempty" yaml:"grammar,omitempty"`
	Cmd          string  `json:"cmd,omitempty" yaml:"cmd,omitempty"`
}

}

type Space struct {
	Identifier Identifier `yaml:"identifier,omitempty" json:"identifier,omitempty"`
	Imports    []Space    `yaml:"imports,omitempty" json:"imports,omitempty"`
	// use identifier instead of string, should reference
	// plugin binaries as a res:// or file://
	EndpointDefinitions []EndpointDefinition `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	channel             chan (*Request)
	logger              hclog.Logger
}

func NewSpace(identifier Identifier, endpoints ...EndpointDefinition) Space {
	s := Space{
		Identifier:          identifier,
		Imports:             []Space{},
		EndpointDefinitions: endpoints,
		channel:             make(chan *Request),
	}

	log.Debug("created space", "identifier", s.Identifier, "endpoints", len(s.EndpointDefinitions))
	return s
}

func (s *Space) BindEndpoint(e EndpointDefinition) {
	s.EndpointDefinitions = append(s.EndpointDefinitions, e)
}

func LoadSpaces(path string) ([]Space, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("failed to read config file", "error", err)
		return []Space{}, nil
	}

	def := SpaceDefinition{}
	err = yaml.Unmarshal(data, &def)
	if err != nil {
		log.Error("failed to parse space definition", err)
        return def.Spaces, fmt.Errorf("failed to parse space definitions")
	}
	return def.Spaces, nil

}

func canResolve(ctx *RequestContext, e EndpointDefinition) bool {
	log.Debug("checking grammar", "grammar", e.Grammar.String(), "identifier", ctx.Request.Identifier)
	resolve := e.Grammar.Match(ctx.Request.Identifier)
	return resolve

}

func (s Space) Resolve(ctx *RequestContext, c chan (Endpoint)) {
	for _, ed := range s.EndpointDefinitions {
		log.Info("interrogating endpoint",
			"space", s.Identifier,
			"endpoint", ed.Name,
		)
		// TODO match grammar in endpoint or in space?
		// e := NewPhysicalEndpoint(ed.Cmd)
		// if e.CanResolve(ctx) {
		if canResolve(ctx, ed) {
			log.Info("resolve affirmed", "endpoint_name", ed.Name, "cmd", ed.Cmd)
			c <- NewPhysicalEndpoint(ed.Cmd)
			close(c)
		}
	}
}

// // Bind binds an endpoint to to the space using it's grammar
// func (s *Space) Bind(endpoint PhysicalEndpoint) {
// 	// TODO map of identifiers -> endpoint?
// 	s.Endpoints = append(s.Endpoints, endpoint)

// }
