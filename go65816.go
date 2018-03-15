// py65816 A 65816 MPU emulator in Go
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 26. Sep 2017
// This version: 15. Mar 2018

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
	// "bufio"
	"flag"
	//	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	// "go65816/config"
	"go65816/mem"
)

const (
	configFile = "config.cfg"
	maxAddr    = 1<<24 - 1
)

var (
	confs  []string
	memory mem.Memory

	beVerbose = flag.Bool("v", false, "Verbose, print more output")

	// Default values for special locations
	// These can be overridden
	getc       = 0x00f000
	getc_block = 0x00f001
	putc       = 0x00f002
)

// -----------------------------------------------------------------

// verbose takes a string and prints it on the standard output through logger if
// the user awants us to be verbose
func verbose(s string) {
	if *beVerbose {
		log.Print(s)
	}
}

// isValidAddr takes an uint and makes sure that as an address, it is not larger
// than can be stated with 24 bits We don't need to test for negative numbers
// because we force uint
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

// convNum Convert a legal number string to an uint. Note we accept ':' and '.'
// as delimiters, use $ or 0x for hex numbers, % for binary numbers, and nothing
// for decimal numbers.
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

// -----------------------------------------------------------------

func main() {

	flag.Parse()

	// --- Load configuration ---

	cf, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer cf.Close()
	/*
		source := bufio.NewScanner(cf)

		for source.Scan() {
				confs = append(confs, source.Text())
			}

			for _, l := range confs {

				if config.IsComment(l) {
					continue
				}

				if config.IsEmpty(l) {
					continue
				}

				ws := strings.Fields(l)

				if config.IsChunkDef(ws[0]) {
					memory = append(memory, makeChunk(ws))
				} else {
					// TODO Test
					fmt.Println(ws)

				}
			}

			fmt.Println(memory) */
}
