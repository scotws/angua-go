// Angua - An Emulator for the 65816 CPU in Native Mode
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 26. Sep 2017
// This version: 06. Jan 2019

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
	"reflect"
	"runtime"
	"strings"
	"unicode"

	"angua/common"
	"angua/cpu"
	"angua/info"
	"angua/mem"
	"angua/specials"

	"gopkg.in/abiosoft/ishell.v2" // https://godoc.org/gopkg.in/abiosoft/ishell.v2
)

const (
	CLEAR = 0 // Convenience for flags
	SET   = 1 // Convenience for flags

	configDir   string = "configs"
	shellBanner string = `Welcome to Angua
An Emulator for the 65816 in Native Mode
Version ALPHA 0.1  06. Jan 2019
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
func readConfig(s string) []string {
	var commands []string
	config := configDir + string(os.PathSeparator) + s

	configFile, err := os.Open(config)
	if err != nil {
		fmt.Println("Couldn't open config file", err)
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

	memory := &mem.Memory{
		// Remember to initialize the special address maps here or
		// machine
		SpecRead:  make(map[common.Addr24]func() (byte, bool)),
		SpecWrite: make(map[common.Addr24]func(byte)),
	}
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
		//	dump <ADDRESS_RANGE>
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
					// TODO make this a list with a pointer;
					// consider taking second parameter for
					// depth
					addrSP := common.Addr24(cpu.SP)
					hexDump(addrSP, addrSP+0xF, memory)

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
			a1, a2, err := parseAddressRange(c.Args)
			if err != nil {
				c.Println("Can't parse address range:", err)
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
			// TODO set up external terminal access

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
		//	load <FILENAME> [ to ] <ADDRESS>
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
			num, err := common.ConvertNum(addrString)
			if err != nil {
				c.Printf("Couldn't convert number %s: %v ", addrString, err)
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

			err = memory.BurnBlock(addr, data)
			if err != nil {
				c.Println(err)
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
		//	 memory <ADDRESS RANGE> ["is"] ("rom" | "ram")
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
			a1, a2, err := parseAddressRange(c.Args)
			if err != nil {
				c.Println("Can't parse address range:", err)
				return
			}

			nc, err := mem.NewChunk(a1, a2, memType)
			if err != nil {
				c.Println("Error: creating chunk with", a1, "to", a2, "as", memType, err)
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
			// The specials name table takes a lower case string
			// (for example "getchar") and returns the related
			// function (GetChar). To add your own special actions,
			// create a function and add the related string and
			// function to this table.
			SpecReadNames := map[string]func() (byte, bool){
				"getchar":        specials.GetChar,
				"getchar-blocks": specials.GetCharBlocks,
			}

			// We need at least two arguments
			if len(c.Args) < 2 {
				c.Println("ERROR: Need at least address and function name")
				return
			}

			addrString := c.Args[0]

			// Skip any "from" if there is one
			if c.Args[0] == "from" {
				addrString = c.Args[1]
			}

			// Convert the address string to an address
			ui, err := common.ConvertNum(addrString)
			if err != nil {
				c.Printf("Could't convert address %s: %v", addrString, err)
				return
			}

			addr := common.Addr24(ui)

			// Make sure address is not already in SpecRead
			_, ok := memory.SpecRead[addr]
			if ok {
				c.Println("ERROR: Address", addrString, "already defined.")
				return
			}

			// Get the function name as the last part of the line.
			// We allow user to have more than one address pointing
			// to the same function
			funcString := c.Args[len(c.Args)-1]

			// Store the information SpecReadNames
			memory.SpecRead[addr] = SpecReadNames[funcString]

		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "reset",
		Help:     "reset the machine (trigger Reset signal)",
		LongHelp: longHelpReset,
		Func: func(c *ishell.Context) {
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
	// TODO see about replacing by RESUME
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

				case "specs", "special", "specials":

					if !haveMachine {
						c.Println("No machine present (use 'init')")
						return
					}

					for a, f := range memory.SpecRead {
						fn := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
						c.Println("Reading:", a.HexString(), "calls", fn)
					}

					for a, f := range memory.SpecWrite {
						fn := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
						c.Println("Writing:", a.HexString(), "calls", fn)
					}

				case "sys", "system":
					c.Println("Use 'status host' to show host information")

				case "v", "vecs", "vectors", "interrupts":

					if !haveMachine {
						c.Println("No machine present (use 'init')")
						return
					}

					for _, vecData := range common.Vectors {
						vec, err := getVector(vecData.Addr, memory)
						if err != nil {
							fmt.Printf("Can't get vector for %s at %s:%v ", vecData.Name, vecData.Addr.HexString(), err)
							return
						}

						c.Printf("%-5s (%s): %s\n", vecData.Name, vecData.Addr.HexString(), vec)
					}

				default:
					c.Println("Option", subcmd, "unknown")
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

	// TODO this word calls BURN and will let you edit ROM
	shell.AddCmd(&ishell.Cmd{
		Name:     "store",
		Help:     "store a byte at a given address of RAM or ROM",
		LongHelp: longHelpStore,
		Func: func(c *ishell.Context) {
			c.Println("CLI: DUMMY store")
		},
	})

	// Format for writing line is
	//	writing [to] <ADDRESS> calls <FUNCNAME>
	shell.AddCmd(&ishell.Cmd{
		Name:     "writing",
		Help:     "define a function to be triggered when address written to",
		LongHelp: longHelpWriting,
		Func: func(c *ishell.Context) {

			// The specials name tables take a lower case string
			// (for example "getchar") and returns the related
			// function (GetChar). To add your own special actions,
			// create a function and add the related string and
			// function to this table.
			SpecWriteNames := map[string]func(byte){
				"putchar": specials.PutChar,
			}

			// We need at least two arguments
			if len(c.Args) < 2 {
				c.Println("ERROR: Need at least address and function name")
				return
			}

			addrString := c.Args[0]

			// Skip any "to" if there is one
			if c.Args[0] == "to" {
				addrString = c.Args[1]
			}

			// Convert the address string to an address
			ui, err := common.ConvertNum(addrString)
			if err != nil {
				c.Printf("Could't convert address %s: %v", addrString, err)
				return
			}

			addr := common.Addr24(ui)

			// Make sure address is not already in SpecWriteName
			_, ok := memory.SpecWrite[addr]
			if ok {
				c.Println("ERROR: Address", addrString, "already defined.")
				return
			}

			// Get the function name as the last part of the line.
			// We allow user to have more than one address pointing
			// to the same function
			funcString := c.Args[len(c.Args)-1]

			// Store the information SpecWritingName
			memory.SpecWrite[addr] = SpecWriteNames[funcString]
		},
	})

	// TODO check for batch mode

	shell.Run()
	shell.Close()

}

// getVectorString takes a 24 bit address and returns a string with the 16 bit
// address of the vector at that memory location in the first bank. This is a
// helper function for "show vectors"

func getVector(addr common.Addr24, m *mem.Memory) (string, error) {

	av, err := m.FetchMore(addr, 2)
	if err != nil {
		return "", fmt.Errorf("getVector: Can't get vector %s: %v", addr, err)
	}

	// FetchMore returns an int
	avAddr := common.Addr24(av)

	return avAddr.HexString(), nil
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

		b, err := m.Fetch(i)
		if err != nil {
			fmt.Println("hexDump:", err)
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
//	[<BANK>[:]]<ADDR16> "to" [<BANK>[:]]<ADDR16>
//	[<BANK>[:]]<ADDR16> [<BANK>[:]]<ADDR16>
//	"bank" <BYTE>
//

// and returns two addresses in the common.Addr24 format and error message.  or
// failure. This function lives here and not in common because it is part of the
// command line interface
func parseAddressRange(ws []string) (addr1, addr2 common.Addr24, err error) {
	// If the first word is "bank", then we are getting a full bank
	if ws[0] == "bank" {

		// Second word must be the bank byte. We brutally cut off
		// everything but the lowest byte
		bankNum, err := common.ConvertNum(ws[1])
		bankByte := common.Addr24(bankNum).Lsb()
		bankAddr := common.Addr24(bankByte) * 0x10000
		addr1 = bankAddr
		addr2 = bankAddr + 0xFFFF

		return addr1, addr2, err
	}

	// We at least need two addresses, so that's two words length. We could
	// parse more carefully, but not at the moment
	if len(ws) < 2 {
		return 0, 0, fmt.Errorf("parseAddrRange: wrong number of parameters")
	}

	num, err := common.ConvertNum(ws[0])
	if err != nil {
		return 0, 0, fmt.Errorf("parseAddrRange: can't convert number: %v", err)
	}

	addr1 = common.Addr24(num)

	// We allow people to slide on the "to" though we don't
	// advertise the fact. Later, once we have the error handling of
	// ConvertNum working, check ws[1] and if there is an error, skip
	// to ws[2].
	if ws[1] == "to" {
		num, err = common.ConvertNum(ws[2])
	} else {
		num, err = common.ConvertNum(ws[1])

	}

	addr2 = common.Addr24(num)

	return addr1, addr2, nil
}

// printCPUStatus print information on the registers, flags and other important
// CPU data. This assumes that the machine has been halted or interesting things
// might happen
// TODO rewrite this is arrays of strings to get rid of the IFs when we're sure
// of what everything is supposed to look like, also get rid of concat operation
// for bytearrays
func printCPUStatus(c *cpu.CPU) {

	// --- Print legend --------------------------------------

	fmt.Print(" PC  PB ")

	// Accumulator: M=1 is 8 bit (and B register), M=0 is 16 bit
	if c.FlagM == SET {
		fmt.Print(" A  B ")
	} else {
		fmt.Print("  A  ")
	}

	// XY Registers: X=1 is 8 bit, X=0 is 16 bit
	if c.FlagX == SET {
		fmt.Print(" X  Y ")
	} else {
		fmt.Print("  X    Y  ")
	}

	fmt.Println("DB  DP   SP  NVMXDIZC E")

	// --- Print data --------------------------------------

	fmt.Print(c.PC.HexString(), " ", c.PBR.HexString())

	// Accumultor: M=1 is 8 bit (and B register), M=0 is 16 bit
	if c.FlagM == SET {
		fmt.Print(" ", c.A8.HexString(), " ", c.B.HexString())
	} else {
		fmt.Print(" ", c.A16.HexString())
	}

	// XY Registers: X=1 is 8 bit, X=0 is 16 bit
	if c.FlagX == SET {
		fmt.Print(" ", c.X8.HexString(), " ", c.Y8.HexString(), " ")
	} else {
		fmt.Print(" ", c.X16.HexString(), " ", c.Y16.HexString(), " ")
	}

	fmt.Print(c.DBR.HexString(), " ",
		c.DP.HexString(), " ",
		c.SP.HexString(), " ",
		c.StringStatReg(), " ",
		c.FlagE)

	// It doesn't make sense to print "halted" because we will usually halt
	// before printing the CPU state
	if c.IsWaiting {
		fmt.Print(" waiting")
	}

	if c.IsStopped {
		fmt.Print(" stopped")
	}

	fmt.Println()

	return
}
