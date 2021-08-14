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
	grammar Grammar
	Mutable bool
}

func NewFilesetRegex(rx string) Fileset {
	grammar, err := NewGrammar(rx, GroupElement{
		GroupElement: &proto.GroupElement{
			Regex: rx,
			Name:  "regex",
		},
	})
	if err != nil {
		// log.Error(err)
		panic(err)
	}

	return Fileset{
		BaseEndpoint: BaseEndpoint{},
		grammar:      grammar,
		Mutable:      false,
	}
}

func (f Fileset) Grammar() Grammar {
	if f.grammar.m.Base != "" {
		return f.grammar
	}
	if f.Regex != "" {
		g, _ := NewGrammar(f.Regex, GroupElement{
			GroupElement: &proto.GroupElement{
				Regex: f.Regex,
				Name:  "regex",
			},
		})
		return g
	}
	return Grammar{}

}

func (e Fileset) Definition() EndpointDefinition {
	return EndpointDefinition{
		EndpointDefinition: &proto.EndpointDefinition{
			Name:    e.Grammar().String(),
			Type:    EndpointTypeFileset,
			Grammar: e.Grammar().m,
		},
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
