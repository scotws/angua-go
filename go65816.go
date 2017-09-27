// py65816 A 65816 MPU emulator in MPU
// Scot W. Stevenson scot.stevenson@gmail.com
// First version: 26. Sep 2017
// Second version: 27. Sep 2017

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

const configFile = "config.sys"

type memBlock struct {
	class  string // "ram" or "rom"
	start  int
	end    int
	size   int
	data   *[]byte
	source string // ROM file path
}

var (
	confs  []string
	memory []memBlock

	// Default values for special locations
	// These can be overridden by config file
	getc       = 0xf000
	getc_block = 0xf001
	putc       = 0xf002
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

		if config.IsMemBlockDef(ws[0]) {
			memory = append(memory, makeMemBlock(ws))
		} else {
			// TODO Test
			fmt.Println(ws)

		}
	}

	fmt.Println(memory)
}

func makeMemBlock(ws []string) memBlock {

	s := convNum(ws[1])
	e := convNum(ws[2])
	sz := e - s + 1
	d := make([]byte, sz)
	prt := &d

	// ROM memory blocks get link to their content
	a := ""
	if len(ws) == 4 {
		a = ws[3]
	}

	return memBlock{class: ws[0], start: s, end: e, size: sz, data: prt, source: a}
}

// Convert a legal number string to an int. Note we accept ':' and '.' as delimiters,
// use $ for hex numbers, % for binary numbers, and nothing for decimal numbers.
// TODO code test
func convNum(s string) int {

	s1 := strings.Replace(s, ":", "", -1)
	s2 := strings.Replace(s1, ".", "", -1)
	s3 := strings.TrimSpace(s2)

	d := s3[0]

	switch d {

	case '$':
		n, err := strconv.ParseInt(s3[1:], 16, 0)
		if err != nil {
			log.Fatal("config.sys: Can't convert ", s3, " as hex number")
		}
		return int(n)

	case '%':
		n, err := strconv.ParseInt(s3[1:], 2, 0)
		if err != nil {
			log.Fatal("config.sys: Can't convert ", s3, " as binary number")
		}
		return int(n)

	default:
		n, err := strconv.ParseInt(s3, 10, 0)
		if err != nil {
			log.Fatal("config.sys: Can't convert ", s3, " as decimal number")
		}
		return int(n)
	}
}
