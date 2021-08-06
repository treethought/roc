package roc

import (
	"net/url"
	"regexp"
	"strings"
)

type grammar interface {
	ParseArguments(Identifier) map[string][]string
	Construct() Identifier
	Match(Identifier) bool
}

type Grammar struct {
	Base   string         `json:"base,omitempty" yaml:"base"`
	Groups []GroupElement `json:"groups,omitempty" yaml:"groups"`
	uri    *url.URL
}

func NewGrammar(base string, elems ...GroupElement) (Grammar, error) {
	uri, err := url.Parse(base)
	if err != nil {
		log.Error(err.Error())
		return Grammar{}, err
	}

	grammar := Grammar{
		Base:   uri.String(),
		Groups: elems,
		uri:    uri,
	}
	return grammar, nil
}

func (g Grammar) String() string {
	if g.uri == nil {
		return ""
	}
	return g.uri.String()
}

func (g Grammar) Parse(i Identifier) (args map[string][]string) {
	args = make(map[string][]string)

	for _, group := range g.Groups {
		for k, v := range group.Parse(g, i) {
			_, exists := args[k]
			if exists {
				args[k] = append(args[k], v...)
			} else {
				args[k] = v
			}
		}
	}

	return args
}

func (g Grammar) Match(i Identifier) bool {
	log.Debug("testing grammar",
		"grammar", g.String(),
		"identitifier", i,
	)
	log.Debug("parsing identitifier")
	uri, err := url.Parse(i.String())
	if err != nil {
		log.Error("failed to parse identifier",
			"identifier", i,
			"error", err,
		)
		return false
	}

	if uri.Scheme != g.uri.Scheme {
		log.Debug("scheme does not match")
		return false
	}

	if uri.Host != g.uri.Host {
		log.Debug("host does not match")
		return false
	}

	for _, p := range g.Groups {
		if !p.Match(g, i) {
			log.Debug("group does not match", "group", p.Name)
			return false
		}
	}

	// log.Debug("checking path")
	// if uri.Path != g.uri.Path {
	// 	log.Debug("path does not match")
	// 	return false
	// }
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
type GroupElement struct {
	Name     string `json:"name,omitempty" yaml:"name,omitempty"`
	Min      uint64 `json:"min,omitempty" yaml:"min,omitempty"`
	Max      uint64 `json:"max,omitempty" yaml:"max,omitempty"`
	Encoding string `json:"encoding,omitempty" yaml:"encoding,omitempty"`
	Regex    string `json:"regex,omitempty" yaml:"regex,omitempty"`
}

func (e GroupElement) Match(g Grammar, i Identifier) bool {
	log.Error("MATCHING GROUP")
	log.Debug("performing match grammar group element")
	part := strings.Replace(i.String(), g.Base, "", 1)
	if e.Regex != "" {
		rx, err := regexp.Compile(e.Regex)
		if err != nil {
			log.Error("regex invalid", "regex", e.Regex, "err", err)
			return false
		}
		if rx.MatchString(part) {
			return true
		}
	}
	return false
}

func (e GroupElement) Parse(g Grammar, i Identifier) (args map[string][]string) {
	args = make(map[string][]string)
	parts := strings.Replace(i.String(), g.Base, "", 1)
	if e.Regex != "" {
		rx := regexp.MustCompile(e.Regex)
		matches := rx.FindAllString(parts, -1)
		args[e.Name] = matches
	}
	return args
}

type optionalGroup struct {
	text string
}

type choiceElement struct {
	groups []GroupElement
}
