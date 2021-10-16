package main

import "strings"

type modifier interface {
	modify(line string) string
}

type trimSpaceModifier struct{}

func (t trimSpaceModifier) modify(line string) string {
	return strings.TrimSpace(line)
}
