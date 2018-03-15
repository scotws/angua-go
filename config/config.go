// Configuration loading routines for go65816
// Part of the go65816 package
// Scot W. Stevenson scot.stevenson@gmail.com
// First version: 26. Sep 2017
// Second version: 15. Mar 2018

package config

import (
	"fmt"
	"log"
	"strings"
)

const (
	keywordChunk   = "chunk"
	keywordSpecial = "special"
)

func IsComment(s string) bool {
	cs := strings.TrimSpace(s)
	return strings.HasPrefix(cs, "#")
}

func IsEmpty(s string) bool {
	cs := strings.TrimSpace(s)
	return cs == ""
}

func DefinesSpecial(s string) bool {
	return s == keywordSpecial
}

func DefinesChunk(s string) bool {
	return s == keywordChunk
}

func IsWriteable(s string) bool {
	var r bool

	switch {
	case s == "ram":
		r = true
	case s == "rom":
		r = false
	default:
		log.Fatal(fmt.Sprintf("Unknown keyword '%s' in chunk config file", s))
	}

	return r
}
