package common

import (
	"log"
	"regexp"
)

type SymbolsTable interface {
	Exists(symbol string) bool
	GetValue(symbol string) (uint16, error)
}

var variableMatcher *regexp.Regexp = createVariableMatcher()

func createVariableMatcher() *regexp.Regexp {
	m, err := regexp.Compile("^[a-zA-Z_$][a-zA-Z_0-9]*$")
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func IsValidVariable(v string) bool {
	return variableMatcher.MatchString(v)
}
