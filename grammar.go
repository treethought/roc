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
	Base  *url.URL
	parts []grammarElement
}

func NewGrammar(base string) Grammar {
	uri, err := url.Parse(base)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	grammar := Grammar{
		Base: uri,
	}
	return grammar
}

func (g Grammar) String() string {
	if g.Base == nil {
		return ""
	}
	return g.Base.String()
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

	if uri.Scheme != g.Base.Scheme {
		return false
	}

	if uri.Host != g.Base.Host {
		return false
	}

	if uri.Path != g.Base.Path {
		return false
	}
	log.Info("grammar matches",
		"grammar", g.Base.String(),
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
