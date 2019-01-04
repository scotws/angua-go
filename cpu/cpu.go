// Angua CPU System
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 06. Nov 2018
// First version: 04. Jan 2019

package cpu

import (
	"fmt"
	"log"
	"time"

	"angua/common"
	"angua/mem"
)

const (
	// Interrupt vectors. Note the reset vector is only for emulated mode.
	// See http://6502.org/tutorials/65c816interrupts.html and Eyes & Lichty
	// p. 195 for details. We store these as 24 bit addresses because that
	// is the way we'll use them during mem.FetchMore
	abortAddr common.Addr24 = 0xFFE8
	brkAddr   common.Addr24 = 0xFFE6
	copAddr   common.Addr24 = 0xFFE4
	irqAddr   common.Addr24 = 0xFFEE
	nmiAddr   common.Addr24 = 0xFFEA
	resetAddr common.Addr24 = 0xFFFC // <-- This is the really important one

	// Width of accumulator and registers
	A8   int = 0
	A16  int = 1
	XY   int = 0
	XY16 int = 1
)

// --------------------------------------------------
// Status Register

// D for decimal mode and E for emulated are provided, but not functional
// If you are coming from the 6502, notice that the BRK instruction is handled
// differently

type StatReg struct {
	FlagN bool // Negative flag, true if highest bit is 1
	FlagV bool // Overflow flag, true if overflow
	FlagM bool // A size, set is 8 bit (and B register), clear is 16 bit
	FlagX bool // XY size, set is 8 bit, clear is 16 bit
	FlagD bool // Decimal mode, true is decimal, false is binary
	FlagI bool // Interrupt disable, true is disabled
	FlagZ bool // Zero flag
	FlagC bool // Carry flag

	FlagB bool // Break Instruction, true if interrupt by BRK (6502)
	FlagE bool // Emulation flag, set signals is 6502 emulation mode (unused)
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
	return 0xFF // TODO dummy
}

// SetStatusReg takes a byte and sets the flags of the Status Register
// accordingly. It is used by the instruction PLP for example.
// TODO code this
func (s *StatReg) SetStatusReg(b byte) {
	fmt.Println("CPU: DUMMY: SetStatusRegister")
}

// TestZ takes an int and sets the Z flag to true if the value is zero and to
// false otherwise
// TODO get serious about this code
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

	DP  common.Addr16 // Direct Page register, yes, 16 bit, not 8
	SP  common.Addr16 // Stack Pointer, 16 bit
	P   byte          // Status Register
	DBR common.Data8  // Data Bank Register
	PBR common.Data8  // Program Bank Register
	PC  common.Addr16 // Program counter

	ModeA  int // Current width of Accumulator, either A16 or A8
	ModeXY int // Current width of X and Y registers, either XY16 or X8

	IsHalted       bool // Signals if CPU stopped by Angua CLI
	IsWaiting      bool // CPU stopped by WAI instruction
	IsStopped      bool // CPU stopped by STP instruction
	SingleStepMode bool // Signals if we are in single step mode

	Mem *mem.Memory // Pointer to the memory we're working on

	StatReg
}

// getFullAddr merges the Bank Byte and the Program Counter and returns a 24
// bit address
func (c *CPU) getFullPC() common.Addr24 {
	bank := common.Addr24(c.PBR) << 16
	addr := bank + common.Addr24(c.PC)

	return common.Ensure24(addr)
}

// Step executes a single instruction from PC. This is called by the Run method
// TODO this is pretty much all fake
func (c *CPU) Step() {

	// Get byte at PC
	ins, ok := c.Mem.Fetch(c.getFullPC())
	if !ok {
		log.Println("ERROR: Can't get instruction at", c.getFullPC().HexString())
		return
	}

	// Execute the instruction by accessing the entry in the Instruction
	// Jump table. We pass a pointer to the CPU struct.
	InsJump[ins](c)

	c.PC += common.Addr16(InsData[ins].Size)
}

// Run is the main loop of the CPU.
func (c *CPU) Run(cmd chan int) {
	c.IsHalted = false  // User freezes execution, resume with 'resume'
	c.IsStopped = false // STP instruction
	c.IsWaiting = false // WAI instruction

	for {
		// We first check if we have received a command from the user
		// via the command channel from the CLI. Otherwise, execute an
		// insruction. We do not check if we got a correct signal from
		// the CLI, that must be checked at that level.
		select {
		case order := <-cmd:

			switch order {

			case common.ABORT:
				fmt.Println("CPU: DUMMY: Received cmd ABORT")

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

			case common.RESET: // Also used for cold boot
				c.reset()
				c.IsHalted = false
				c.SingleStepMode = false

			case common.RESUME:
				fmt.Println("CPU: DUMMY: Received cmd RESUME")
				c.IsHalted = false

			case common.RUN:
				fmt.Println("CPU: DUMMY: Received cmd RUN")
				c.IsHalted = false
				c.SingleStepMode = false

			case common.STEP:
				fmt.Println("CPU: DUMMY: Received cmd STEP")
				c.SingleStepMode = true

			case common.TRACE:
				fmt.Println("CPU: DUMMY: Received cmd TRACE")
				trace = true

			case common.VERBOSE:
				fmt.Println("CPU: DUMMY: Received cmd VERBOSE")
				verbose = true

				// No default clause because we have the CLI check the
				// signals that we send
			}

		default:
			// This is where the CPU actually runs an
			// instruction.

			if !c.IsHalted && !c.IsStopped {
				c.Step()
				time.Sleep(1 * time.Second) // TODO for testing

				if c.SingleStepMode {
					<-cmd
				}

			} else {
				lock := <-cmd

				if lock == common.RESUME {
					c.IsHalted = false
				}

			}

		}
	}

}

// Reset the machine. This is handled by the RESET signal and is also how we
// cold boot the machine after INIT. Details on the reset procedure are on page
// 201 of Eyes & Lichty.
func (c *CPU) reset() {
	var ok bool

	// Set Direct Page to 0000
	c.DP = 0

	// Set Program Bank Register to 00
	c.PBR = 0

	// Set Data Bank Register to 00
	c.DBR = 0

	// Set Stack high byte to 01
	c.SP = 0x0100

	// TODO Set X Register high to 00 (through x Flag = 1)
	// TODO Set Y Register high to 00 (through x Flag = 1)
	// TODO Status Register: m=1, x=1, d=0, i=1
	// TODO Emulation Flag: 1

	// Get address at 0xFFFC (Reset Vector)
	rv, ok := c.Mem.FetchMore(resetAddr, 2)
	if !ok {
		log.Println("ERROR: Couldn't get RESET vector from", resetAddr)
		return
	}

	c.PC = common.Addr16(rv)

	// TODO Make sure we have "magic number" at address

}
