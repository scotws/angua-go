// Angua CPU System - Emulated Mode (8-bit) CPU
// Scot W. Stevenson
// First version: 06. Nov 2018
// First version: 09. Nov 2018

package cpu8

import (
	"fmt"

	"angua/common"
)

const (
	// Interrupt vectors
	irqAddr   = 0xFFFE
	resetAddr = 0xFFFC
	nmiAddr   = 0xFFFA
	copAddr   = 0xFFF4
)

type reg8 uint8
type reg16 uint16

// --------------------------------------------------
// Status Register

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
	enable8 = make(chan struct{})
	cmd     = make(chan int, 2)

	Verbose bool // Print lots of information
	Trace   bool // Print even more information
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

// TestZ takes a byte and sets the Z flag to true if the value is zero and to
// false otherwise
func (s *StatReg) TestZ(b byte) {
	if b == 0 {
		s.FlagZ = true
	} else {
		s.FlagZ = false
	}
}

// TestN takes a byte and sets the N flag to true if bit 7 is one else to flase
func (s *StatReg) TestN(b byte) {
	// TODO
}

// --------------------------------------------------
// CPU

type Cpu8 struct {
	A   reg8  // 8 bit accumulator
	B   reg8  // Special hidden register of the 65816
	X   reg8  // index register
	Y   reg8  // index register
	DP  reg16 // Direct Page register, yes, 16 bit, not 8
	SP  reg8  // Stack Pointer, 8 bit
	P   byte  // Status Register
	DBR reg8  // Data Bank Register, yes, available in emulated mode
	PBR reg8  // Program Bank Register, yes, available in emulated mode
	PC  reg16 // Program counter

	Halted     bool // Signals if CPU stopped by CLI
	SingleStep bool // Signals if we are in single step mode

	StatReg
}

// Step executes a single instruction from PC. This is called by the Run method
func (c *Cpu8) Step() {
	fmt.Println("CPU8: DUMMY: Step")
}

// Run is the main loop of the Cpu8. It takes two channels from the CLI: A
// boolean which enables running the processor and blocks it when waiting for
// input (which means the other CPU is running or everything is halted).
func (c *Cpu8) Run() {

	fmt.Println("CPU8: DUMMY: Run")

	for {
		// This channel is used to block the CPU until it receives the
		// signal to run again
		<-enable8
		fmt.Println("CPU8: DUMMY: CPU8 enabled by CLI")

		// If we have received a signal to run, then we're not halted
		c.Halted = false

		// If we are not halted, we run the main CPU loop: See if we
		// received a command from the CLI, if not, single step an
		// instruction.
		// TODO figure out single step mode
		for !c.Halted {

			select {
			case order := <-cmd:
				// If we were given a command by the operating system,
				// we execute it first
				// TODO add a switch that handles our input,
				// especially the HALT signal
				fmt.Println("CPU8: DUMMY: Received command", order, "from CLI")

				switch order {

				case common.HALT:
					fmt.Println("CPU8: DUMMY: Received cmd HALT")
					c.Halted = true
				case common.RESUME, common.RUN:
					fmt.Println("CPU8: DUMMY: Received cmd RESUME/RUN")
					SingleStep = false
					c.Halted = false

				case common.STEP:
					fmt.Println("CPU8: DUMMY: Received cmd STEP")
					StingleStep = true

				case common.STATUS:
					fmt.Println("CPU8: DUMMY: Received cmd STATUS")
					c.Status()

				case common.BOOT:
					fmt.Println("CPU8: DUMMY: Received cmd BOOT")
				case common.RESET:
					fmt.Println("CPU8: DUMMY: Received cmd RESET")
				case common.IRQ:
					fmt.Println("CPU8: DUMMY: Received cmd IRQ")
				case common.NMI:
					fmt.Println("CPU8: DUMMY: Received cmd NMI")
				case common.ABORT:
					fmt.Println("CPU8: DUMMY: Received cmd ABORT")

				case common.VERBOSE:
					fmt.Println("CPU8: DUMMY: Received cmd VERBOSE")
					verbose = true
				case common.LACONIC:
					fmt.Println("CPU8: DUMMY: Received cmd LACONIC")
					verbose = false
				case common.TRACE:
					fmt.Println("CPU8: DUMMY: Received cmd TRACE")
					trace = true
				case common.NOTRACE:
					fmt.Println("CPU8: DUMMY: Received cmd NOTRACE")
					trace = false

				default:
					log.Fatal("ERROR: cpu8: Got unknown command", order, "from CLI")

				}

			default:
				// This is where the CPU actually runs an
				// instruction
				c.Step()
			}
		}
	}

}

// Status prints the status of the machine
func (c *Cpu8) Status() {
	fmt.Println("CPU8: DUMMY: Status")
}

// Status prints the status of the machine
func (c *Cpu8) Status() {
	fmt.Println("CPU8: DUMMY: Status")
}

/*
// TODO test version to Execute one opcode
func (c *Cpu8) Execute(b byte) {
	opcodes8[b](c)
}
*/
