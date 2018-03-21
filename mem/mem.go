// go65816 Memory System
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 09. Mar 2018
// This version: 15. Mar 2018

package mem

import (
	"fmt"
	"log"
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

// Size returns the, uh, size of a chunk in bytes
func (c Chunk) Size() uint {
	return c.End - c.Start + 1
}

// Store takes a byte and an address and stores the byte at the address in the
// chunk. Assumes that we already checked that the address is in fact in this
// chunk
func (c Chunk) Store(addr uint, b byte) {
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
// the byte with a true flag for success. If the address is not memory, it
// returns a zero value and a false flag
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

// FetchMore takes a 65816 address and the number of bytes to get -- 1, 2 or 3
// -- should be fetched and returned as an integer, retrieving those bytes as
// little endian. Also returns a bool to show if all fetches were to legal
// addresses. Assumes that the address itself was vetted.
func (m Memory) FetchMore(addr uint, num uint) (uint, bool) {

	const maxint = 3
	var legal bool = true
	var sum uint

	// This is a bit harsh, but it will help us find errors during the
	// development process. Consider doing something less drastic later
	if num > maxint {
		log.Fatal(fmt.Sprintf("Illegal attempt to read %d bytes from %06X", num, addr))
	}

	for i := uint(0); i <= num-1; i++ {
		b, ok := m.Fetch(addr + i)

		if !ok {
			legal = false
		}

		// Shift eight bits to the left for every byte we go further to
		// the right
		sum += (uint(b) << (8 * i))
	}
	return sum, legal
}

// Hexdump prints the contents of a memory range in a nice hex table. If the
// addresses do not exist, we just print a zero without any fuss. We could use
// the library encoding/hex for this, but we want to print the first address of
// the line, and the library function starts the count with zero, not the
// address. Also, we want uppercase letters for hex values
// TODO Return result as a string instead of printing it
func (m Memory) Hexdump(addr1, addr2 uint) {

	var r rune
	var count uint
	var hb strings.Builder // hex part
	var cb strings.Builder // char part
	var template string = "%-58s%s\n"

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

		// We ignore the ok flag here
		b, _ := m.Fetch(i)

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

// List returns a list of all chunks in memory, as a string
func (m Memory) List() string {
	var r string
	var template string = "%s %s %06X-%06X (%d bytes)\n"

	for _, c := range m.Chunks {
		r += fmt.Sprintf(template, c.Label, c.Type, c.Start, c.End, c.Size())
	}

	return r
}

// Read returns a slice of memory as bytes and a flag to show if all bytes were
// part of legal memory when given a starting address and the a size. Assumes
// that the addresses have been correctly formatted and vetted
func (m Memory) Read(addr uint, size uint) ([]byte, bool) {

	var allLegal bool = true
	var bs []byte

	for i := addr; i <= addr+size; i++ {
		b, ok := m.Fetch(i)

		if !ok {
			allLegal = false
		}
		bs = append(bs, b)
	}
	return bs, allLegal
}

// Size returns the total size of the system memory, RAM and ROM, in bytes
func (m Memory) Size() uint {

	var sum uint

	for _, c := range m.Chunks {
		sum += c.Size()
	}

	return sum
}

// Store takes an address and a byte and saves them to memory. If the addr is
// not part of legal memory, we return a false flag, otherwise a true. If the
// addr is part of ROM, do the same thing
func (m Memory) Store(addr uint, b byte) bool {
	var f bool = false

	for _, c := range m.Chunks {

		if c.Type == "ram" && c.Contains(addr) {
			c.Store(addr, b)
			f = true
			break
		}
	}
	return f
}

// StoreMore takes an address, a number and the number of bytes to store little
// endian at that address. If any part of the address is not a part of legal
// memory, return a false flag, otherwise true. At most, numbers up to 24 bit
// length (three bytes) are stored. Anything above that is silently discarded.
// If the number of bytes to store is anything but 1, 2, or 3, we return a false
// flag with memory untouched
func (m Memory) StoreMore(addr uint, num uint, len uint) bool {

	if len < 1 || len > 3 {
		return false
	}

	f := true

	lsb := byte(num & 0xff)
	msb := byte((num & 0xff00) >> 8)
	bank := byte((num & 0xff0000) >> 16)

	ba := [...]byte{lsb, msb, bank}

	for i := 0; i < 3; i++ {
		ok := m.Store(addr+uint(i), ba[i])
		f = f && ok
	}

	return f
}

// Write takes a 65816 address and a slice of bytes and story those bytes start
// at that address. If all addresses were legal, it returns a true flag,
// otherwise a false
func (m Memory) Write(addr uint, bs []byte) bool {
	var legal bool = true

	for i := addr; i < addr+uint(len(bs)); i++ {
		ok := m.Store(i, bs[i-addr])

		if !ok {
			legal = false
		}
	}
	return legal
}
