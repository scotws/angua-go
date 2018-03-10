// go65816 Memory Routines
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 09. Mar 2018
// This version: 09. Mar 2018

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package mem

import "fmt"

const (
	maxAddr = 1<<24 - 1
)

type Chunk struct {
	start     uint
	end       uint
	writeable bool   // true for RAM, false for ROM
	file      string // file path
	name      string // For user only
	data      []byte
}

// Return size of a chunk
func (c Chunk) size() uint {
	return c.end - c.start
}

// Clear the memory region. This is also used to initilize the data
func (c Chunk) erase() {
	c.data = make([]byte, c.end-c.start)
}

// Dump the chunk's memory in hex
func (c Chunk) hexdump() {
	for _, b := range c.data {
		fmt.Print(b)
	}
}

var (
	// The whole memory of the machine is a list of chunks
	memory []Chunk

	// Default values for special locations
	// These can be overridden
	getc       = 0x00f000
	getc_block = 0x00f001
	putc       = 0x00f002
)

// Make sure address is not larger than can be stated with 24 bits
func isValidAddr(a uint) bool {
	return 0 <= a && a <= maxAddr
}
