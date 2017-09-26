// py65816 A 65816 MPU emulator in MPU
// Scot W. Stevenson scot.stevenson@gmail.com
// First version: 26. Sep 2017
// Second version: 26. Sep 2017

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

	"go65816/config"
)

const configFile = "config.sys"

type memBlock struct {
	class  string
	start  int
	end    int
	source string
}

var (
	confs []string
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

	// TODO Testing print lines

	for _, l := range confs {

		if config.IsComment(l) {
			continue
		}

		if config.IsEmpty(l) {
			continue
		}

		fmt.Println(l)
	}
}
