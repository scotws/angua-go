// Angua - A 65816 MPU emulator in Go
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 26. Sep 2017
// This version: 11. Nov 2018

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
	"runtime"
	"strings"

	"angua/common"
	"angua/cpu16"
	"angua/cpu8"
	// "angua/mem"
	"angua/switcher"

	"gopkg.in/abiosoft/ishell.v2"
)

const (
	// TODO make this pretty
	shellBanner = "Angua 65816 Emulator\n(c) 2018 Scot W. Stevenson"
)

var (
	haveMachine bool = false

	// Flags passed.
	// TODO Add "-c" to load config file
	// TODO Add "-d" to print debug information
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

	// memory := &mem.Memory{}
	cpuEmu := &cpu8.Cpu8{}
	cpuNat := &cpu16.Cpu16{}

	// We communicate with the system through the command channel, which is
	// buffered because lots of other stuff might be going on. Both CPUs see
	// the same channel, but because one of them is inactive, that's not a
	// problem
	cmd := make(chan int, 2)

	// Start interactive shell. Note that by default, this provides the
	// directives "exit", "help", and "clear"
	shell := ishell.New()
	shell.Println(shellBanner)
	shell.SetPrompt("> ")

	// We create a history file
	// TODO point this out in the documentation
	shell.SetHomeHistoryPath(".angua_shell_history")

	shell.AddCmd(&ishell.Cmd{
		Name: "abort",
		Help: "Trigger the ABORT vector",
		Func: func(c *ishell.Context) {
			c.Println("Sending ABORT signal to machine ...")
			cmd <- common.ABORT
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "beep",
		Help:     "Print a beeping noise",
		LongHelp: "No, seriously, it only produces a beeping sound.",
		Func: func(c *ishell.Context) {
			c.Println("\a")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "boot",
		Help: "Boot the machine. Same effect as turning on the power",
		Func: func(c *ishell.Context) {
			c.Println("Sending BOOT signal to machine ...")
			cmd <- common.BOOT
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "destroy",
		Help: "destroy the machine",
		Func: func(c *ishell.Context) {
			if !haveMachine {
				c.Println("ERROR: No machine present")
			} else {
				c.Println("CLI: DUMMY: destroy the machine")
				haveMachine = false
				shell.Process("beep")
			}

			// TODO Call HALT
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "disasm",
		Help: "Disassemble a range of memory",
		Func: func(c *ishell.Context) {
			c.Println("CLI: DUMMY: disassemble memory")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "dump",
		Help: "Print hex dump of range",
		Func: func(c *ishell.Context) {
			c.Println("CLI: DUMMY: dump")
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
		Name:     "fill",
		Help:     "Fill a bock of memory with a byte",
		LongHelp: "Format: 'fill <ADDRESS RANGE> with <BYTE>'",
		Func: func(c *ishell.Context) {
			c.Println("CLI: DUMMY: fill")
		},
	})

	// TODO Set a flag to signal that system is halted so commands
	// like status make more sense. Also, we need some way of remembering
	// which CPU we last left off with so we can start in the right mode
	// Note that RESUME is currently broken
	shell.AddCmd(&ishell.Cmd{
		Name: "halt",
		Help: "Halt the machine",
		Func: func(c *ishell.Context) {
			c.Println("Requesting HALT from machine ...")
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

			c.Println("Initializing machine ...")
			haveMachine = true

			// Start the Switcher Daemon which in turn launches the
			// CPUs. It will handle any requests to switch from
			// native to emulated mode and back again without our
			// intervention
			go switcher.Daemon(cpuEmu, cpuNat, cmd)

			if *beVerbose {
				c.Println("Switcher daemon launched")
			}

			c.Println("*** System initialized, start with 'run' or 'boot' ***")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "irq",
		Help: "trigger an Interrupt Request (IRQ)",
		Func: func(c *ishell.Context) {
			c.Println("Triggering maskable interrupt request (IRQ) ...")
			cmd <- common.IRQ
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "load",
		Help:     "load contents of a file to memory",
		LongHelp: "Format: 'load <FILENAME> to <ADDRESS RANGE>'",
		Func: func(c *ishell.Context) {
			c.Println("CLI: DUMMY: load")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "man",
		Help:     "Print information on 65816 instructions",
		LongHelp: "Format 'man [ <OPCODE> | <MNEMONIC> ]'",
		Func: func(c *ishell.Context) {
			c.Println("CLI: DUMMY: man")
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
			c.Println("Triggering non-maskable interrupt (NMI) ...")
			cmd <- common.NMI
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "reading",
		Help: "set a special address",
		Func: func(c *ishell.Context) {
			c.Println("CLI: DUMMY reading")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "reset",
		Help: "trigger a RESET signal",
		Func: func(c *ishell.Context) {
			c.Println("Triggering RESET signal ...")
			cmd <- common.RESET
		},
	})

	// TODO this doesn't work at all, see RUN
	shell.AddCmd(&ishell.Cmd{
		Name: "resume",
		Help: "resume after a halt",
		Func: func(c *ishell.Context) {
			c.Println("Resuming at current PC location ...")
			cmd <- common.RESUME
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "run",
		Help:     "Run machine created with 'init'",
		LongHelp: "System starts as a boot in Emulated Mode.",
		Func: func(c *ishell.Context) {
			fmt.Println("(CLI: DUMMY run: Triggering cpuEmu)")
			switcher.Run(cpuEmu)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "save",
		Help: "save an address range to file",
		Func: func(c *ishell.Context) {
			c.Println("(CLI: DUMMY run)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "set",
		Help: "set various parameters",
		Func: func(c *ishell.Context) {
			c.Println("(CLI: DUMMY set)")
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
				case "system":
					c.Println("(Use 'status system' to show host information)")
				case "vectors":
					c.Println("(DUMMY show vectors)")
				default:
					c.Println("ERROR: Option", subcmd, "unknown")
				}
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "status",
		Help:     "display status of the machine",
		LongHelp: "Options: 'system' or ''",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				cmd <- common.STATUS
			} else {
				subcmd := c.Args[0]

				switch subcmd {
				case "system":
					fmt.Println("Host architecture:", runtime.GOARCH)
					fmt.Println("Host operating system:", runtime.GOOS)
					fmt.Println("Host system CPU cores available:", runtime.NumCPU())
					fmt.Println("Host system goroutines running:", runtime.NumGoroutine())
					fmt.Println("Host system Go version:", runtime.Version())
				default:
					c.Println("ERROR: Option", subcmd, "unknown")
				}
			}

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
