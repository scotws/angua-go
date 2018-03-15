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
	Start uint   // stores 65816 addr
	End   uint   // stores 65816 addr
	Type  string // "ram" or "rom"
	Label string // internal use only
	Data  []byte
}

// Memory is the total system memory, which is basically just a bunch of chunks
type Memory struct {
	Chunks []Chunk
}

// --- CHUNK METHODS ---

// contains takes a memory address and checks if it is in this chunk,
// returning a bool. Assumes that the address has been confirmed to be a valid
// 65816 address as a uint
func (c Chunk) Contains(addr uint) bool {
	return c.Start <= addr && addr <= c.End
}

// fetch gets one byte of memory from the data of a chunk and returns it.
// Assumes we have already made sure that the address is in this chunk
func (c Chunk) Fetch(addr uint) byte {
	return c.Data[addr-c.Start]
}

// hexdump prints the chunk's memory contents in a nice hex table
// We could use the library encoding/hex for this, but we want to print the
// first address of the line, and the library function starts the count with
// zero, not the address. Also, we want uppercase letters for hex values
func (c Chunk) Hexdump(addr1, addr2 uint) {

	var r rune
	var count uint
	var hb strings.Builder // hex part
	var cb strings.Builder // char part
	var template string = "%-58s%s\n"

	if !c.Contains(addr1) {
		fmt.Printf("Address %X not in chunk %s", addr1, c.Label)
		return
	}

	if !c.Contains(addr2) {
		fmt.Printf("Address %X not in chunk %s", addr2, c.Label)
		return
	}

	for i := addr1; i <= addr2; i++ {

		// The first run produces a blank line because this if is
		// triggered, however, the strings are empty because of the way
		// Go initializes things
		if count%16 == 0 {
			fmt.Printf(template, hb.String(), cb.String())
			hb.Reset()
			cb.Reset()

			fmt.Fprintf(&hb, "%06X ", addr1+count)
		}

		b := c.Fetch(i)

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
func (c Chunk) Size() uint {
	return c.End - c.Start + 1
}

// Store takes a byte and an address and stores the byte at the address in the
// chunk. Assumes that we already checked that the address is in fact in this
// chunk
func (c Chunk) Store(b byte, addr uint) {
	c.Data[addr-c.Start] = b
}

// --- MEMORY METHODS ---

// contains takes an 65816 address as an uint and checks to see if it is
// valid, returning a bool
func (m Memory) Contains(addr uint) bool {

	result := false

	for _, c := range m.Chunks {

		if c.Contains(addr) {
			result = true
			break
		}
	}
	return result
}

// Fetch takes an address and gets a byte from the appropriate chunk and returns
// the byte with a true flag for success. If the byte is not memory, it returns
// a zero and a false flag
func (m Memory) Fetch(addr uint) (byte, bool) {
	var b byte
	var found bool = false

	for _, c := range m.Chunks {

		if c.Contains(addr) {
			b = c.Fetch(addr)
			found = true
			break
		}
	}
	return b, found
}

// List returns a list of all chunks in memory, as a string
func (m Memory) List() string {
	var r string
	var template string = "%s %s %06X-%06X (%d bytes)\n"

	for _, c := range m.Chunks {
		r += fmt.Sprintf(template, c.Label, c.Type, c.Start, c.End, c.Size())
	}

	return r
}

// Size returns the total size of the system memory, RAM and ROM, in bytes
func (m Memory) Size() uint {

	var sum uint

	for _, c := range m.Chunks {
		sum += c.Size()
	}

	return sum
}
