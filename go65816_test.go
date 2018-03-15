// Test file for go65816
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 15. Mar 2018
// This version: 15. Mar 2018

package main

import "testing"

func TestIsValidAddr(t *testing.T) {
	var tests = []struct {
		input uint
		want  bool
	}{
		{0, true},
		{1 << 24, false},
		{1<<24 - 1, true},
	}

	for _, test := range tests {
		got := isValidAddr(test.input)
		if got != test.want {
			t.Errorf("isValidAddr(%q) = %v", test.input, got)
		}
	}
}
