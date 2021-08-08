package roc

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v3"
)

var log = hclog.New(&hclog.LoggerOptions{
	DisableTime: true,
	Level:       LogLevel,
	Color:       hclog.AutoColor,
})

type SpaceDefinition struct {
	Spaces []Space `json:"spaces" yaml:"spaces"`
}

type EndpointDefinition struct {
	Name         string         `json:"name,omitempty" yaml:"name,omitempty"`
	Grammar      Grammar        `json:"grammar,omitempty" yaml:"grammar,omitempty"`
	Cmd          string         `json:"cmd,omitempty" yaml:"cmd,omitempty"`
	EndpointType string         `json:"type,omitempty" yaml:"type,omitempty"`
	Literal      Representation `json:"literal,omitempty" yaml:"literal,omitempty"`

	// TODO generalize endpoint def for any endpoint/prototype
	Regex string `json:"regex,omitempty" yaml:"regex,omitempty"`
	// overlay wrapped space
	Space Space `json:"space,omitempty" yaml:"space,omitempty"`
}

func (ed EndpointDefinition) Type() string {
	if ed.EndpointType != "" {
		return ed.EndpointType
	}
	return EndpointTypeAccessor
}

type Space struct {
	Identifier Identifier `yaml:"identifier,omitempty" json:"identifier,omitempty"`
	Imports    []Space    `yaml:"imports,omitempty" json:"imports,omitempty"`
	// use identifier instead of string, should reference
	// plugin binaries as a res:// or file://
	EndpointDefinitions []EndpointDefinition `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
}

func NewSpace(identifier Identifier, endpoints ...EndpointDefinition) Space {
	s := Space{
		Identifier:          identifier,
		Imports:             []Space{},
		EndpointDefinitions: endpoints,
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

func (s Space) Resolve(ctx *RequestContext, c chan (EndpointDefinition)) {
	for _, ed := range s.EndpointDefinitions {
		log.Debug("interrogating endpoint",
			"space", s.Identifier,
			"endpoint", ed.Name,
		)
		// TODO match grammar in endpoint or in space?
		// e := NewPhysicalEndpoint(ed.Cmd)
		// if e.CanResolve(ctx) {
		if canResolve(ctx, ed) {
			log.Debug("resolve affirmed", "endpoint_name", ed.Name, "cmd", ed.Cmd)
			c <- ed
			close(c)
		}
	}
}
