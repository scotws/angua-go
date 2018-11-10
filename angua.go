// Angua - A 65816 MPU emulator in Go
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 26. Sep 2017
// This version: 10. Nov 2018

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
	"log"
	"strings"

	"angua/common"
	"angua/cpu16"
	"angua/cpu8"
	"angua/mem"
	"angua/xo" // TODO Remove this

	"gopkg.in/abiosoft/ishell.v2"
)

const (
	// TODO make this pretty
	shellBanner = "Angua 65816 Emulator\n(c) 2018 Scot W. Stevenson"
)

var (
	memory *mem.Memory
	cpuEmu *cpu8.Cpu8
	cpuNat *cpu16.Cpu16

	haveMachine bool = false

	// Flags passed.
	// TODO Add "-c" to load config file
	beVerbose   = flag.Bool("v", false, "Verbose, print more output")
	inBatchMode = flag.Bool("b", false, "Start in batch mode")
)

// verbose takes a string and prints it on the standard output through logger if
// the user awants us to be verbose
func verbose(s string) {
	if *beVerbose {
		log.Print(s)
	}
}

// -----------------------------------------------------------------
// MAIN ROUTINE

func main() {

	flag.Parse()

	// We communicate with the XO through the command channel once its main
	// routine (xo.MakeItSo) is up and running
	// TODO see if we really need to have this buffered
	cmd := make(chan int, 2)

	// The enable channels single the various processors that it is time for
	// them to run
	// enable8 := make(chan struct{})
	// enable16 := make(chan struct{})

	// Start interactive shell. Note that by default, this provides the
	// directives "exit", "help", and "clear"
	shell := ishell.New()
	shell.Println(shellBanner)
	shell.SetPrompt("> ")

	// We create a history file
	shell.SetHomeHistoryPath(".angua_shell_history")

	shell.AddCmd(&ishell.Cmd{
		Name: "abort",
		Help: "Trigger the abort vector",
		Func: func(c *ishell.Context) {
			c.Println("CLI: DUMMY: trigger abort vector")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "beep",
		Help: "Print a beeping noise",
		Func: func(c *ishell.Context) {
			c.Println("\a")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "boot",
		Help: "Boot the machine. Same effect as turning on the power",
		Func: func(c *ishell.Context) {
			c.Println("CLI: DUMMY: boot the machine")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "destroy",
		Help: "destroy the machine",
		Func: func(c *ishell.Context) {
			if !haveMachine {
				c.Println("ERROR: No machine present")
			} else {
				c.Println("(DUMMY destroy the machine)")
				haveMachine = false
				shell.Process("beep")
			}

			// TODO Call HALT and close channel to XO
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "disasm",
		Help: "Disassemble a range of memory",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY disassemble memory)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "dump",
		Help: "Print hex dump of range",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY dump)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "echo",
		Help: "Print following text to end of line",
		Func: func(c *ishell.Context) {
			c.Println(strings.Join(c.Args, " "))
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "halt",
		Help: "Halt the machine",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY halt the machine)")

			// Send XO the halt signal
			cmd <- common.HALT
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "init",
		Help:     "initialize a new machine",
		LongHelp: "Options: '<CONFIG FILE>', 'default', ''",
		Func: func(c *ishell.Context) {

			if haveMachine {
				c.Println("ERROR: Already have machine")
				return
			}

			// Three variants: Without a parameter or with the words
			// "defalt", load the default.cfg file from configs;
			// with a filename, load the file cfom configs

			// TODO set up memory by reading cfg file
			// TODO Send pointers to cpuEmu and cpuNat to the xo
			// TODO Start the xo as a go routine

			c.Println("CLI: DUMMY: init")
			haveMachine = true

			// TODO remove the xo package

			xo.Init(cpuEmu, cpuNat)

			go xo.MakeItSo(cpuEmu, cpuNat, cmd)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "irq",
		Help: "trigger an Interrupt Request (IRQ)",
		Func: func(c *ishell.Context) {
			c.Println("CLI: DUMMY: IRQ")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "load",
		Help: "load contents of a file to memory",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY load)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "memory",
		Help: "define a memory chunk",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY Memory)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "mode",
		Help:     "set CPU mode",
		LongHelp: "Options: 'native', 'emulated'",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY mode )")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "nmi",
		Help: "trigger a Non Maskable Interrupt (NMI)",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY NMI)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "reading",
		Help: "set a special address",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY reading)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "reset",
		Help: "trigger a RESET signal",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY reset)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "resume",
		Help: "resume after a halt",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY resume)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "run",
		Help: "run from a given address",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY run)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "save",
		Help: "save an address range to file",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY run)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "set",
		Help: "set various parameters",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY set)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: "display information on various parts of the system",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("ERROR: Need an argument")
			} else {
				subcmd := c.Args[0]

				switch subcmd {
				case "config":
					c.Println("(DUMMY show config)")
				case "memory":
					c.Println("(DUMMY show memory)")
				case "specials":
					c.Println("(DUMMY show specials)")
				case "vectors":
					c.Println("(DUMMY show vectors)")
				default:
					c.Println("ERROR: Option", subcmd, "unknown")
				}
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "status",
		Help: "display status of the machine",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY status)")

		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "store",
		Help: "store byte at a given address",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY store)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "writing",
		Help: "define a function to be triggered when address written to",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY writing)")
		},
	})

	// TODO check for batch mode
	shell.Run()
	shell.Close()

}
