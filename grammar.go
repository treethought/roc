package roc

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/treethought/roc/proto"
)

// https://mixedanalytics.com/blog/regex-match-number-subdirectories-url/
// var RegexTypes = map[string]regexp.Regexp{
// 	"anything":     regexp.MustCompile(".*"),
// 	"path-segment": regexp.MustCompile("^[^/]+/[^/]+[a-zA-Z0-9]$"),
// }

type grammar interface {
	ParseArguments(Identifier) map[string][]string
	Construct() Identifier
	Match(Identifier) bool
}

type Grammar struct {
	m   *proto.Grammar
	uri *url.URL
}

func NewGrammar(base string, elems ...GroupElement) (Grammar, error) {
	uri, err := url.Parse(base)
	if err != nil {
		log.Error(err.Error())
		return Grammar{}, err
	}

	grammar := Grammar{
		m: &proto.Grammar{
			Base:   uri.String(),
			Groups: []*proto.GroupElement{},
		},
		uri: uri,
	}

	for _, g := range elems {
		grammar.m.Groups = append(grammar.m.Groups, g.GroupElement)
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

	for _, group := range g.m.Groups {

		wrap := GroupElement{group}

		for k, v := range wrap.Parse(g, i) {
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
	log.Debug("checking grammar",
		"grammar", g.String(),
		"identitifier", i.String(),
	)
	log.Debug("parsing identitifier", "identifier", i.String())
	uri, err := url.Parse(i.String())
	if err != nil {
		log.Error("failed to parse identifier",
			"identifier", i,
			"error", err,
		)
		return false
	}

	log.Info("checking scheme", "uri_scheme", uri.Scheme, "grammar_scheme", g.uri.Scheme)
	if uri.Scheme != g.uri.Scheme {
		log.Debug("scheme does not match")
		return false
	}

	log.Trace("checking host", "uri_host", uri.Host, "grammar_host", g.uri.Host)
	if uri.Host != g.uri.Host {
		log.Debug("host does not match")
		return false
	}

	for _, p := range g.m.Groups {
		wrap := GroupElement{p}
		if !wrap.Match(g, i) {
			log.Debug("group does not match", "group", p.Name)
			return false
		}
	}

	// log.Info("checking path", "uri_path", uri.Path, "grammar_path", g.uri.Path)
	// if !(strings.HasPrefix(uri.Path, g.uri.Path)) {
	// 	// if uri.Path != g.uri.Path {
	// 	log.Info("path does not match")
	// 	return false
	// }
	log.Debug("grammar matches",
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
	*proto.GroupElement
}

func (e GroupElement) Match(g Grammar, i Identifier) bool {
	log.Debug("performing match grammar group element")
	part := strings.Replace(i.String(), g.m.Base, "", 1)
	if e.Regex != "" {
		rx, err := regexp.Compile(e.Regex)
		if err != nil {
			log.Error("regex invalid", "regex", e.Regex, "err", err)
			return false
		}
		if rx.MatchString(part) {
			log.Debug("grammar group regex match", "regex", e.Regex, "identifier", i)
			return true
		}
	}
	return false
}

func (e GroupElement) Parse(g Grammar, i Identifier) (args map[string][]string) {
	args = make(map[string][]string)
	parts := strings.Replace(i.String(), g.m.Base, "", 1)
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
