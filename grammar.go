package roc

import (
	"fmt"
	"regexp"
	"strings"

	proto "github.com/treethought/roc/proto/v1"
)

// https://mixedanalytics.com/blog/regex-match-number-subdirectories-url/
// var RegexTypes = map[string]regexp.Regexp{
// 	"anything":     regexp.MustCompile(".*"),
// 	"path-segment": regexp.MustCompile("^[^/]+/[^/]+[a-zA-Z0-9]$"),
// }

var activeURIRegex = regexp.MustCompile(`\+(?P<name>[^@]+)@(?P<value>[^\+]+)`)

func constructIdentifier(g *proto.Grammar, args map[string][]string) string {
	log.Debug("building identifier from grammar", "base", g.GetBase())
	if g.GetActive() != nil {
		i := g.Active.GetIdentifier()
		for k, v := range args {
			i = fmt.Sprintf("%s+%s@%s", i, k, v[0])
		}
		return i
	}
	// TODO handle different types of groyps
	// like path vs regex
	// right now assuming always path based
	// if !strings.HasSuffix(i, "/") {
	// i := fmt.Sprintf("%s/", i)
	// return
	// }

	if g.GetBase() != "" {
		i := g.GetBase()
		for _, g := range g.GetGroups() {
			val, ok := args[g.Name]
			if ok {
				i = fmt.Sprintf("%s%s/", i, val[0])
			}
		}
		return i
	}

	return ""
}

// parseGrammar returns the group or active arguments from n identifier
func parseGrammar(g *proto.Grammar, i string) (args map[string][]string) {
	log.Trace("parsing grammar", "grammar", g, "identitifier", i)
	args = make(map[string][]string)
	if g == nil {
		log.Error("grammar was nil")
		return args
	}

	// strip the base and perform match on remaining portion
	path := strings.Replace(i, g.Base, "", 1)

	if g.Active != nil {
		return parseActive(g.Active, path)

	}

	for _, group := range g.GetGroups() {
		for k, v := range parseGroupElement(group, path) {
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

func matchGrammar(g *proto.Grammar, i string) bool {
	log.Trace("matching grammar",
		"grammar", g.GetBase(),
		"identitifier", i,
	)

	if !strings.HasPrefix(i, g.GetBase()) {
		log.Trace("identifier does not start with base", "id", i, "g", g.GetBase())
		return false
	}

	// strip the base and perform match on remaining portion
	path := strings.Replace(i, g.Base, "", 1)

	for _, p := range g.Groups {
		if !matchGroupElement(p, path) {
			log.Trace("group does not match", "group", p.Name)
			return false
		}
	}

	if len(g.GetActive().GetArguments()) > 0 {
		if !matchActive(g.Active, path) {
			log.Debug("active element does not match", "active", g.Active)
			return false
		}
	}

	log.Trace("grammar matches",
		"grammar", g.Base,
		"identifier", i,
	)

	return true
}

func matchGroupElement(g *proto.GroupElement, i string) bool {
	log.Trace("performing match on grammar group element")
	// TODO remove base before passing to this func
	// part := strings.Replace(i, g.Base, "", 1)
	if g.GetRegex() != "" {
		rx, err := regexp.Compile(g.Regex)
		if err != nil {
			log.Error("regex invalid", "regex", g.Regex, "err", err)
			return false
		}
		if rx.MatchString(i) {
			log.Debug("grammar group regex match", "regex", g.Regex, "identifier", i)
			return true
		}
	}
	return false
}

func parseGroupElement(g *proto.GroupElement, i string) (args map[string][]string) {
	log.Trace("parsing group element")
	args = make(map[string][]string)
	// TODO remove base before passing to this func
	// parts := strings.Replace(i, g.m.Base, "", 1)
	if g.Regex != "" {
		rx := regexp.MustCompile(g.Regex)
		matches := rx.FindAllString(i, -1)
		args[g.Name] = matches
	}
	return args
}

// active:toUpper+operand@file:/example.txt

func matchActive(a *proto.ActiveElement, i string) bool {
	log.Debug("performing match on active element")
	return activeURIRegex.MatchString(i)
}

func parseActive(a *proto.ActiveElement, i string) (args map[string][]string) {
	log.Trace("parsing active grammar")
	args = make(map[string][]string)
	matches := activeURIRegex.FindAllStringSubmatch(i, -1)

	// TODO fix nested active - prbly include base in regex
	for _, m := range matches {
		name, val := m[1], m[2]

		for _, arg := range a.Arguments {
			if arg.Name == name {
				log.Debug("parsed active arg", "name", m[1], "val", m[2])
				_, ok := args[name]
				if !ok {
					args[name] = []string{}
				}
				args[name] = append(args[name], val)
			}
		}

	}
	return args
}
