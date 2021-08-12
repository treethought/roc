package roc

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/go-hclog"
	"github.com/treethought/roc/proto"
	"gopkg.in/yaml.v3"
)

var log = hclog.New(&hclog.LoggerOptions{
	DisableTime: true,
	Level:       LogLevel,
	Color:       hclog.AutoColor,
})

type SpaceDefinition struct {
	Spaces []*proto.Space `json:"spaces" yaml:"spaces"`
}

type Space struct {
	m *proto.Space
}

type EndpointDefinition struct {
	*proto.EndpointDefinition
}

// func (ed EndpointDefinition) Type() string {
// 	return ed.EndpointType
// }

func (ed *EndpointDefinition) grammar() Grammar {
	elems := []GroupElement{}
	for _, g := range ed.Grammar.Groups {
		elems = append(elems, GroupElement{g})
	}
	g, err := NewGrammar(ed.Grammar.Base, elems...)
	if err != nil {
		panic(err)
	}
	return g
}

func NewSpace(identifier Identifier, endpoints ...EndpointDefinition) Space {
	s := Space{
		m: &proto.Space{
			Identifier: identifier.String(),
			Imports:    []*proto.Space{},
		},
	}
	for _, e := range endpoints {
		s.m.Endpoints = append(s.m.Endpoints, e.EndpointDefinition)
	}

	log.Debug("created space", "identifier", s.m.Identifier, "endpoints", len(s.m.Endpoints))
	return s
}

func (s *Space) BindEndpoint(e *proto.EndpointDefinition) {
	s.m.Endpoints = append(s.m.Endpoints, e)
}

func LoadSpaces(path string) ([]*proto.Space, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("failed to read config file", "error", err)
		return []*proto.Space{}, nil
	}
	// jsonBytes, _ := json.Unmarshal(data)

	def := SpaceDefinition{}
	err = yaml.Unmarshal(data, &def)
	if err != nil {
		log.Error("failed to parse space definition", err)
		return def.Spaces, fmt.Errorf("failed to parse space definitions")
	}

	out, _ := yaml.Marshal(def)
	fmt.Println(string(out))

	return def.Spaces, nil

}

func canResolve(ctx *RequestContext, e *proto.EndpointDefinition) bool {
	log.Trace(fmt.Sprintf("%+v", e))
	if e.Type == "transport" {
		return false
	}

	ed := EndpointDefinition{e}

	resolve := ed.grammar().Match(NewIdentifier(ctx.Request().m.Identifier))
	return resolve

}

func (s Space) Resolve(ctx *RequestContext, c chan (EndpointDefinition)) {
	for _, ed := range s.m.Endpoints {
		log.Debug("interrogating endpoint",
			"space", s.m.Identifier,
			"endpoint", ed.Name,
		)
		// TODO match grammar in endpoint or in space?
		// e := NewPhysicalEndpoint(ed.Cmd)
		// if e.CanResolve(ctx) {
		if canResolve(ctx, ed) {
			log.Debug("resolve affirmed", "endpoint_name", ed.Name, "cmd", ed.Cmd)
			c <- EndpointDefinition{ed}
			// close(c)
		}
	}
}
