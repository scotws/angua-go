// Angua - An Emulator for the 65816 CPU in Native Mode
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 26. Sep 2017
// This version: 04. Jan 2019

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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unicode"

	"angua/common"
	"angua/cpu"
	"angua/info"
	"angua/mem"

	"gopkg.in/abiosoft/ishell.v2" // https://godoc.org/gopkg.in/abiosoft/ishell.v2
)

const (
	configDir   string = "configs"
	shellBanner string = `Welcome to Angua
An Emulator for the 65816 in Native Mode
Version ALPHA 0.1  02. Jan 2019
Copyright (c) 2018-2019 Scot W. Stevenson
Angua comes with absolutely NO WARRANTY
Type 'help' for more information
`
)

var (
	haveMachine bool = false

	// Flags passed.
	// TODO Add "-d" to print debug information
	beVerbose   = flag.Bool("v", false, "Verbose, print more output")
	inBatchMode = flag.Bool("b", false, "Start in batch mode") // TODO see if this should be R for "run"
	configFile  = flag.String("c", "default.cfg", "Configuration file")
)

// readConfig takes the address of a configuration file in the form of
// "configs/<NAME>.cfg" and reads the content, returning it stripped of empty
// lines and comments as a list of strings. We use bufio.Scanner here because we
// want to read the test as lines
// TODO pass bool back to check if stuff went wrong, removing log.Fatal()
func readConfig(s string) []string {
	var commands []string
	config := configDir + string(os.PathSeparator) + s

	configFile, err := os.Open(config)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	source := bufio.NewScanner(configFile)

	for source.Scan() {
		line := strings.TrimSpace(source.Text())

		if line == "" || line[0] == ';' {
			continue
		}

		commands = append(commands, line)
	}

	return commands
}

// verbose takes a string and prints it on the standard output through logger if
// the user wants us to be verbose
func verbose(s string) {
	if *beVerbose {
		log.Println(s)
	}
}

// -----------------------------------------------------------------

func main() {

	// Generate Dictionaries for 'info' system in the background
	go info.GenerateDicts()

	flag.Parse()

	memory := &mem.Memory{}
	cpu := &cpu.CPU{}

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
		Name:     "destroy",
		Help:     "de-initialize a machine (start over)",
		LongHelp: longHelpDestroy,
		Func: func(c *ishell.Context) {
			if !haveMachine {
				c.Println("ERROR: No machine present")
				return
			}

			c.Println("CLI: DUMMY: destroy the machine")
			haveMachine = false
			cmd <- common.HALT

			err := shell.Process("beep")
			if err != nil {
				log.Fatal(err)
			}
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
		// Format is either
		//
		// 	dump <ADDRESS_RANGE>
		//	dump stack
		//	dump dp
		Name:     "dump",
		Help:     "print hex dump of a memory range",
		LongHelp: longHelpDump,
		Func: func(c *ishell.Context) {

			const ruler string = "         00 01 02 03 04 05 06 07  08 09 0A 0B 0C 0D 0E 0F"

			if !haveMachine {
				c.Println("ERROR: No machine present")
				return
			}

			if len(c.Args) == 0 {
				c.Println("ERROR: Need at least one argument (see 'dump help')")
				return
			}

			if len(c.Args) == 1 {

				switch c.Args[0] {

				case "sp", "stack", "stackpointer":
					// TODO make this a list with a pointer,
					// consider taking second parameter for
					// depth
					addrDP := common.Addr24(cpu.SP)
					hexDump(addrDP, addrDP+0xF, memory)

				case "dp", "direct", "directpage":
					addrDP := common.Addr24(cpu.DP)
					fmt.Println(ruler)
					hexDump(addrDP, addrDP+0xFF, memory)

				default:
					fmt.Println("ERROR: Unknown option:", c.Args[0])
				}
				return
			}

			// The arg string is passed without "dump"
			a1, a2, ok := parseAddressRange(c.Args)
			if !ok {
				c.Println("ERROR parsing address range")
				return
			}
			hexDump(a1, a2, memory)
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
			printCPUStatus(cpu) // TODO see if we are really giving current status
		},
	})

	shell.AddCmd(&ishell.Cmd{
		// TODO allow reading new configuration from a file given as a
		// parameter
		Name:     "init",
		Help:     "initialize a new machine",
		LongHelp: longHelpInit,
		Func: func(c *ishell.Context) {
			var commands []string

			if haveMachine {
				c.Println("ERROR: Machine already initialized")
				return
			}

			commands = readConfig(*configFile)

			// Process configuration file
			for _, cmd := range commands {
				ws := strings.Fields(cmd)

				err := shell.Process(ws[0:]...)
				if err != nil {
					log.Fatal(err)
				}
			}

			c.Println("Processed configuration file", *configFile, "...")

			// TODO check for chunk overlap in memory
			// TODO set up special memory

			// Set up CPU
			cpu.Mem = memory

			c.Println("Initializing machine ...")
			haveMachine = true
			cpu.IsHalted = true

			c.Println("System initialized, start with 'reset'")
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
		// Load a binary file to memory. Memory has to already exist.
		// Format is
		//
		// 	load <FILENAME> [ to ] <ADDRESS>
		//
		// The name of the file is not in quotation marks
		Name:     "load",
		Help:     "load binary file to memory",
		LongHelp: longHelpLoad,
		Func: func(c *ishell.Context) {

			if !haveMachine {
				c.Println("Machine must be initialized first by 'init'")
				return
			}

			// The filename must be the first parameter and the
			// address the last one. This lets us ignore any "to" in
			// the middle. One way or another, we need at least two
			// parameters
			if len(c.Args) < 2 {
				c.Println("ERROR: Need filename and address to load.")
				return
			}

			fileName := c.Args[0]
			addrString := c.Args[len(c.Args)-1]
			num, ok := common.ConvertNum(addrString)
			if !ok {
				c.Println("ERROR: Couldn't convert number", addrString)
				return
			}

			addr := common.Addr24(num)

			// We can use ioutil.ReadFile here because we want
			// binary data
			data, err := ioutil.ReadFile(fileName)
			if err != nil {
				c.Println(err)
				return
			}

			ok = memory.BurnBlock(addr, data)
			if !ok {
				c.Println("ERROR: Couldn't write binary data to address")
				return
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "info",
		Help:     "print information on 65816 instructions",
		LongHelp: longHelpInfo,
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("ERROR: Need opcode, SAN mnemonic, or 'all'")
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
				// another command such as 'all'

				c.Println("ERROR: Opcode or mnemonic", subcmd, "unknown")
			}

		},
	})

	shell.AddCmd(&ishell.Cmd{
		// Memory commands take the form of
		// 	 memory <ADDRESS RANGE> ["is"] ("rom" | "ram")
		//       memory bank <NUMBER> ["is"] ("rom" | "ram")
		//       memory
		Name:     "memory",
		Help:     "define a memory chunk",
		LongHelp: longHelpMemory,
		Func: func(c *ishell.Context) {

			// If we weren't given any parameters, we just print the
			// current memory configuration. This is the same
			// command as "show memory"
			if len(c.Args) == 0 {
				c.Println(memory.List())
				return
			}

			// We can break this up from the end by making sure that the
			// last word is either "ram" or "rom", and feeding the beginning
			// to the address range finder. Note that mem.NewChunk
			// tests if the sting passed is either "rom" or "ram" so
			// we don't have to do that here
			memType := c.Args[len(c.Args)-1]

			// The arg string is passed without "memory"
			a1, a2, ok := parseAddressRange(c.Args)
			if !ok {
				c.Println("ERROR parsing address range")
				return
			}

			nc, ok := mem.NewChunk(a1, a2, memType)
			if !ok {
				c.Println("ERROR creating chunk with", a1, "to", a2, "as", memType)
				return
			}

			memory.Chunks = append(memory.Chunks, nc)
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
			shell.Process("run")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "resume",
		Help:     "resume after a halt",
		LongHelp: longHelpResume,
		Func: func(c *ishell.Context) {
			c.Println("Resuming at current PC location ...")
			cmd <- common.RESUME
		},
	})

	// Run is only used after HALT. To start the machine, use INIT followed
	// by RESET. After INIT, the PC is not in a defined state.
	shell.AddCmd(&ishell.Cmd{
		Name:     "run",
		Help:     "run a machine created with 'init'",
		LongHelp: longHelpRun,
		Func: func(c *ishell.Context) {

			if !haveMachine {
				c.Println("Machine must be initialized first")
				return
			}

			// TODO print message about interfacing with machine
			c.Println("Running machine ...")
			go cpu.Run(cmd)
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
			c.Println("CLI: DUMMY set")
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
					c.Println("CLI: DUMMY show config")
				case "memory": // This is the same as just calling memory
					c.Println(memory.List())
				case "specials":
					c.Println("CLI: DUMMY show specials")
				case "system":
					c.Println("Use 'status host' to show host information")
				case "vectors":
					c.Println("CLI: DUMMY show vectors")
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

			// "status host" prints data on the host computer,
			// including number of goroutines
			if len(c.Args) == 1 && c.Args[0] == "host" {
				fmt.Println("Host architecture:", runtime.GOARCH)
				fmt.Println("Host operating system:", runtime.GOOS)
				fmt.Println("Host system CPU cores available:", runtime.NumCPU())
				fmt.Println("Host system goroutines running:", runtime.NumGoroutine())
				fmt.Println("Host system Go version:", runtime.Version())

				return
			}

			if !haveMachine {
				c.Println("No machine present (use 'init')")
				return
			}

			// TODO remember if machine is halted or not, then halt
			// machine

			// If we were told to print all data with "status all",
			// well, print all data
			if len(c.Args) == 1 && c.Args[0] == "all" {

				if cpu.IsHalted {
					fmt.Println("Machine is halted.")
				}

				fmt.Println("Total memory (ROM and RAM):", memory.Size(), "bytes")
				printCPUStatus(cpu)

				return
			}

			// In all other cases, we just print the CPU status
			printCPUStatus(cpu)

			// TODO restore old status

			return
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
			c.Println("CLI: DUMMY writing")
		},
	})

	// TODO check for batch mode
	shell.Run()
	shell.Close()

}

// Hexdump prints the contents of a memory range in a nice hex table. If the
// addresses do not exist, we just print a zero without any fuss. We could use
// the library encoding/hex for this, but we want to print the first address of
// the line, and the library function starts the count with zero, not the
// address. Also, we want uppercase letters for hex values. This is kept here
// because it is part of the command-line interface
func hexDump(addr1, addr2 common.Addr24, m *mem.Memory) {
	var r rune
	var count uint
	var hb strings.Builder // hex part
	var cb strings.Builder // char part
	var template string = "%-58s%s\n"

	for i := addr1; i <= addr2; i++ {

		// The first run produces a blank line because this if is
		// triggered, however, the strings are empty because of the way
		// Go initializes things
		if count%16 == 0 {
			fmt.Printf(template, hb.String(), cb.String())
			hb.Reset()
			cb.Reset()

			nextAddr := addr1 + common.Addr24(count)
			fmt.Fprintf(&hb, nextAddr.HexString()+" ")
		}

		b, ok := m.Fetch(i)
		if !ok {
			fmt.Println("ERROR: No memory at address", i.HexString(), "(see 'memory')")
			return
		}

		// Build the hex string
		fmt.Fprintf(&hb, " %02X", b)

		// Build the string list. This is the 21. century so we hexdump
		// in Unicode, not ASCII, though this doesn't make a different
		// if we just have byte values
		r = rune(b)
		if !unicode.IsPrint(r) {
			r = rune('.')
		}

		fmt.Fprintf(&cb, string(r))
		count += 1

		// We put one extra blank line after the first eight entries to
		// make the dump more readable
		if count%8 == 0 {
			fmt.Fprintf(&hb, " ")
		}

	}

	// If the loop is all done, we might still have stuff left in the
	// buffers
	fmt.Printf(template, hb.String(), cb.String())
}

// parseAddressRange takes a list of strings that in the formats of
//
// 	[<BANK>[:]]<ADDR16> "to" [<BANK>[:]]<ADDR16>
// 	[<BANK>[:]]<ADDR16> [<BANK>[:]]<ADDR16>
// 	"bank" <BYTE>
//
// and returns two addresses in the common.Addr24 format and a bool for success
// or failure. This function lives here and not in common because it is part of
// the command line interface
func parseAddressRange(ws []string) (addr1, addr2 common.Addr24, ok bool) {
	ok = true

	// If the first word is "bank", then we are getting a full bank
	if ws[0] == "bank" {

		// Second word must be the bank byte. We brutally cut off
		// everything but the lowest byte
		bankNum, ok := common.ConvertNum(ws[1])
		bankByte := common.Addr24(bankNum).Lsb()
		bankAddr := common.Addr24(bankByte) * 0x10000
		addr1 = bankAddr
		addr2 = bankAddr + 0xFFFF

		return addr1, addr2, ok
	}

	// We at least need two addresses, so that's two words length. We could
	// parse more carefully, but not at the moment
	if len(ws) < 2 {
		addr1 = 0
		addr2 = 0
		ok = false
		return addr1, addr2, ok
	}

	num, ok := common.ConvertNum(ws[0])
	addr1 = common.Addr24(num)

	// We allow people to slide on the "to" though we don't
	// advertise the fact. Later, once we have the error handling of
	// ConvertNum working, check ws[1] and if there is an error, skip
	// to ws[2].
	if ws[1] == "to" {
		num, ok = common.ConvertNum(ws[2])
	} else {
		num, ok = common.ConvertNum(ws[1])

	}

	addr2 = common.Addr24(num)

	return addr1, addr2, ok
}

// printCPUStatus print information on the registers, flags and other important
// CPU data. This assumes that the machine has been halted or interesting things
// might happen
func printCPUStatus(c *cpu.CPU) {

	// --- Print legend --------------------------------------

	fmt.Print(" PC  PB ")

	// Accumulator: M=1 is 8 bit (and B register), M=0 is 16 bit
	if c.FlagM {
		fmt.Print(" A  B ")
	} else {
		fmt.Print("  A  ")
	}

	// XY Registers: X=1 is 8 bit, X=0 is 16 bit
	if c.FlagX {
		fmt.Print(" X  Y ")
	} else {
		fmt.Print("  X    Y  ")
	}

	fmt.Println("DB  DP   SP  NVMXDIZC")

	// --- Print data --------------------------------------

	fmt.Print(c.PC.HexString(), " ", c.PBR.HexString())

	// Accumultor: M=1 is 8 bit (and B register), M=0 is 16 bit
	if c.FlagM {
		fmt.Print(" ", c.A8.HexString(), " ", c.B.HexString())
	} else {
		fmt.Print(" ", c.A16.HexString())
	}

	// XY Registers: X=1 is 8 bit, X=0 is 16 bit
	if c.FlagX {
		fmt.Print(" ", c.X8.HexString(), " ", c.Y8.HexString(), " ")
	} else {
		fmt.Print(" ", c.X16.HexString(), " ", c.Y16.HexString(), " ")
	}

	p := c.GetStatusReg()
	statusReg := strconv.FormatUint(uint64(p), 2)

	fmt.Println(c.DBR.HexString(), c.DP.HexString(), c.SP.HexString(), statusReg)

	return
}
