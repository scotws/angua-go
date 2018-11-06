// Angua - A 65816 MPU emulator in Go
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 26. Sep 2017
// This version: 06. Nov 2018

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

package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"angua/mem"
)

const (
	maxAddr = 1<<24 - 1
)

var (
	confs  []string
	memory mem.Memory

	specials  = make(map[uint]string)
	beVerbose = flag.Bool("v", false, "Verbose, print more output")
)

// -----------------------------------------------------------------
// TOP LEVEL HELPER FUNCTIONS

// verbose takes a string and prints it on the standard output through logger if
// the user awants us to be verbose
func verbose(s string) {
	if *beVerbose {
		log.Print(s)
	}
}

// convNum Converts a number string -- hex, binary, or decimal -- to an uint.
// We accept ':' and '.' as delimiters, use $ or 0x for hex numbers, % for
// binary numbers, and nothing for decimal numbers.
func convNum(s string) uint {

	ss := stripDelimiters(s)

	switch {

	case strings.HasPrefix(ss, "$"):
		n, err := strconv.ParseInt(ss[1:], 16, 0)
		if err != nil {
			log.Fatal("config.sys: Can't convert ", ss, " as hex number")
		}
		return uint(n)

	case strings.HasPrefix(ss, "0x"):
		n, err := strconv.ParseInt(ss[2:], 16, 0)
		if err != nil {
			log.Fatal("config.sys: Can't convert ", ss, " as hex number")
		}
		return uint(n)

	case strings.HasPrefix(ss, "%"):
		n, err := strconv.ParseInt(ss[1:], 2, 0)
		if err != nil {
			log.Fatal("config.sys: Can't convert ", ss, " as binary number")
		}
		return uint(n)

	default:
		n, err := strconv.ParseInt(ss, 10, 0)
		if err != nil {
			log.Fatal("config.sys: Can't convert ", ss, " as decimal number")
		}
		return uint(n)
	}
}

// fmtAddr takes a 65816 address as a uint and returns a hex number string with
// a ':' between the bank byte and the rest of the address. Hex digits are
// capitalized. Assumes we are sure that the address is valid
func fmtAddr(addr uint) string {
	s1 := fmt.Sprintf("%06X", addr)
	s2 := s1[0:2] + ":" + s1[2:len(s1)]
	return s2
}

// isValidAddr takes an uint and makes sure that as an address, it is not larger
// than can be stated with 24 bits We don't need to test for negative numbers
// because we force uint
// TODO Move this to CPU and make special for 16 and 24 bits
func isValidAddr(a uint) bool {
	return a <= maxAddr
}

// stripDelimiters removes '.' and ':' which users can use as number delimiters.
// Also removes spaces
func stripDelimiters(s string) string {
	s1 := strings.Replace(s, ":", "", -1)
	s2 := strings.Replace(s1, ".", "", -1)
	return strings.TrimSpace(s2)
}

// -----------------------------------------------------------------
// MAIN ROUTINE

func main() {

	flag.Parse()

	// TODO --- Load configuration ---

	fmt.Println(" ---- (TESTING) ----")

	bs := []byte{0xFF, 0xee, 0xdd}

	_ = memory.Write(0xa000, bs)

	n, ok := memory.FetchMore(0xa000, 3)

	if ok {
		fmt.Println("All is well")
	} else {
		fmt.Println("This sucks")
	}

	fmt.Printf("Number: %d (%06X)\n", n, n)

	memory.Hexdump(0xa000, 0xa020)

}
