// Test file for mem.og
// Part of the go65816 packages
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version 09. Mar 2018
// This version 09. Mar 2018

package mem

import "testing"

// Test if test for given address is valid
func TestIsValidAddr(t *testing.T) {
	var tests = []struct {
		input uint
		want  bool
	}{
		{0x0, true},
		{1024, true},
		{0xffffff, true},
		{0x1000000, false},
	}

	for _, test := range tests {
		if got := isValidAddr(test.input); got != test.want {
			t.Errorf("isValidAddr(%q) = %v", test.input, got)
		}
	}
}

// Test if we get the right size of a Chunk
func TestChunkSize(t *testing.T) {
	type ip struct {
		start uint
		end   uint
	}
	var tests = []struct {
		input ip
		want  uint
	}{
		{ip{0, 0}, 0},
		{ip{0, 1024}, 1024},
		{ip{1024, 2048}, 1024},
		{ip{0x100, 0x200}, 0x100},
	}

	for _, test := range tests {
		tc := Chunk{start: test.input.start, end: test.input.end}
		got := tc.size()
		if got != test.want {
			t.Errorf("Chunk size(%q) = %v", test.input, got)
		}
	}
}

// Test hexdump of data from chunk's memory
func TestChunkHexdump(t *testing.T) {
	tc := Chunk{start: 0, end: 10}
	tc.erase()
	tc.hexdump()
}
