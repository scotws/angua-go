// Test file for mem.go
// Part of the Angua package
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version 09. Mar 2018
// This version 01. Jan 2019

package mem

import (
	"testing"

	"angua/common"
)

// ==== TESTING CHUNKS ====

// Test if we get the right size of a Chunk
func TestChunkSize(t *testing.T) {
	type ip struct {
		start common.Addr24
		end   common.Addr24
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

// Test if our address is in range in a chunk
func TestContainsAddr(t *testing.T) {
	var (
		tc = Chunk{Start: 0x400, End: 0x7FF}

		tests = []struct {
			input common.Addr24
			want  bool
		}{
			{0x100, false},
			{0x3FF, false},
			{0x400, true},
			{0x500, true},
			{0x7FF, true},
			{0x800, false},
		}
	)

	for _, test := range tests {
		got := tc.Contains(test.input)

		if got != test.want {
			t.Errorf("Contains Addr(%q) = %v", test.input, got)
		}
	}
}

// Test fetching of a byte from a chunk. Note that chunk.Fetch does not test if
// value is legal and does not return a flag
func TestFetch(t *testing.T) {
	var (
		mydata = make([]byte, 0x400) // 1 KiB length

		tests = []struct {
			input common.Addr24
			want  byte
		}{
			{0x100, 0},
		}
	)
	tc := Chunk{Start: 0x100, End: 0x5FF, Data: mydata}

	for _, test := range tests {
		got := tc.Fetch(test.input)

		if got != test.want {
			t.Errorf("Fetch (%q) = %v", test.input, got)
		}
	}
}

// Test storing of a byte in a chunk
func TestStoreNFetch(t *testing.T) {
	var (
		mydata = make([]byte, 0x400) // 1 KiB buffer

		tests = []struct {
			addr common.Addr24
			b    byte
		}{
			{0x100, 0xEE},
			{0x400, 0xEE},
		}
	)

	tc := Chunk{Start: 0x100, End: 0x4FF, Data: mydata}

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
		addr common.Addr24
		num  uint
		len  uint
	}

	var mydata = make([]byte, 0x400) // 1 KiB length
	var mymem Memory
	tc := Chunk{Start: 0x100, End: 0x4FF, Data: mydata, Type: "ram"}
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
	//	mymem.Hexdump(0x100, 0x4FF)
}
