// Test file for config.go
// Part of the Rlyeh package
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version 26. Sep 2017
// This version 15. Mar 2018

package config

import "testing"

func TestIsComment(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"", false},      // empty line
		{"#", true},      // comment at beginning of line
		{"# ----", true}, // comment with stuff
		{" #", true},     // comment after indent
	}

	for _, test := range tests {
		if got := IsComment(test.input); got != test.want {
			t.Errorf("IsComment(%q) = %v", test.input, got)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"tali", false},
		{"", true},
		{" ", true},
		{"\t", true},
	}

	for _, test := range tests {
		if got := IsEmpty(test.input); got != test.want {
			t.Errorf("IsEmpty(%q) = %v", test.input, got)
		}
	}
}
