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

func NewSpace(identifier Identifier, endpoints ...*proto.EndpointMeta) *proto.Space {
	space := &proto.Space{
		Identifier: identifier.String(),
		Imports:    []*proto.Space{},
		Endpoints:  endpoints,
		Clients:    make(map[string]*proto.ClientConfig),
	}

	log.Debug("created space", "identifier", space.GetIdentifier(), "endpoints", len(space.GetEndpoints()))
	return space
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

func startAccessors(s *proto.Space) {
	if s.Clients == nil {
		s.Clients = make(map[string]*proto.ClientConfig)
	}
	log.Warn("starting accessors", "space", s.Identifier)
	for _, e := range s.Endpoints {

		if e.Space != nil {
			startAccessors(e.Space)
		}

		if e.Type == EndpointTypeAccessor {
			log.Warn("starting accessor", "id", e.Identifier)

			phys := NewPhysicalEndpoint(e, nil)

			reconf := phys.Client.ReattachConfig()
			log.Trace("setting reattach config", "config", reconf)
			config := &proto.ClientConfig{
				Protocol:        string(reconf.Protocol),
				ProtocolVersion: uint32(reconf.ProtocolVersion),
				AddressNetwork:  reconf.Addr.Network(),
				AddressString:   reconf.Addr.String(),
				Pid:             uint32(reconf.Pid),
			}
			s.Clients[e.Identifier] = config
		}
	}

}

func canResolve(ctx *RequestContext, e *proto.EndpointMeta) bool {
	log.Trace(fmt.Sprintf("%+v", e))
	if e.Type == "transport" {
		return false
	}

	return matchGrammar(e.Grammar, ctx.m.Request.Identifier)
}

func resolveToEndpoint(s *proto.Space, ctx *RequestContext) (*proto.EndpointMeta, bool) {
	for _, ed := range s.GetEndpoints() {
		log.Trace("interrogating endpoint",
			"space", s.GetIdentifier(),
			"endpoint", ed.Identifier,
		)
		// TODO match grammar in endpoint or in space?
		// e := NewPhysicalEndpoint(ed.Cmd)
		// if e.CanResolve(ctx) {
		if canResolve(ctx, ed) {
			log.Debug("resolve affirmed", "endpoint_name", ed.Identifier, "cmd", ed.Cmd)
			return ed, true
		}
	}
	return &proto.EndpointMeta{}, false
}
