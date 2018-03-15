// go65816 Memory System
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 09. Mar 2018
// This version: 15. Mar 2018

package mem

import (
	"fmt"
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

// containsAddr takes a memory address and checks if it is in this chunk,
// returning a bool. Assumes that the address has been confirmed to be a valid
// 65816 address as a uint
func (c Chunk) containsAddr(addr uint) bool {
	return c.start <= addr && addr <= c.end
}

// fetch gets one byte of memory from the data of a chunk and returns it.
// Assumes we have already made sure that the address is in this chunk
func (c Chunk) fetch(addr uint) byte {
	index := addr - c.start
	return c.data[index]
}

// hexdump prints the chunk's memory contents in a nice hex table
// TODO add ASCII dump
func (c Chunk) hexdump() {

	fmt.Printf("%06X: ", c.start)

	count := 0

	for _, b := range c.data {
		fmt.Printf(" %02x", b)

		count += 1

		if count%8 == 0 {
			fmt.Print(" ")
		}

		if count%16 == 0 {
			fmt.Print("\n")
			fmt.Printf("%06X: ", c.start+uint(count))
		}
	}

	fmt.Println()
}

// Size returns the, uh, size of a chunk in bytes
func (c Chunk) size() uint {
	return c.end - c.start
}

// Store takes a byte and an address and stores the byte at the address in the
// chunk. Assumes that we already checked that the address is in fact in this
// chunk
func (c Chunk) store(b byte, addr uint) {
	c.data[addr] = b
}
