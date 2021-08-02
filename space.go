package roc

import (
	"io/ioutil"
	"os"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v3"
)

var log = hclog.New(&hclog.LoggerOptions{
	DisableTime: true,
})

type SpaceDefinition struct {
	Spaces []Space `json:"spaces" yaml:"spaces"`
}

type GrammarDefinition struct {
	Base string `json:"base" yaml:"base"`
	// parts []grammarElement `json:"parts,omitempty" yaml:"parts,omitempty"`
}

type EndpointDefinition struct {
	Name    string            `json:"name,omitempty" yaml:"name,omitempty"`
	Grammar GrammarDefinition `json:"grammar,omitempty" yaml:"grammar,omitempty"`
	Cmd     string            `json:"cmd,omitempty" yaml:"cmd,omitempty"`
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

func LoadSpaces(path string) []Space {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("failed to read config file", "error", err)
		os.Exit(1)
	}

	def := SpaceDefinition{}
	err = yaml.Unmarshal(data, &def)
	if err != nil {
		log.Error("failed to parse space definition", err)
		os.Exit(1)
	}
	return def.Spaces

}

func canResolve(ctx *RequestContext, e EndpointDefinition) bool {
	grammar := NewGrammar(e.Grammar.Base)
	log.Debug("checking grammar", "grammar", grammar.String(), "identifier", ctx.Request.Identifier)
	resolve := grammar.Match(ctx.Request.Identifier)
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
