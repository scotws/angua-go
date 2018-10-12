// Configuration loading routines for Angua
// Scot W. Stevenson scot.stevenson@gmail.com
// First version: 26. Sep 2017
// Second version: 15. Mar 2018

package config

import (
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
