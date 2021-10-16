package main

import (
	"fmt"
	"regexp"
)

type rule struct {
	PathFormat *regexp.Regexp `json:"path_format"`
	Modifiers  []modifier     `json:"modifiers"`
	Filters    []filterer     `json:"filters"`
}

type rules map[string]*rule

func newRule(pathFormat string, modifiers []modifier, filters []filterer) (*rule, error) {
	reg, err := regexp.Compile(pathFormat)
	if err != nil {
		return nil, fmt.Errorf("file format: %w", err)
	}

	if len(filters) == 0 {
		return nil, fmt.Errorf("at list one filter must be specified")
	}

	if modifiers == nil {
		modifiers = []modifier{}
	}

	return &rule{
		PathFormat: reg,
		Modifiers:  modifiers,
		Filters:    filters,
	}, nil
}

func (r *rule) checkPath(filename string) bool {
	return r.PathFormat.MatchString(filename)
}

func (r *rule) checkLine(line string) bool {
	for _, m := range r.Modifiers {
		line = m.modify(line)
	}

	ok := true
	for _, f := range r.Filters {
		if !f.filter(line) {
			ok = false
			break
		}
	}

	return ok
}
