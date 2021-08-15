package roc

import (
	"io/ioutil"
	"strings"

	proto "github.com/treethought/roc/proto/v1"
)

const EndpointTypeFileset string = "fileset"

type Fileset struct {
	BaseEndpoint
	Regex   string
	grammar *proto.Grammar
	Mutable bool
}

func NewFilesetRegex(rx string) Fileset {
	grammar := &proto.Grammar{
		Base: rx,
		Groups: []*proto.GroupElement{
			{Regex: rx, Name: "regex"},
		},
	}

	return Fileset{
		BaseEndpoint: BaseEndpoint{},
		grammar:      grammar,
		Mutable:      false,
	}
}

func (f Fileset) Grammar() *proto.Grammar {
	if f.grammar.Base != "" {
		return f.grammar
	}

	if f.Regex != "" {
		grammar := &proto.Grammar{
			Base: f.Regex,
			Groups: []*proto.GroupElement{
				{Regex: f.Regex, Name: "regex"},
			},
		}
		return grammar
	}
	return &proto.Grammar{}

}

func (e Fileset) Definition() *proto.EndpointDefinition {
	return &proto.EndpointDefinition{
		Name:    e.Grammar().GetBase(),
		Type:    EndpointTypeFileset,
		Grammar: e.Grammar(),
	}
}

func (e Fileset) Source(ctx *RequestContext) interface{} {
	path := strings.Replace(ctx.Request().Identifier().String(), "res://", "", 1)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return NewRepresentation(&proto.ErrorMessage{Message: err.Error()})
	}
	return string(data)
}
