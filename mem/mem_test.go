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

// ==== CHUNK TESTS ====

// Make sure bad chunk parameters are flagged during creation. Note these will
// print error strings.
func TestBadNewChunks(t *testing.T) {
	type params struct {
		start   common.Addr24
		end     common.Addr24
		memType string
	}
	var tests = []struct {
		input params
	}{
		// {params{1, 0, "ram"}},             // first addr must be smaller
		{params{0, 0, "ram"}}, // single-byte chunk is legal
		// {params{0, 0xffff, "frogbreath"}}, // need "rom" or "ram"
	}

	for _, test := range tests {
		_, err := NewChunk(test.input.start, test.input.end, test.input.memType)

		if err != nil {
			t.Errorf("Bad NewChunk test(%q) = %v", test.input, err)
		}
	}
}

// Test if we get the right size of a Chunk
func TestChunkSize(t *testing.T) {
	type params struct {
		start   common.Addr24
		end     common.Addr24
		memType string
	}
	var tests = []struct {
		input params
		want  uint
	}{
		{params{0, 0, "ram"}, 1}, // chunks of size 1 are legal
		{params{0, 0x3FF, "ram"}, 1024},
		{params{0x400, 0x7FF, "ram"}, 1024},
		{params{0x100, 0x1FF, "ram"}, 0x100},
		{params{0x100100, 0x1001FF, "ram"}, 0x100},
	}

	for _, test := range tests {
		tc, _ := NewChunk(test.input.start, test.input.end, test.input.memType)
		got := tc.size()

		if got != test.want {
			t.Errorf("chunk.Size(%q) = %v", test.input, got)
		}
	}
}

// Test if our address is in range in a chunk
func TestChunkContainsAddr(t *testing.T) {
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
		got := tc.contains(test.input)

		if got != test.want {
			t.Errorf("Contains Addr(%q) = %v", test.input, got)
		}
	}
}

// Test fetching of a byte from a chunk. Note that chunk.Fetch does not test if
// value is legal and does not return a flag
func TestChunkFetch(t *testing.T) {
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
		got := tc.fetch(test.input)

		if got != test.want {
			t.Errorf("Fetch (%q) = %v", test.input, got)
		}
	}
}

// Test storing of a byte in a chunk. Assumes this is RAM.
func TestChunkStoreNFetch(t *testing.T) {
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

		tc.store(test.addr, test.b)
		got := tc.fetch(test.addr)
		if got != test.b {
			t.Errorf("Store and Fetch (%q) = %v", test.addr, test.b)
		}
	}
}

// Test storing of a multi-byte number in little-endian format
func TestChunkStoreMore(t *testing.T) {
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
	}{
		// {ip{0x100, 0xee, 0}}, // can't ask to store zero bytes
		// {ip{0x100, 0xee, 4}}, // can't ask to store four bytes
		{ip{0x100, 0xaa, 1}},
		{ip{0x100, 0xaabb, 2}},
		{ip{0x100, 0xaabbcc, 3}},
	}

	for _, test := range tests {
		err := mymem.StoreMore(test.input.addr, test.input.num, test.input.len)

		if err != nil {
			t.Errorf("StoreMore(%v) = %v", test.input, err)
		}

	}
}

// TODO test Chunk.StoreBlock

// ==== MEMORY TESTS ====

// TODO test Memory.Contains
// TODO test Memory.Store
// TODO test Memory.Fetch
// TODO test Memory.StoreBlock
