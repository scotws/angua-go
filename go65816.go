// py65816 A 65816 MPU emulator in Go
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 26. Sep 2017
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

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"go65816/config"
)

const (
	configFile = "config.sys"
	maxAddr    = 1<<24 - 1
)

type chunk struct {
	class  string // "ram" or "rom"
	start  uint
	end    uint
	size   uint
	data   *[]byte
	source string // ROM file path
}

var (
	confs  []string
	memory []chunk

	// Default values for special locations
	// These can be overridden
	getc       = 0x00f000
	getc_block = 0x00f001
	putc       = 0x00f002
)

func main() {

	// *** CONFIGURATION FILE ***

	cf, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer cf.Close()

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

	fmt.Println(memory)
}

func makeChunk(ws []string) chunk {

	s := convNum(ws[1])

	if !isValidAddr(s) {
		log.Fatal("Can't use ", s, " as start address")
	}

	e := convNum(ws[2])

	if !isValidAddr(e) {
		log.Fatal("Can't use ", e, " as end address")
	}

	sz := e - s + 1
	d := make([]byte, sz)
	prt := &d

	// ROM memory blocks get link to their content
	a := ""
	if len(ws) == 4 {
		a = ws[3]
	}

	return chunk{class: ws[0], start: s, end: e, size: sz, data: prt, source: a}
}

// Make sure address is not larger than can be stated with 24 bits
// TODO code test
func isValidAddr(a uint) bool {
	return a <= maxAddr
}

// Remove '.' and ':' which users can use as number delimiters. Also removed
// spaces
// TODO code test
func stripDelim(s string) string {
	s1 := strings.Replace(s, ":", "", -1)
	s2 := strings.Replace(s1, ".", "", -1)
	return strings.TrimSpace(s2)
}

// Convert a legal number string to an int. Note we accept ':' and '.' as delimiters,
// use $ for hex numbers, % for binary numbers, and nothing for decimal numbers.
// TODO code test
func convNum(s string) uint {

	ss := stripDelim(s)

	d := ss[0]

	switch d {

	case '$':
		n, err := strconv.ParseInt(ss[1:], 16, 0)
		if err != nil {
			log.Fatal("config.sys: Can't convert ", ss, " as hex number")
		}
		return uint(n)

	case '%':
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
