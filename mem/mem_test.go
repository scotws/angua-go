// Test file for mem.og
// Part of the go65816 packages
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version 09. Mar 2018
// This version 15. Mar 2018

package mem

import (
	"testing"
)

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
		{ip{0, 0}, 1},
		{ip{0, 0x3FF}, 1024},
		{ip{0x400, 0x7FF}, 1024},
		{ip{0x100, 0x1FF}, 0x100},
	}

	for _, test := range tests {
		tc := Chunk{Start: test.input.start, End: test.input.end}
		got := tc.Size()
		if got != test.want {
			t.Errorf("Chunk size(%q) = %v", test.input, got)
		}
	}
}

// Test if our address is in range
func TestContainsAddr(t *testing.T) {
	var (
		tc = Chunk{Start: 0x400, End: 0x800}

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
		got := tc.Contains(test.input)

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
	tc := Chunk{Start: 0, End: 0x400, Label: "Test", Data: mydata}

	for _, test := range tests {
		got := tc.Fetch(test.input)

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
			{0x400, 0xEE},
		}
	)

	tc := Chunk{Start: 0x100, End: 0x500, Label: "Test", Data: mydata}

	for _, test := range tests {

		tc.Store(test.addr, test.b)
		got := tc.Fetch(test.addr)

		if got != test.b {
			t.Errorf("Store and Fetch (%q) = %v", test.addr, test.b)
		}
	}
}

// Test storing of a multi-byte number in little-endian format
func TestStoreMore(t *testing.T) {
	type ip struct {
		addr uint
		num  uint
		len  uint
	}

	var mydata = make([]byte, 0x400) // 1 KiB length
	var mymem Memory
	tc := Chunk{Start: 0x100, End: 0x4FF, Label: "Test", Data: mydata, Type: "ram"}
	mymem.Chunks = append(mymem.Chunks, tc)

	var tests = []struct {
		input ip
		want  bool
	}{
		{ip{0x100, 0xee, 0}, false}, // can't ask to store zero bytes
		{ip{0x100, 0xee, 4}, false}, // can't ask to store four bytes
		{ip{0x100, 0xaa, 1}, true},
		{ip{0x100, 0xaabb, 2}, true},
		{ip{0x100, 0xaabbcc, 3}, true},
	}

	for _, test := range tests {
		got := mymem.StoreMore(test.input.addr, test.input.num, test.input.len)

		if got != test.want {
			t.Errorf("StoreMore(%v) = %v", test.input, got)
		}

	}
}
