// Angua CPU System
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 06. Nov 2018
// First version: 18. Nov 2018

package cpu

import (
	"fmt"
	"log"
	"time"

	"angua/common"
)

const (
	// Interrupt vectors TODO check
	irqAddr   = 0xFFFE
	resetAddr = 0xFFFC
	nmiAddr   = 0xFFFA
	copAddr   = 0xFFF4

	// Width of accumulator and registers. We would love to use bools or
	// such, but this is Go
	A8   = 0
	A16  = 1
	XY   = 0
	XY16 = 1
)

// --------------------------------------------------
// Status Register

// D for decimal mode and E for emulated are provided, but not functional

type StatReg struct {
	FlagN bool // Negative flag, true if highest bit is 1
	FlagV bool // Overflow flag, true if overflow
	FlagB bool // Break Instruction, true if interrupt by BRK
	FlagD bool // Decimal mode, true is decimal, false is binary
	FlagI bool // Interrupt disable, true is disabled
	FlagZ bool // Zero flag
	FlagC bool // Carry bit

	FlagE bool // Emulation, true is 6502 emulation mode
}

var (
	cmd = make(<-chan int, 2) // Receive commands from CLI

	verbose bool // Print lots of information
	trace   bool // Print even more information
)

// GetStatusReg creates a status byte out of the flags of the Status Register
// and returns it to the caller. It is used by the instuction PHP for example.
// TODO code this
func (s *StatReg) GetStatusReg() byte {
	fmt.Println("DUMMY GetStatusRegister")
	return 0xFF // TODO dummy
}

// SetStatusReg takes a byte and sets the flags of the Status Register
// accordingly. It is used by the instruction PLP for example.
// TODO code this
func (s *StatReg) SetStatusReg(b byte) {
	fmt.Println("DUMMY SetStatusRegister")
}

// TestZ takes an int and sets the Z flag to true if the value is zero and to
// false otherwise
func (s *StatReg) TestZ(i int) {
	if i == 0 {
		s.FlagZ = true
	} else {
		s.FlagZ = false
	}
}

// TestN takes a int and sets the N flag to true if bit 7 is one else to flase
func (s *StatReg) TestN(i int) {
	// TODO
}

// --------------------------------------------------
// CPU

type CPU struct {
	A8  common.Data8  // Accumulator 8 bit
	A16 common.Data16 // Accumulator 16 bit
	B   common.Data8  // B register always 8 bit
	X8  common.Data8  // X register 8 bit
	X16 common.Data16 // X register 16 bit
	Y8  common.Data8  // Y register 8 bit
	Y16 common.Data16 // Y register 16 bit
	DP  common.Data16 // Direct Page register, yes, 16 bit, not 8
	SP  common.Data16 // Stack Pointer, 16 bit
	P   byte          // Status Register
	DBR common.Data8  // Data Bank Register, yes, available in emulated mode
	PBR common.Data8  // Program Bank Register, yes, available in emulated mode
	PC  common.Data16 // Program counter

	ModeA  int // Current width of Accumulator, either A16 or A8
	ModeXY int // Current width of X and Y registers, either XY16 or X8

	Halted     bool // Signals if CPU stopped by CLI
	SingleStep bool // Signals if we are in single step mode

	StatReg
}

// Step executes a single instruction from PC. This is called by the Run method
func (c *CPU) Step() {
	fmt.Println("CPU: DUMMY: <EXECUTING ONE INSTRUCTION>")
}

// Run is the main loop of the Cpu8. It takes two channels from the CLI: A
// boolean which enables running the processor and blocks it when waiting for
// input (which means the other CPU is running or everything is halted).
func (c *CPU) Run(cmd <-chan int) {

	fmt.Println("CPU: DUMMY: Run")
	c.Halted = true

	// TODO REWRITE THIS WITHOUT MODES

	for {
		// If we are not halted, we run the main CPU loop: See if we
		// received a command from the CLI, if not, single step an
		// instruction.
		for !c.Halted {

			select {
			case order := <-cmd:
				// If we were given a command by the operating system,
				// we execute it first
				switch order {

				case common.HALT:
					fmt.Println("CPU: DUMMY: Received *** HALT ***")
					c.Halted = true

				case common.RESUME, common.RUN:
					fmt.Println("CPU: DUMMY: Received cmd RESUME/RUN")
					c.SingleStep = false
					c.Halted = false

				case common.STEP:
					fmt.Println("CPU: DUMMY: Received cmd STEP")
					c.SingleStep = true

				case common.STATUS:
					c.Status()

				case common.BOOT:
					fmt.Println("CPU: DUMMY: Received *** BOOT ***")

				case common.RESET:
					fmt.Println("CPU: DUMMY: Received cmd RESET")

				case common.IRQ:
					fmt.Println("CPU: DUMMY: Received cmd IRQ")

				case common.NMI:
					fmt.Println("CPU: DUMMY: Received cmd NMI")

				case common.ABORT:
					fmt.Println("CPU: DUMMY: Received cmd ABORT")

				case common.VERBOSE:
					fmt.Println("CPU: DUMMY: Received cmd VERBOSE")
					verbose = true

				case common.LACONIC:
					fmt.Println("CPU: DUMMY: Received cmd LACONIC")
					verbose = false

				case common.TRACE:
					fmt.Println("CPU: DUMMY: Received cmd TRACE")
					trace = true

				case common.NOTRACE:
					fmt.Println("CPU: DUMMY: Received cmd NOTRACE")
					trace = false

				default:
					log.Fatal("ERROR: CPU: Got unknown command", order, "from CLI")

				}

			default:
				// This is where the CPU actually runs an
				// instruction. We pretend for testing that at
				// some point we are told to switch
				c.Step()
				fmt.Println("CPU: DUMMY: Main loop")
				time.Sleep(10 * time.Second)
			}
		}
	}

}

func (c *CPU) Status() {
	fmt.Println("CPU: DUMMY: Request of status received")
}
