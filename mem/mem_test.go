// Test file for mem.og
// Part of the go65816 packages
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version 09. Mar 2018
// This version 15. Mar 2018

package mem

import "testing"

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

// Test if our address is in range
func TestContainsAddr(t *testing.T) {
	var (
		tc = Chunk{start: 0x400, end: 0x800}

		tests = []struct {
			input uint
			want  bool
		}{
			{0x400, true},
			{0x3FF, false},
			{0x801, false},
		}
	)

	for _, test := range tests {
		got := tc.containsAddr(test.input)

		if got != test.want {
			t.Errorf("Contains Addr(%q) = %v", test.input, got)
		}
	}
}

// Test fetching of a byte
func TestFetch(t *testing.T) {
	var (
		mydata = make([]byte, 0x400) // 1 KiB length

		tests = []struct {
			input uint
			want  byte
		}{
			{0x100, 0},
		}
	)
	tc := Chunk{start: 0, end: 0x400, label: "Test", data: mydata}

	for _, test := range tests {
		got := tc.fetch(test.input)

		if got != test.want {
			t.Errorf("Fetch (%q) = %v", test.input, got)
		}
	}
}

// Test storing of a byte
func TestStoreNFetch(t *testing.T) {
	var (
		mydata = make([]byte, 0x400) // 1 KiB buffer

		tests = []struct {
			addr uint
			b    byte
		}{
			{0x100, 0xEE},
		}
	)

	tc := Chunk{start: 0, end: 0x400, label: "Test", data: mydata}

	for _, test := range tests {

		tc.store(test.b, test.addr)
		got := tc.fetch(test.addr)

		if got != test.b {
			t.Errorf("Store and Fetch (%q) = %v", test.addr, test.b)
		}

		tc.hexdump()
	}
}
