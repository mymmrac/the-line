package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"regexp"
)

type rule struct {
	PathFormat *regexp.Regexp `yaml:"path-format"`
	Modifiers  []modifier     `yaml:"modifiers"`
	Filters    []filterer     `yaml:"filters"`
}

type imap map[interface{}]interface{}

func parseFilter(kind string, value imap) (filterer, error) {
	var f filterer
	switch kind {
	case "any":
		f = &anyFilter{}
	case "blank":
		f = &blankFilter{}
	case "match":
		m := &matchFilter{}
		bytes, _ := yaml.Marshal(value)
		if err := yaml.Unmarshal(bytes, m); err != nil {
			return nil, err
		}
		f = m
	case "contains":
		m := &containsFilter{}
		bytes, _ := yaml.Marshal(value)
		if err := yaml.Unmarshal(bytes, m); err != nil {
			return nil, err
		}
		f = m
	case "length":
		m := &lengthFilter{}
		bytes, _ := yaml.Marshal(value)
		if err := yaml.Unmarshal(bytes, m); err != nil {
			return nil, err
		}
		f = m
	case "regexp":
		m := &regexpFilter{}
		bytes, _ := yaml.Marshal(value)
		if err := yaml.Unmarshal(bytes, m); err != nil {
			return nil, err
		}
		f = m
	case "multiline":
		m := &multilineFilter{}
		bytes, _ := yaml.Marshal(value)
		if err := yaml.Unmarshal(bytes, m); err != nil {
			return nil, err
		}
		f = m
	case "union":
		m := &unionFilter{}
		bytes, _ := yaml.Marshal(value)
		if err := yaml.Unmarshal(bytes, m); err != nil {
			return nil, err
		}
		f = m
	case "intersection":
		m := &intersectionFilter{}
		bytes, _ := yaml.Marshal(value)
		if err := yaml.Unmarshal(bytes, m); err != nil {
			return nil, err
		}
		f = m
	case "not":
		m := &intersectionFilter{}
		bytes, _ := yaml.Marshal(value)
		if err := yaml.Unmarshal(bytes, m); err != nil {
			return nil, err
		}
		f = m
	default:
		return nil, fmt.Errorf("unknown filter: %q", kind)
	}

	return f, nil
}

func (r *rule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ruleUnmarshal struct {
		PathFormat string `yaml:"path-format"`
		Modifiers  []imap `yaml:"modifiers"`
		Filters    []imap `yaml:"filters"`
	}

	u := ruleUnmarshal{}
	err := unmarshal(&u)
	if err != nil {
		return err
	}

	modifiers := make([]modifier, len(u.Modifiers))
	for i, m := range u.Modifiers {
		kind, ok := m["kind"]
		if !ok {
			return fmt.Errorf("no kind")
		}

		switch kind {
		case "trim-spaces":
			modifiers[i] = &trimSpaceModifier{}
		default:
			return fmt.Errorf("unknown modifier: %q", kind)
		}
	}

	filters := make([]filterer, len(u.Filters))
	for i, f := range u.Filters {
		kind, ok := f["kind"]
		if !ok {
			return fmt.Errorf("no kind")
		}

		filters[i], err = parseFilter(kind.(string), f)
		if err != nil {
			return err
		}
	}

	parsedRule, err := newRule(u.PathFormat, modifiers, filters)
	if err != nil {
		return err
	}

	r.PathFormat = parsedRule.PathFormat
	r.Modifiers = parsedRule.Modifiers
	r.Filters = parsedRule.Filters

	return nil
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
