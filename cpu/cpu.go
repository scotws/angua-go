// Angua CPU System
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 06. Nov 2018
// First version: 26. Dec 2018

package cpu

import (
	"fmt"
	"log"
	"time"

	"angua/common"
)

const (
	// Interrupt vectors. Note the reset vector is only for emulated mode.
	// See http://6502.org/tutorials/65c816interrupts.html for details
	abortAddr = 0xFFE8
	brkAddr   = 0xFFE6
	copAddr   = 0xFFE4
	irqAddr   = 0xFFEE
	nmiAddr   = 0xFFEA
	resetAddr = 0xFFFC // Routine must move to native mode ASAP

	// Width of accumulator and registers
	A8   = 0
	A16  = 1
	XY   = 0
	XY16 = 1
)

// --------------------------------------------------
// Status Register

// D for decimal mode and E for emulated are provided, but not functional
// If you are coming from the 6502, notice that the BRK instruction is handled
// differently

type StatReg struct {
	FlagN bool // Negative flag, true if highest bit is 1
	FlagV bool // Overflow flag, true if overflow
	FlagB bool // Break Instruction, true if interrupt by BRK
	FlagD bool // Decimal mode, true is decimal, false is binary
	FlagI bool // Interrupt disable, true is disabled
	FlagZ bool // Zero flag
	FlagC bool // Carry flag

	FlagE bool // Emulation flag, set signals is 6502 emulation mode
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

// TestN takes a int and sets the N flag to true if highest bit is set else to flase
func (s *StatReg) TestN(i int) {
	// TODO check based on register size
}

// --------------------------------------------------
// CPU

type CPU struct {
	A8  common.Data8  // Accumulator 8 bit
	A16 common.Data16 // Accumulator 16 bit
	B   common.Data8  // B register (always 8 bit)
	X8  common.Data8  // X register 8 bit
	X16 common.Data16 // X register 16 bit
	Y8  common.Data8  // Y register 8 bit
	Y16 common.Data16 // Y register 16 bit

	DP  common.Data16 // Direct Page register, yes, 16 bit, not 8
	SP  common.Data16 // Stack Pointer, 16 bit
	P   byte          // Status Register
	DBR common.Data8  // Data Bank Register
	PBR common.Data8  // Program Bank Register
	PC  common.Data16 // Program counter

	ModeA  int // Current width of Accumulator, either A16 or A8
	ModeXY int // Current width of X and Y registers, either XY16 or X8

	IsHalted       bool // Signals if CPU stopped by Angua CLI
	IsWaiting      bool // CPU stopped by WAI instruction
	IsStopped      bool // CPU stopped by STP instruction
	SingleStepMode bool // Signals if we are in single step mode

	StatReg
}

// Step executes a single instruction from PC. This is called by the Run method
func (c *CPU) Step() {
	fmt.Println("CPU: DUMMY: <EXECUTING ONE INSTRUCTION>")
}

// Run is the main loop of the CPU.
func (c *CPU) Run(cmd <-chan int) {

	fmt.Println("CPU: DUMMY: Run")
	c.IsHalted = false  // User freezes execution, resume with 'resume'
	c.IsStopped = false // STP instruction
	c.IsWaiting = false // WAI instruction

	for {
		// If we are not halted, we run the main CPU loop: See if we
		// received a command from the CLI; if not, single step an
		// instruction.
		for !c.IsHalted {

			select {
			case order := <-cmd:
				// If we were given a command by the operating system,
				// we execute it first
				switch order {

				case common.ABORT:
					fmt.Println("CPU: DUMMY: Received cmd ABORT")

				case common.BOOT:
					fmt.Println("CPU: DUMMY: Received cmd BOOT")

				case common.HALT:
					fmt.Println("CPU: DUMMY: Received *** HALT ***")
					c.IsHalted = true

				case common.IRQ:
					fmt.Println("CPU: DUMMY: Received cmd IRQ")

				case common.NOTRACE:
					fmt.Println("CPU: DUMMY: Received cmd NOTRACE")
					trace = false

				case common.NOVERBOSE:
					fmt.Println("CPU: DUMMY: Received cmd NOVERBOSE")
					verbose = false

				case common.NMI:
					fmt.Println("CPU: DUMMY: Received cmd NMI")

				case common.RESET:
					fmt.Println("CPU: DUMMY: Received cmd RESET")

				case common.RESUME:
					fmt.Println("CPU: DUMMY: Received cmd RESUME")
					c.IsHalted = false

				case common.RUN:
					fmt.Println("CPU: DUMMY: Received cmd RUN")
					c.IsHalted = false

				case common.STATUS:
					fmt.Println("CPU: DUMMY: Received cmd STATUS")

				case common.STEP:
					fmt.Println("CPU: DUMMY: Received cmd STEP")
					c.SingleStepMode = true

				case common.TRACE:
					fmt.Println("CPU: DUMMY: Received cmd TRACE")
					trace = true

				case common.VERBOSE:
					fmt.Println("CPU: DUMMY: Received cmd VERBOSE")
					verbose = true

				default:
					// TODO make this less brutal
					log.Fatal("ERROR: CPU: Got unknown command", order, "from CLI")

				}

			default:
				// This is where the CPU actually runs an
				// instruction
				c.Step()
				fmt.Println("CPU: DUMMY: Main loop")
				time.Sleep(2 * time.Second)
			}
		}
	}

}
