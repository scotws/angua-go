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

	// "github.com/fatih/color"
	"gopkg.in/abiosoft/ishell.v2"
)

const (
	maxAddr     = 1<<24 - 1
	shellBanner = "Angua 65816 Emulator\n(c) 2018 Scot W. Stevenson"
)

var (
	memory mem.Memory

	haveMachine bool = false

	specials = make(map[uint]string)

	// Flags passed. Add "-c" to load config file
	beVerbose   = flag.Bool("v", false, "Verbose, print more output")
	inBatchMode = flag.Bool("b", false, "Start in batch mode")
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
			c.Println("(DUMMY trigger abort vector)")
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
			c.Println("(DUMMY boot the machine)")
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
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "init",
		Help: "initialize a new machine",
		Func: func(c *ishell.Context) {

			if haveMachine {
				c.Println("ERROR: Already have machine")
				return
			}

			c.Println("(DUMMY init)")
			haveMachine = true
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "irq",
		Help: "trigger an Interrupt Request (IRQ)",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY IRQ)")
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
