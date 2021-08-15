package roc

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/go-hclog"
	proto "github.com/treethought/roc/proto/v1"
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

func NewSpace(identifier Identifier, endpoints ...*proto.EndpointDefinition) Space {
	s := Space{
		m: &proto.Space{
			Identifier: identifier.String(),
			Imports:    []*proto.Space{},
			Endpoints:  endpoints,
		},
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

	return matchGrammar(e.Grammar, ctx.m.Request.Identifier)

}

func (s Space) Resolve(ctx *RequestContext) (*proto.EndpointDefinition, bool) {
	for _, ed := range s.m.Endpoints {
		log.Trace("interrogating endpoint",
			"space", s.m.Identifier,
			"endpoint", ed.Name,
		)
		// TODO match grammar in endpoint or in space?
		// e := NewPhysicalEndpoint(ed.Cmd)
		// if e.CanResolve(ctx) {
		if canResolve(ctx, ed) {
			log.Debug("resolve affirmed", "endpoint_name", ed.Name, "cmd", ed.Cmd)
			return ed, true
		}
	}
	return &proto.EndpointDefinition{}, false
}
