package roc

import (
	"net/url"
	"os"
	"regexp"
)

type grammar interface {
	Parse(Identifier)
	Construct() Identifier
	Match(Identifier) bool
}

type Grammar struct {
	Base   string `json:"base,omitempty" yaml:"base,omitempty"`
	Groups []groupElement `json:"groups,omitempty" yaml:"groups,omitempty"`
	uri    *url.URL
}

func NewGrammar(base string, elems ...groupElement) (Grammar, error) {
	uri, err := url.Parse(base)
	if err != nil {
		log.Error(err.Error())
		return Grammar{}, err
	}

	grammar := Grammar{
		Base:   uri.String(),
		Groups: elems,
        uri: uri,
	}
	return grammar, nil
}

func (g Grammar) String() string {
	if g.uri == nil {
		return ""
	}
	return g.uri.String()
}

func (g Grammar) Match(i Identifier) bool {
	log.Debug("testing grammar",
		"grammar", g.String(),
		"identitifier", i,
	)
	uri, err := url.Parse(string(i))
	if err != nil {
		log.Error("failed to parse identifier",
			"identifier", i,
			"error", err,
		)

		return false
	}

	if uri.Scheme != g.uri.Scheme {
		return false
	}

	if uri.Host != g.uri.Host {
		return false
	}

	if uri.Path != g.uri.Path {
		return false
	}
	log.Info("grammar matches",
		"grammar", g.uri.String(),
		"identifier", i,
	)

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
