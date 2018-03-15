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
	start     uint   // stores 65816 addr
	end       uint   // stores 65816 addr
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
	return c.data[addr-c.start]
}

// hexdump prints the chunk's memory contents in a nice hex table
// We could use the library encoding/hex for this, but we want to print the
// first address of the line, and the library function starts the count with
// zero, not the address. Also, we want uppercase letters for hex values
func (c Chunk) hexdump(addr1, addr2 uint) {

	var r rune
	var count uint
	var hb strings.Builder // hex part
	var cb strings.Builder // char part
	var template string = "%-58s%s\n"

	if !c.contains(addr1) {
		fmt.Printf("Address %X not in chunk %s", addr1, c.label)
		return
	}

	if !c.contains(addr2) {
		fmt.Printf("Address %X not in chunk %s", addr2, c.label)
		return
	}

	for i := addr1; i < addr2; i++ {

		// The first run produces a blank line because this if is
		// triggered, however, the strings are empty because of the way
		// Go initializes things
		if count%16 == 0 {
			fmt.Printf(template, hb.String(), cb.String())
			hb.Reset()
			cb.Reset()

			fmt.Fprintf(&hb, "%06X ", addr1+count)
		}

		b := c.fetch(i)

		// Build the hex string
		fmt.Fprintf(&hb, " %02X", b)

		// Build the string list. This is the 21. century so we hexdump
		// in Unicode, not ASCII
		r = rune(b)
		if !unicode.IsPrint(r) {
			r = rune('.')
		}

		fmt.Fprintf(&cb, string(r))
		count += 1

		// We put one extra blank line after the first eight entries to
		// make the dump more readable
		if count%8 == 0 {
			fmt.Fprintf(&hb, " ")
		}

	}

	// If the loop is all done, we might still have stuff left in the
	// buffers
	fmt.Printf(template, hb.String(), cb.String())
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
