// Angua Memory System
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 09. Mar 2018
// This version: 31. Dec 2018

package mem

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"unicode"

	"angua/common"
)

// Chunks are the basic memory unit in Angua, a continuous region of memory that
// can either be read-only (ROM) or read and write (RAM). The memory contents
// itself is stored as a list of bytes
type Chunk struct {
	Start      common.Addr24 // stores 65816 addr
	End        common.Addr24 // stores 65816 addr
	Type       string        // "ram" or "rom"
	sync.Mutex               // Make chunks threadsafe
	Data       []byte
}

// Memory is the total system memory, which is just a list of chunks
type Memory struct {
	Chunks []Chunk
}

// --- CHUNK METHODS ---

// Contains takes a memory address and checks if it is in this chunk,
// returning a bool
func (c Chunk) Contains(addr common.Addr24) bool {
	return c.Start <= addr && addr <= c.End
}

// Fetch gets one byte of memory from the data of a chunk and returns it.
// Assumes we have already made sure that the address is in this chunk
func (c Chunk) Fetch(addr common.Addr24) byte {
	c.Lock()
	b := c.Data[addr-c.Start]
	c.Unlock()
	return b
}

// Size returns the, uh, size of a chunk in bytes as an uint. Does not check to
// see if chunk addresses are valid, that is, c.End is larger than c.Start
// TODO see if this should be an int
func (c Chunk) Size() uint {
	return uint(c.End - c.Start + 1)
}

// Store takes a byte and an address and stores the byte at the address in the
// chunk. Assumes that we already checked that the address is in fact in this
// chunk.
func (c Chunk) Store(addr common.Addr24, b byte) {
	c.Lock()
	c.Data[addr-c.Start] = b
	c.Unlock()
}

// --- MEMORY METHODS ---

// Contains takes an 65816 address and checks to see if it is
// valid, returning a bool
func (m Memory) Contains(addr common.Addr24) bool {
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
// returns a zero value and a false flag. We use a Mutex at chunk level, not
// memory level
func (m Memory) Fetch(addr common.Addr24) (byte, bool) {
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
// -- to be fetched and returned as an integer, retrieving those bytes as
// little endian. Also returns a bool to show if all fetches were to legal
// addresses. Assumes that the address itself was vetted.
func (m Memory) FetchMore(addr common.Addr24, num uint) (uint, bool) {
	const maxint uint = 3
	var legal bool = true
	var sum uint

	// This is a bit harsh, but it will help us find errors during the
	// development process. Consider doing something less drastic later
	if num > maxint {
		log.Fatal(fmt.Sprintf("Illegal attempt to read %d bytes from %06X", num, addr))
	}

	for i := common.Addr24(0); i <= common.Addr24(num-1); i++ {
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
// TODO Move this to the interface code and just take an address
func (m Memory) Hexdump(addr1, addr2 common.Addr24) {
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

			// TODO move this to common print function
			fmt.Fprintf(&hb, "%06X ", addr1+common.Addr24(count))
		}

		// We ignore the ok flag here
		b, ok := m.Fetch(i)
		if !ok {
			log.Fatal("ERROR fetching byte", i, "from memory")
		}

		// Build the hex string
		fmt.Fprintf(&hb, " %02X", b)

		// Build the string list. This is the 21. century so we hexdump
		// in Unicode, not ASCII, though this doesn't make a different
		// if we just have byte values
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
	var template string = "%s %s %s (%d bytes)\n"

	for _, c := range m.Chunks {
		r += fmt.Sprintf(template,
			c.Start.HexString(), c.End.HexString(),
			strings.ToUpper(c.Type), c.Size())
	}

	if r == "" {
		r = "No memory defined. Use 'memory' command."
	}

	return r
}

// Read returns a slice of memory as bytes and a flag to show if all bytes were
// part of legal memory when given a starting address and the a size. Assumes
// that the addresses have been correctly formatted and vetted
func (m Memory) Read(addr common.Addr24, size uint) ([]byte, bool) {

	var allLegal bool = true
	var bs []byte

	for i := addr; i <= addr+common.Addr24(size); i++ {
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
func (m Memory) Store(addr common.Addr24, b byte) bool {
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
func (m Memory) StoreMore(addr common.Addr24, num uint, len uint) bool {

	if len < 1 || len > 3 {
		return false
	}

	f := true

	lsb := byte(num & 0xff)
	msb := byte((num & 0xff00) >> 8)
	bank := byte((num & 0xff0000) >> 16)

	ba := [...]byte{lsb, msb, bank}

	for i := 0; i < 3; i++ {
		ok := m.Store(addr+common.Addr24(i), ba[i])
		f = f && ok
	}

	return f
}

// Write takes a 65816 address and a slice of bytes and stores those bytes
// starting at that address. If all addresses were legal, it returns a true
// flag, otherwise a false
func (m Memory) Write(addr common.Addr24, bs []byte) bool {
	var legal bool = true

	for i := addr; i < addr+common.Addr24(len(bs)); i++ {
		ok := m.Store(i, bs[i-addr])

		if !ok {
			legal = false
		}
	}
	return legal
}
