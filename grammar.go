package roc

import (
	"log"
	"net/url"
	"regexp"
)

type grammar interface {
	Parse(Identifier)
	Construct() Identifier
	Match(Identifier) bool
}

type Grammar struct {
	Base  *url.URL
	parts []grammarElement
}

func (g Grammar) String() string {
	return g.Base.String()
}

func (g Grammar) Match(i Identifier) bool {
	log.Printf("matching grammar %s against %s", g.String(), i)
	uri, err := url.Parse(string(i))
	if err != nil {
		log.Printf("failed to parse identifier %s", uri.String())
		return false
	}

	if uri.Scheme != g.Base.Scheme {
		return false
	}

	if uri.Host != g.Base.Host {
		return false
	}

	if uri.Path != g.Base.Path {
		return false
	}
	log.Printf("%s matches %s", uri.String(), g.Base.String())

	return true
}

type grammarElement struct {
	values []string
}

// groupElement defines segments of an identifier token
type groupElement struct {
	grammarElement
	name     string
	min      uint64
	max      uint
	encoding string
	regex    regexp.Regexp
}

type optionalGroup struct {
	text string
}

type choiceElement struct {
	groups []groupElement
}
