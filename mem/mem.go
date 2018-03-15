// go65816 Memory System
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 09. Mar 2018
// This version: 15. Mar 2018

package mem

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	maxAddr = 1<<24 - 1
)

type Chunk struct {
	start     uint
	end       uint
	writeable bool   // true for RAM, false for ROM
	file      string // file path
	label     string // internal use only
	data      []byte
}

// contains takes a memory address and checks if it is in this chunk,
// returning a bool. Assumes that the address has been confirmed to be a valid
// 65816 address as a uint
func (c Chunk) contains(addr uint) bool {
	return c.start <= addr && addr <= c.end
}

// fetch gets one byte of memory from the data of a chunk and returns it.
// Assumes we have already made sure that the address is in this chunk
func (c Chunk) fetch(addr uint) byte {
	index := addr - c.start
	return c.data[index]
}

// hexdump prints the chunk's memory contents in a nice hex table
// We could use the library encoding/hex for this, but we want to print the
// first address of the line, and the library function starts the count with
// zero, not the address. Also, we want uppercase letters for hex values
func (c Chunk) hexdump() {

	var r rune
	var count uint = 0
	var hb strings.Builder // hex part
	var sb strings.Builder // char part

	for _, b := range c.data {

		// The first run produces a blank line because this if is
		// triggered, however, the strings are empty because of the way
		// Go initializes things
		if count%16 == 0 {
			fmt.Print(hb.String())
			fmt.Println(" ", sb.String())
			hb.Reset()
			sb.Reset()

			fmt.Fprintf(&hb, "%06X ", c.start+count)
		}

		// Build the hex string
		fmt.Fprintf(&hb, " %02X", b)

		// Build the string list. This is the 21. century so we hexdump
		// in Unicode, not ASCII
		r = rune(b)
		if !unicode.IsPrint(r) {
			r = rune('.')
		}

		fmt.Fprintf(&sb, string(r))
		count += 1

		// We put one extra blank line after the first eight entries to
		// make the dump more readable
		if count%8 == 0 {
			fmt.Fprintf(&hb, " ")
		}
	}
}

// Size returns the, uh, size of a chunk in bytes
func (c Chunk) size() uint {
	return c.end - c.start
}

// Store takes a byte and an address and stores the byte at the address in the
// chunk. Assumes that we already checked that the address is in fact in this
// chunk
func (c Chunk) store(b byte, addr uint) {
	c.data[addr-c.start] = b
}

// Memory is the total system memory, which is basically just a bunch of chunks
type Memory struct {
	chunks []Chunk
}

// contains takes an 65816 address as an uint and checks to see if it is
// valid, returning a bool
func (m Memory) contains(addr uint) bool {

	result := false

	for _, c := range m.chunks {

		if c.contains(addr) {
			result = true
			break
		}
	}

	return result
}

// size returns the total size of the system memory, RAM and ROM, in bytes
func (m Memory) size() uint {

	var sum uint = 0

	for _, c := range m.chunks {
		sum += c.size()
	}

	return sum
}
