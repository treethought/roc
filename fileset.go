package roc

import (
	"io/ioutil"
	"strings"
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
		Regex: rx,
		Name:  "regex",
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
	if f.grammar.Base != "" {
		return f.grammar
	}
	if f.Regex != "" {
		g, _ := NewGrammar(f.Regex, GroupElement{
			Regex: f.Regex,
			Name:  "regex",
		})
		return g
	}
	return Grammar{}

}

func (e Fileset) Definition() EndpointDefinition {
	return EndpointDefinition{
		Name:         e.Grammar().String(),
		EndpointType: EndpointTypeFileset,
		Grammar:      e.Grammar(),
	}
}

func (e Fileset) Source(ctx *RequestContext) Representation {
	path := strings.Replace(ctx.Request.Identifier.String(), "res://", "", 1)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return string(data)
}
