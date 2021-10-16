package main

import (
	"fmt"
	"regexp"
)

type filterer interface {
	filter(line string) bool
}

type modifier interface {
	modify(line string) string
}

type rule struct {
	name       string
	pathFormat *regexp.Regexp
	modifiers  []modifier
	filters    []filterer
}

func newRule(name, pathFormat string, modifiers []modifier, filters []filterer) (*rule, error) {
	reg, err := regexp.Compile(pathFormat)
	if err != nil {
		return nil, fmt.Errorf("file format: %w", err)
	}

	if filters == nil || len(filters) == 0 {
		return nil, fmt.Errorf("at list one filter must be specified")
	}

	if modifiers == nil {
		modifiers = []modifier{}
	}

	return &rule{
		name:       name,
		pathFormat: reg,
		modifiers:  modifiers,
		filters:    filters,
	}, nil
}

func (r *rule) checkPath(filename string) bool {
	return r.pathFormat.MatchString(filename)
}

func (r *rule) checkLine(line string) bool {
	for _, m := range r.modifiers {
		line = m.modify(line)
	}

	ok := true
	for _, f := range r.filters {
		if !f.filter(line) {
			ok = false
			break
		}
	}

	return ok
}
