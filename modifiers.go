package main

import (
	"strings"
)

type modifier interface {
	modify(line string) string
}

type trimSpaceModifier struct{}

func (m *trimSpaceModifier) modify(line string) string {
	return strings.TrimSpace(line)
}

type trimPrefixModifier struct {
	Prefix string `yaml:"prefix"`
}

func (m *trimPrefixModifier) modify(line string) string {
	return strings.TrimPrefix(line, m.Prefix)
}

type trimSuffixModifier struct {
	Suffix string `yaml:"suffix"`
}

func (m *trimSuffixModifier) modify(line string) string {
	return strings.TrimSuffix(line, m.Suffix)
}

type toLowerModifier struct{}

func (m *toLowerModifier) modify(line string) string {
	return strings.ToLower(line)
}
