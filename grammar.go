package roc

import (
	"net/url"
	"regexp"
	"strings"

	proto "github.com/treethought/roc/proto/v1"
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
	if g.m.Active != nil {
		return parseActive(g.m.Active, i.String())
	}

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
	log.Trace("checking grammar",
		"grammar", g.String(),
		"identitifier", i.String(),
	)

	if !strings.HasPrefix(i.String(), g.m.GetBase()) {
		return false
	}

	if len(g.m.GetActive().GetArguments()) > 0 {
		return matchActive(g.m.Active, i.String())
	}

	for _, p := range g.m.Groups {
		wrap := GroupElement{p}
		if !wrap.Match(g, i) {
			log.Trace("group does not match", "group", p.Name)
			return false
		}
	}

	log.Trace("grammar matches",
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
	log.Trace("performing match on grammar group element")
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

func matchActive(a *proto.ActiveElement, i string) bool {
	log.Debug("performing match on active element")
	regex := `\+(?P<name>[^@]+)@(?P<value>[^\+]+)`
	rx := regexp.MustCompile(regex)
	return rx.MatchString(i)
}

func parseActive(a *proto.ActiveElement, i string) (args map[string][]string) {
	log.Debug("parsing active grammar")
	args = make(map[string][]string)
	regex := `\+(?P<name>[^@]+)@(?P<value>[^\+]+)`
	rx := regexp.MustCompile(regex)
	matches := rx.FindAllStringSubmatch(i, -1)

	// TODO fix nested active - prbly include base in regex
	for _, m := range matches {
		name, val := m[1], m[2]
		log.Debug("parsed active arg", "name", m[1], "val", m[2])
		_, ok := args[name]
		if !ok {
			args[name] = []string{}
		}
		args[name] = append(args[name], val)
	}
	return args
}
