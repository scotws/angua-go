// Angua - An Emulator for the 65816 CPU in Native Mode
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 26. Sep 2017
// This version: 24. May 2019

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

// This is the initial Angua program. It reads the configuration file and starts
// the command line interface in a different go routine to be accessed through a
// separate terminal. This terminal is where the actual emulator runs, allowing
// us to use trickery to input single characters and such.

package main

import (
	"fmt"

	"angua/cli"
	"angua/info"
)

func main() {

	// Generate Dictionaries for the information system in the
	// background
	go info.GenerateDicts()

	// Create three channels for the Command Line Interface (CLI). 'cmd' are
	// the commands sent from the CLI running in a separate go routine to the
	// acutal emulator in this routine; 'resp' is the response string
	// sent to the CLI; done signals to the emulator that everything is over
	// and it is time to quit.
	cmd := make(chan int, 2)
	resp := make(chan string, 2)
	done := make(chan struct{})

	// Start the Command Line Interface (CLI) as a separate go routine.
	go cli.Start(cmd, resp, done)

command_loop:
	for {
		select {
		// The done channel is used to end the imput and quit. We need
		// the label so break takes us not only out of the select
		// construct but also the loop
		case <-done:
			break command_loop

		case s := <-cmd:
			fmt.Println("Got command:", s)
		}
	}
}
