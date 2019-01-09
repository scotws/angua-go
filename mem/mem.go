// Angua Memory System
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 09. Mar 2018
// This version: 06. Jan 2019

package mem

import (
	"fmt"
	"strings"
	"sync"

	"angua/common"
)

// Chunk is the basic memory unit in Angua, a continuous region of memory that
// can either be read-only (ROM) or read and write (RAM). The memory contents
// itself is stored as a list of bytes
type Chunk struct {
	Start      common.Addr24 // stores 65816 addr
	End        common.Addr24 // stores 65816 addr
	Type       string        // "ram" or "rom"
	sync.Mutex               // Make chunks threadsafe
	Data       []byte
}

// Memory is the total system memory, which is just a list of chunks. This is
// why there is no function NewMemory to parallel NewChunk. We include the lists
// of special addresses here
// TODO add a map of addresses to bools to memoize checks if the address is
// present in memory; see Memory.Contains()
type Memory struct {
	Chunks    []Chunk
	SpecRead  map[common.Addr24]func() (byte, bool)
	SpecWrite map[common.Addr24]func(byte)
}

// NewChunk takes the start and end address for a new chunk as well as its type
// and returns a new chunk with initialized memory and a bool that describes
// success or failure. This routine checks to make sure that the addresses are
// sane, and the strings are correct. Errors are handled here by printing to the
// log. This is the only way that new chunks should be created.
func NewChunk(addr1, addr2 common.Addr24, cType string) (Chunk, error) {

	// Limit size of addresses to 24 bit. We don't have to check if this is
	// an unsigned int because the type system does that for us, see
	// common.Addr24
	addr1ok := common.Ensure24(addr1)
	addr2ok := common.Ensure24(addr2)

	// Make sure addr1 is really smaller than addr2. We don't accept chunks
	// with a length of zero bytes. We do, however, accept a chunk with one
	// byte, which is what you get when addr1 and addr2 are the same.
	if addr2ok < addr1ok {
		return Chunk{}, fmt.Errorf("mem.NewChunk: Invalid addresses for new chunk")
	}

	// Make sure memType is either "ram" or "rom"
	if cType != "ram" && cType != "rom" {
		return Chunk{}, fmt.Errorf("ERROR: Chunk type must either be 'ram' or 'rom', not %s", cType)
	}

	// The tendency of computer people to count from 0 makes the size a bit
	// more difficult. If we start at the address 0 and end at (15), then a
	// naive calculation of 15-0 gives us 15 addresses, though of course we
	// have 16. For this reason, we have to add one futher byte by hand.
	size := (addr2ok - addr1ok) + 1

	nc := Chunk{addr1ok, addr2ok, cType, sync.Mutex{}, make([]byte, size)}

	return nc, nil
}

// --- CHUNK METHODS ---

// Chunk methods are only to be used internally; the other parts of Angua
// interact through the Memory on a higher level.

// Contains takes a memory address and checks if it is in this chunk,
// returning a bool
func (c Chunk) contains(addr common.Addr24) bool {
	return c.Start <= addr && addr <= c.End
}

// Fetch gets one byte of memory from the data of a chunk and returns it.
// Assumes we have already made sure that the address is in this chunk
func (c Chunk) fetch(addr common.Addr24) byte {
	c.Lock()
	b := c.Data[addr-c.Start]
	c.Unlock()
	return b
}

// Size returns the size of a chunk in bytes as an uint. Does not check to
// see if chunk addresses are valid, that is, c.End is larger than c.Start
func (c Chunk) size() uint {
	return uint(c.End - c.Start + 1)
}

// Store takes a byte and an address and stores the byte at the address in the
// chunk. Assumes that we already checked that the address is in fact in this
// chunk. Note this doesn't care if memory is RAM or ROM, this is handled at the
// Memory level. This way, we can use this word for both Write and Burn routines
// TODO see if we need the lock
func (c Chunk) store(addr common.Addr24, b byte) {
	c.Lock()
	c.Data[addr-c.Start] = b
	c.Unlock()
}

// --- MEMORY METHODS ---

// Contains takes an 65816 address and checks to see if it is valid, returning a
// bool
// TODO we are going to be checking the same addresses over and over again, so we
// should speed this up by including a map of addresses to bools as a
// memoization device in the Memory struct.
func (m Memory) Contains(addr common.Addr24) bool {
	result := false

	for _, c := range m.Chunks {

		if c.contains(addr) {
			result = true
			break
		}
	}
	return result
}

// Fetch takes an address and gets a byte from the appropriate chunk and returns
// it with a true flag for success. If the address is not memory, it
// returns a zero value and an error. We use a Mutex at chunk level, not
// memory level. This is the basic fetch operation
func (m Memory) Fetch(addr common.Addr24) (byte, error) {
	var b byte
	var found bool

	// First, see if this triggers a special reading address
	// Make sure address is not already in SpecRead
	_, ok := m.SpecRead[addr]
	if ok {
		// TODO add real special function
		fmt.Println("MEM: DUMMY: SpecialRead from", addr.HexString())
	}

	for _, c := range m.Chunks {

		if c.contains(addr) {
			b = c.fetch(addr)
			found = true
			break
		}
	}

	if !found {
		return 0, fmt.Errorf("mem.Fetch: Address %s not in RAM or ROM", addr.HexString())
	}

	return b, nil
}

// FetchMore takes a 65816 address and the number of bytes to get -- 1, 2 or 3
// -- to be fetched and returned as an integer, retrieving those bytes as
// little endian. Also returns an error to show if all fetches were to legal
// addresses. Assumes that the address itself was vetted.
func (m Memory) FetchMore(addr common.Addr24, num uint) (uint, error) {
	const maxint uint = 3
	var sum uint

	if num > maxint {
		return 0, fmt.Errorf("mem.FetchMore: Can fetch up to 3 bytes, not %d", num)
	}

	for i := common.Addr24(0); i <= common.Addr24(num-1); i++ {
		b, err := m.Fetch(addr + i)
		if err != nil {
			return 0, fmt.Errorf("mem.FetchMore: Can't get byte from %s: %v", i.HexString(), err)
		}

		// Shift eight bits to the left for every byte we go further to
		// the right
		sum += (uint(b) << (8 * i))
	}
	return sum, nil
}

// List returns a list of all chunks in memory, as a string
func (m Memory) List() string {
	const template string = "%s %s %s (%d bytes)\n"
	var r string

	for _, c := range m.Chunks {
		r += fmt.Sprintf(template,
			c.Start.HexString(), c.End.HexString(),
			strings.ToUpper(c.Type), c.size())
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
		b, err := m.Fetch(i)
		if err != nil {
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
		sum += c.size()
	}

	return sum
}

// NOTE: Store and Burn are different only in the test for the memory type. We
// could combine them to one routine, but that adds complexity and this is
// easier to maintain, DRY be damned.

// Store takes a 24-bit address and a byte. If the address is legal for a write,
// that is, actually part of a RAM chunk, the byte is stored at the address and
// a nil error is returned to signal success. If the address is outside the
// defined range or in ROM, there is no write, and an error message is returned.
// This is the main store routine for the emulator. See mem.Burn for a function
// that ignores ROM/RAM differences.
func (m Memory) Store(addr common.Addr24, b byte) error {

	// First, see if this triggers a special write address
	// Make sure address is not already in SpecRead
	_, ok := m.SpecWrite[addr]
	if ok {
		// TODO add real special function
		fmt.Println("MEM: DUMMY: SpecialWrite to", addr.HexString())
		return nil
	}

	// Walk through the memory chunks looking for the address. If we don't
	// find it, return an error code
	var found bool

	for _, c := range m.Chunks {

		// Assumes we're short-circuiting thought it seems that this is
		// not explicitly in the Go specs
		if c.Type == "ram" && c.contains(addr) {
			c.store(addr, b)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Store: address %s not found in RAM", addr.HexString())
	}

	return nil
}

// Burn takes a 24-bit address and a byte. If the address is part of a chunk,
// RAM or ROM, the byte is stored at the address and a true flag is returned to
// signal success. If the address is outside the defined range, an error is
// returned. Burn is used to write to memory during intialization and the load
// commands, as well as storing bytes from the CLI. The main routine for
// assembler instructions is Store.
func (m Memory) Burn(addr common.Addr24, b byte) error {

	// First, see if this triggers a special write address
	// Make sure address is not already in SpecRead
	_, ok := m.SpecWrite[addr]
	if ok {
		// TODO add real special function
		fmt.Println("MEM: DUMMY: SpecialWrite to", addr.HexString())
		return nil
	}

	var found bool

	for _, c := range m.Chunks {

		if c.contains(addr) {
			c.store(addr, b)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Burn: address %s not found in RAM or ROM", addr.HexString())
	}

	return nil
}

// StoreMore takes an address, a number and the number of bytes to store little
// endian at that address. If any part of the address is not a part of legal
// memory, return a false flag, otherwise true. At most, numbers up to 24 bit
// length (three bytes) are stored. Anything above that is silently discarded.
// If the number of bytes to store is anything but 1, 2, or 3, we return an
// error flag
func (m Memory) StoreMore(addr common.Addr24, num uint, len uint) error {

	if len < 1 || len > 3 {
		return fmt.Errorf("StoreMore: can only store 1, 2, or 3 bytes, not %d", len)
	}

	lsb := byte(num & 0xff)
	msb := byte((num & 0xff00) >> 8)
	bank := byte((num & 0xff0000) >> 16)

	ba := [...]byte{lsb, msb, bank}

	for i := 0; i < 3; i++ {
		err := m.Store(addr+common.Addr24(i), ba[i])
		if err != nil {
			return fmt.Errorf("StoreMore: %v", err)
		}
	}

	return nil
}

// NOTE: StoreBlock and BurnBlock are different only in the call to Store or
// Burn for the individual bytes. We could combine them to one routine and pass
// the function (technically a method) as a parameter, but that adds complexity
// and this is easier to maintain, DRY be damned.

// StoreBlock takes a 65816 address and a slice of bytes and stores those bytes
// starting at that address. If all addresses were legal, it returns a true
// flag, otherwise a false. StoreBlock will not write to ROM, use BurnBlock for
// that.
func (m Memory) StoreBlock(addr common.Addr24, bs []byte) error {

	for i := addr; i < addr+common.Addr24(len(bs)); i++ {
		err := m.Store(i, bs[i-addr])

		if err != nil {
			return fmt.Errorf("StoreBlock: %v", err)
		}
	}

	return nil
}

// BurnBlock takes a 65816 address and a slice of bytes and stores those bytes
// starting at that address. If all addresses were legal, it returns a true
// flag, otherwise a false
func (m Memory) BurnBlock(addr common.Addr24, bs []byte) error {

	for i := addr; i < addr+common.Addr24(len(bs)); i++ {

		err := m.Burn(i, bs[i-addr])
		if err != nil {
			return fmt.Errorf("BurnBlock: %v", err)
		}
	}

	return nil
}
