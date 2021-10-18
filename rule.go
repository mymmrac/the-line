package main

import (
	"fmt"
	"regexp"

	"gopkg.in/yaml.v2"
)

type rule struct {
	PathFormat *regexp.Regexp `yaml:"path-format"`
	Modifiers  []modifier     `yaml:"modifiers"`
	Filters    []filterer     `yaml:"filters"`
}

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

type rules map[string]*rule

type imap map[interface{}]interface{}

func (r *rule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type uRule struct {
		PathFormat string `yaml:"path-format"`
		Modifiers  []imap `yaml:"modifiers"`
		Filters    []imap `yaml:"filters"`
	}

	u := uRule{}
	err := unmarshal(&u)
	if err != nil {
		return err
	}

	modifiers := make([]modifier, len(u.Modifiers))
	for i, m := range u.Modifiers {
		kind, ok := m["kind"]
		if !ok {
			return fmt.Errorf("no kind for modifier")
		}

		modifiers[i], err = unmarshalYAMLModifier(kind, m)
		if err != nil {
			return err
		}
	}

	filters := make([]filterer, len(u.Filters))
	for i, f := range u.Filters {
		kind, ok := f["kind"]
		if !ok {
			return fmt.Errorf("no kind for filter")
		}

		filters[i], err = unmarshalYAMLFilter(kind, f)
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

func unmarshalYAMLModifier(kind interface{}, value imap) (modifier, error) {
	var m modifier

	bytes, err := yaml.Marshal(value)
	if err != nil {
		return nil, err
	}

	switch kind {
	case "trim-spaces":
		m = &trimSpaceModifier{}
	case "trim-prefix":
		u := &trimPrefixModifier{}
		if err = yaml.Unmarshal(bytes, u); err != nil {
			return nil, err
		}
		m = u
	case "trim-suffix":
		u := &trimSuffixModifier{}
		if err = yaml.Unmarshal(bytes, u); err != nil {
			return nil, err
		}
		m = u
	case "to-lower":
		m = &toLowerModifier{}
	default:
		return nil, fmt.Errorf("unknown modifier: %q", kind)
	}

	return m, nil
}

//nolint:funlen,gocyclo,gocognit
func unmarshalYAMLFilter(kind interface{}, value imap) (filterer, error) {
	var f filterer

	bytes, err := yaml.Marshal(value)
	if err != nil {
		return nil, err
	}

	switch kind {
	case "any":
		f = &anyFilter{}
	case "blank":
		f = &blankFilter{}
	case "match":
		u := &matchFilter{}
		if err = yaml.Unmarshal(bytes, u); err != nil {
			return nil, err
		}
		f = u
	case "contains":
		u := &containsFilter{}
		if err = yaml.Unmarshal(bytes, u); err != nil {
			return nil, err
		}
		f = u
	case "prefix":
		u := &prefixFilter{}
		if err = yaml.Unmarshal(bytes, u); err != nil {
			return nil, err
		}
		f = u
	case "suffix":
		u := &suffixFilter{}
		if err = yaml.Unmarshal(bytes, u); err != nil {
			return nil, err
		}
		f = u
	case "length":
		u := &lengthFilter{}
		if err = yaml.Unmarshal(bytes, u); err != nil {
			return nil, err
		}
		f = u
	case "regexp":
		u := &regexpFilter{}
		if err = yaml.Unmarshal(bytes, u); err != nil {
			return nil, err
		}
		f = u
	case "multiline":
		u := &multilineFilter{}
		if err = yaml.Unmarshal(bytes, u); err != nil {
			return nil, err
		}
		f = u
	case "union":
		u := &unionFilter{}
		if err = yaml.Unmarshal(bytes, u); err != nil {
			return nil, err
		}
		f = u
	case "intersection":
		u := &intersectionFilter{}
		if err = yaml.Unmarshal(bytes, u); err != nil {
			return nil, err
		}
		f = u
	case "not":
		u := &notFilter{}
		if err = yaml.Unmarshal(bytes, u); err != nil {
			return nil, err
		}
		f = u
	default:
		return nil, fmt.Errorf("unknown filter: %q", kind)
	}

	return f, nil
}
