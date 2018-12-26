// Angua - A partial 65816 MPU emulator in Go
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 26. Sep 2017
// This version: 26. Dec 2018

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
	"sync"

	"angua/common"
	"angua/cpu"
	"angua/info"
	"angua/mem"

	"gopkg.in/abiosoft/ishell.v2"
)

const (
	shellBanner string = `The Angua partial 65816 emulator
Version ALPHA 0.1  26. Dec 2018
Copyright (c) 2018 Scot W. Stevenson
Angua comes with absolutely NO WARRANTY
Type 'help' for more information`
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

// parseAddressRange takes a list of strings that either has a format in form of
// "<BANK>:<ADDR16> to <BANK>:<ADDR16>" or "bank <BYTE>" and returns two
// addresses in the common.Addr24 format and a bool for success for failure.
func parseAddressRange(ws []string) (addr1, addr2 common.Addr24, ok bool) {

	ok = true

	// If the first word is "bank", then we are getting a full bank
	if ws[0] == "bank" {

		// Second word must be the bank byte. We brutally cut off
		// everything but the lowest byte
		bankNum := common.ConvNum(ws[1]) // Returns uint
		bankByte := common.Addr24(bankNum).Lsb()
		bankAddr := common.Addr24(bankByte) * 0x10000
		addr1 = bankAddr
		addr2 = bankAddr + 0xFFFF

	} else {
		// We at least need two addresses and the memory type, so that's
		// three words length. We could parse more carefully, but not at
		// the moment
		if len(ws) < 3 {
			addr1 = 0
			addr2 = 0
			ok = false
			return addr1, addr2, ok
		}

		addr1 = common.Addr24(common.ConvNum(ws[0]))

		// We allow people to slide on the "to" though we don't
		// advertise the fact. Later, once we have the error handling of
		// ConvNum working, check ws[1] and if there is an error, skip
		// to ws[2].

		if ws[1] == "to" {
			addr2 = common.Addr24(common.ConvNum(ws[2]))
		} else {
			addr2 = common.Addr24(common.ConvNum(ws[1]))
		}
	}
	return addr1, addr2, ok
}

// -----------------------------------------------------------------

func main() {

	// Generate Dictionaries for 'info' system in the background
	go info.GenerateDicts()

	flag.Parse()

	memory := &mem.Memory{}

	// We communicate with the system through the command channel, which is
	// buffered because other stuff might be going on.
	cmd := make(chan int, 2)

	// Start interactive shell. Note that by default, this provides the
	// directives "exit", "help", and "clear".
	shell := ishell.New()
	shell.Println(shellBanner)
	shell.SetPrompt("> ")

	// We create a history file
	// TODO point this out in the documentation
	shell.SetHomeHistoryPath(".angua_shell_history")

	// Individual commands. Normal help is lower case with no punctuation.

	shell.AddCmd(&ishell.Cmd{
		Name:     "abort",
		Help:     "trigger the ABORT vector",
		LongHelp: longHelpAbort,
		Func: func(c *ishell.Context) {
			c.Println("Sending ABORT signal to machine ...")
			cmd <- common.ABORT
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "beep",
		Help:     "make a beeping noise",
		LongHelp: longHelpBeep,
		Func: func(c *ishell.Context) {
			c.Println("\a")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "boot",
		Help:     "boot the machine (cold restart)",
		LongHelp: longHelpBoot,
		Func: func(c *ishell.Context) {
			c.Println("Sending BOOT signal to machine ...")
			cmd <- common.BOOT
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "destroy",
		Help:     "destroy the machine",
		LongHelp: longHelpDestroy,
		Func: func(c *ishell.Context) {
			if !haveMachine {
				c.Println("ERROR: No machine present")
				return
			}

			c.Println("CLI: DUMMY: destroy the machine")
			haveMachine = false
			shell.Process("beep")

			// TODO Call HALT
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "disasm",
		Help:     "disassemble a range of memory",
		LongHelp: longHelpDisasm,
		Func: func(c *ishell.Context) {
			c.Println("CLI: DUMMY: disassemble memory")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "dump",
		Help:     "print hex dump of a memory range",
		LongHelp: longHelpDump,
		Func: func(c *ishell.Context) {

			if !haveMachine {
				c.Println("ERROR: No machine present")
				return
			}

			// The arg string is passed without "dump"
			a1, a2, ok := parseAddressRange(c.Args)
			if !ok {
				c.Println("ERROR parsing address range")
				return
			}
			memory.Hexdump(a1, a2)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "echo",
		Help:     "print a string of text",
		LongHelp: longHelpEcho,
		Func: func(c *ishell.Context) {
			c.Println(strings.Join(c.Args, " "))
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "fill",
		Help:     "fill a block of memory with a byte",
		LongHelp: longHelpFill,
		Func: func(c *ishell.Context) {
			c.Println("CLI: DUMMY: fill")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "halt",
		Help:     "halt the machine (freeze)",
		LongHelp: longHelpHalt,
		Func: func(c *ishell.Context) {
			c.Println("Telling machine to halt ...")
			cmd <- common.HALT
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "init",
		Help:     "initialize a new machine",
		LongHelp: longHelpInit,
		Func: func(c *ishell.Context) {

			if haveMachine {
				c.Println("ERROR: Already have machine")
				return
			}

			// Three variants: Without a parameter or with the words
			// "default", load the default.cfg file from configs;
			// with a filename, load the file cfom configs

			// TODO set up memory by reading cfg file

			c.Println("Initializing machine ...")
			haveMachine = true

			c.Println("*** System initialized, start with 'run' or 'boot' ***")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "irq",
		Help:     "trigger an interrupt request (IRQ)",
		LongHelp: longHelpIRQ,
		Func: func(c *ishell.Context) {
			c.Println("Triggering maskable interrupt request (IRQ) ...")
			cmd <- common.IRQ
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "load",
		Help:     "load binary file to memory",
		LongHelp: longHelpLoad,
		Func: func(c *ishell.Context) {
			c.Println("CLI: DUMMY: load")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "info",
		Help:     "print information on 65816 instructions",
		LongHelp: longHelpInfo,
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("ERROR: Need opcode or SAN mnemonic")
			} else {
				subcmd := c.Args[0]

				// First see if this is a mnemonic
				opc, ok := info.SANDict[subcmd]
				if ok {
					info.PrintOpcodeInfo(opc.Opcode)
					return
				}

				// TODO Okay, it isn't. Then see if it is an
				// opcode

				// TODO Okay, still not good. Then see about
				// another command such as 'all' or 'list'

				c.Println("ERROR: Opcode or mnemonic", subcmd, "unknown")
			}

		},
	})

	shell.AddCmd(&ishell.Cmd{
		// Memory commands take the form of
		// 	 memory <ADDRESS RANGE> ["is"] ("rom" | "ram")
		// TODO Add labels
		Name:     "memory",
		Help:     "define a memory chunk",
		LongHelp: longHelpMemory,
		Func: func(c *ishell.Context) {

			// We can break this up from the end by making sure that the
			// last word is either "ram" or "rom", and feeding the beginning
			// to the address range finder
			memType := c.Args[len(c.Args)-1]

			if memType != "ram" && memType != "rom" {
				c.Println("ERROR: Last word must be memory type ('ram' or 'rom')")
				return
			}

			// The arg string is passed without "memory"
			a1, a2, ok := parseAddressRange(c.Args)
			if !ok {
				c.Println("ERROR parsing address range")
				return
			}

			newChunk := mem.Chunk{a1, a2, memType, "<NONE>", sync.Mutex{}, []byte{}}
			memory.Chunks = append(memory.Chunks, newChunk)

			// We can at least allow stuff like hexdumps of memory
			haveMachine = true
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "nmi",
		Help:     "trigger a non-maskable interrupt (NMI)",
		LongHelp: longHelpNMI,
		Func: func(c *ishell.Context) {
			c.Println("Triggering non-maskable interrupt (NMI) ...")
			cmd <- common.NMI
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "reading",
		Help:     "set a special address for reading",
		LongHelp: longHelpReading,
		Func: func(c *ishell.Context) {
			c.Println("CLI: DUMMY reading")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "reset",
		Help:     "rest the machine (RESET)",
		LongHelp: longHelpReset,
		Func: func(c *ishell.Context) {
			c.Println("Triggering RESET signal ...")
			cmd <- common.RESET
		},
	})

	// TODO this doesn't work at all, see RUN
	shell.AddCmd(&ishell.Cmd{
		Name:     "resume",
		Help:     "resume after a halt",
		LongHelp: longHelpResume,
		Func: func(c *ishell.Context) {
			c.Println("Resuming at current PC location ...")
			cmd <- common.RESUME
		},
	})

	// TODO make this work quickly
	shell.AddCmd(&ishell.Cmd{
		Name:     "run",
		Help:     "run a machine created with 'init'",
		LongHelp: longHelpRun,
		Func: func(c *ishell.Context) {
			// TODO make sure we have an initialized machine
			// TODO print message about interfacing with machine
			fmt.Println("(CLI: DUMMY run)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "save",
		Help:     "save an address range to a file",
		LongHelp: longHelpSave,
		Func: func(c *ishell.Context) {
			c.Println("(CLI: DUMMY save)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "set",
		Help:     "set various parameters",
		LongHelp: longHelpSet,
		Func: func(c *ishell.Context) {
			c.Println("(CLI: DUMMY set)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "show",
		Help:     "display information on various parts of the system",
		LongHelp: longHelpShow,
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("ERROR: Need an argument")
			} else {
				subcmd := c.Args[0]

				switch subcmd {
				case "breakpoints":
					c.Println("CLI: SHOW: DUMMY show breakpoints")
				case "config":
					c.Println("(DUMMY show config)")
				case "memory":
					c.Println(memory.List())
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
		Help:     "display status of the machine or the host computer",
		LongHelp: longHelpStatus,
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				cmd <- common.STATUS
			} else {
				subcmd := c.Args[0]

				switch subcmd {
				case "host":
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
		Name:     "store",
		Help:     "store a byte at a given address",
		LongHelp: longHelpStore,
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY store)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "writing",
		Help:     "define a function to be triggered when address written to",
		LongHelp: longHelpWriting,
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY writing)")
		},
	})

	// TODO check for batch mode
	shell.Run()
	shell.Close()

}
