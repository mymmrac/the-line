package main

import (
	"fmt"
	"regexp"
	"strings"
)

type filterer interface {
	filter(line string) bool
}

type anyFilter struct{}

func (f *anyFilter) filter(_ string) bool {
	return true
}

type blankFilter struct{}

func (f *blankFilter) filter(line string) bool {
	return line == ""
}

type matchFilter struct {
	Line string `yaml:"line"`
}

func (f *matchFilter) filter(line string) bool {
	return line == f.Line
}

type containsFilter struct {
	Line string `yaml:"line"`
}

func (f *containsFilter) filter(line string) bool {
	return strings.Contains(line, f.Line)
}

type prefixFilter struct {
	Prefix string `yaml:"prefix"`
}

func (p *prefixFilter) filter(line string) bool {
	return strings.HasPrefix(line, p.Prefix)
}

type suffixFilter struct {
	Suffix string `yaml:"suffix"`
}

func (s *suffixFilter) filter(line string) bool {
	return strings.HasSuffix(line, s.Suffix)
}

type lengthFilter struct {
	Length int `yaml:"length"`
}

func (f *lengthFilter) filter(line string) bool {
	return len(line) == f.Length
}

type regexpFilter struct {
	Reg *regexp.Regexp `yaml:"pattern"`
}

func (f *regexpFilter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type fu struct {
		Pattern string `yaml:"pattern"`
	}
	u := fu{}

	err := unmarshal(&u)
	if err != nil {
		return err
	}
	f.Reg, err = regexp.Compile(u.Pattern)
	return err
}

func (f *regexpFilter) filter(line string) bool {
	return f.Reg.MatchString(line)
}

type multilineFilter struct {
	StartFilter filterer `yaml:"start-filter"`
	EndFilter   filterer `yaml:"end-filter"`
	matching    bool     `yaml:"-"`
}

func (f *multilineFilter) filter(line string) bool {
	if f.StartFilter.filter(line) {
		f.matching = true
		return true
	}

	if f.EndFilter.filter(line) {
		f.matching = false
		return true
	}

	return f.matching
}

func (f *multilineFilter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type fu struct {
		StartFilter imap `yaml:"start-filter"`
		EndFilter   imap `yaml:"end-filter"`
	}
	u := fu{}

	err := unmarshal(&u)
	if err != nil {
		return err
	}

	kind, ok := u.StartFilter["kind"]
	if !ok {
		return fmt.Errorf("no kind for multiline filert start")
	}
	f.StartFilter, err = unmarshalYAMLFilter(kind, u.StartFilter)
	if err != nil {
		return err
	}

	kind, ok = u.EndFilter["kind"]
	if !ok {
		return fmt.Errorf("no kind for multiline filert end")
	}
	f.EndFilter, err = unmarshalYAMLFilter(kind, u.EndFilter)
	if err != nil {
		return err
	}

	return nil
}

type unionFilter struct {
	FilterA filterer `yaml:"filter-a"`
	FilterB filterer `yaml:"filter-b"`
}

func (f *unionFilter) filter(line string) bool {
	return f.FilterA.filter(line) || f.FilterB.filter(line)
}

func (f *unionFilter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type fu struct {
		FilterA imap `yaml:"filter-a"`
		FilterB imap `yaml:"filter-b"`
	}
	u := fu{}

	err := unmarshal(&u)
	if err != nil {
		return err
	}

	kind, ok := u.FilterA["kind"]
	if !ok {
		return fmt.Errorf("no kind for union filert a")
	}
	f.FilterA, err = unmarshalYAMLFilter(kind, u.FilterA)
	if err != nil {
		return err
	}

	kind, ok = u.FilterB["kind"]
	if !ok {
		return fmt.Errorf("no kind for union filert b")
	}
	f.FilterB, err = unmarshalYAMLFilter(kind, u.FilterB)
	if err != nil {
		return err
	}

	return nil
}

type intersectionFilter struct {
	FilterA filterer `yaml:"filter-a"`
	FilterB filterer `yaml:"filter-b"`
}

func (f *intersectionFilter) filter(line string) bool {
	return f.FilterA.filter(line) && f.FilterB.filter(line)
}

func (f *intersectionFilter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type fu struct {
		FilterA imap `yaml:"filter-a"`
		FilterB imap `yaml:"filter-b"`
	}
	u := fu{}

	err := unmarshal(&u)
	if err != nil {
		return err
	}

	kind, ok := u.FilterA["kind"]
	if !ok {
		return fmt.Errorf("no kind for intersection filert a")
	}
	f.FilterA, err = unmarshalYAMLFilter(kind, u.FilterA)
	if err != nil {
		return err
	}

	kind, ok = u.FilterB["kind"]
	if !ok {
		return fmt.Errorf("no kind for intersection filert b")
	}
	f.FilterB, err = unmarshalYAMLFilter(kind, u.FilterB)
	if err != nil {
		return err
	}

	return nil
}

type notFilter struct {
	Filter filterer `yaml:"filter"`
}

func (f *notFilter) filter(line string) bool {
	return !f.Filter.filter(line)
}

func (f *notFilter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type fu struct {
		Filter imap `yaml:"filter"`
	}
	u := fu{}

	err := unmarshal(&u)
	if err != nil {
		return err
	}

	kind, ok := u.Filter["kind"]
	if !ok {
		return fmt.Errorf("no kind for not filter")
	}
	f.Filter, err = unmarshalYAMLFilter(kind, u.Filter)

	return err
}
