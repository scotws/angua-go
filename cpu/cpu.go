// Angua CPU System
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 06. Nov 2018
// This version: 05. Jan 2019

package cpu

import (
	"fmt"
	"log"
	"math/bits"
	"time"

	"angua/common"
	"angua/mem"
)

const (
	// Interrupt vectors. Note the reset vector is only for emulated mode.
	// See http://6502.org/tutorials/65c816interrupts.html and Eyes & Lichty
	// p. 195 for details. We store these as 24 bit addresses because that
	// is the way we'll use them during mem.FetchMore().
	abortAddr common.Addr24 = 0xFFE8
	brkAddr   common.Addr24 = 0xFFE6
	copAddr   common.Addr24 = 0xFFE4
	irqAddr   common.Addr24 = 0xFFEE
	nmiAddr   common.Addr24 = 0xFFEA
	resetAddr common.Addr24 = 0xFFFC // <-- This is the important one

	// Width of accumulator and registers. We follow the flag convention for
	// M and X flags: 0 (clear) is 16 bits, 1 (set) 8 bit
	W8  = 1
	W16 = 0

	// Convenience definitions for working with flags
	SET   = 1
	CLEAR = 0

	// "Magic Number": Because we only emulate native mode, the first
	// two instructions after a reset must be CLC XCE, which in hex is 0x18
	// 0xFB, or loaded as a little-endian 16 bit number 0xFB18. We check for
	// this after a reboot
	MAGICNUMBER common.Data16 = 0xFB18
)

// --------------------------------------------------
// Status Register

/*
	You would think that using bools is the logical way to go with the
	flags, but it turns out that numbers are easier to move the bits around.
	D for decimal mode and E for emulated are provided, but not functional.
	If you are coming from the 6502/65c02, note that the BRK instruction is
	handled differently.
*/
type StatReg struct {
	FlagN byte // Negative flag
	FlagV byte // Overflow flag
	FlagM byte // A size. SET (1) is 8 bit, CLEAR (0) is 16 bit
	FlagX byte // XY size. SET (1) is 8 bit, CLEAR (0) is 16 bit
	FlagD byte // Decimal mode. Ignored with this emulator.
	FlagI byte // Interrupt disable. SET (1) is disabled.
	FlagZ byte // Zero flag. SET (1) means zero was found.
	FlagC byte // Carry flag. SET (1) means we have a carry.

	FlagB byte // Break Instruction, SET (1) if IRQ came from BRK (6502, unused)
	FlagE byte // Emulation flag. SET (1) means emulated, CLEAR (0) native
}

var (
	cmd = make(<-chan int, 2) // Receive commands from CLI

	verbose bool // Print lots of information
	trace   bool // Show us what is going on
)

// GetStatReg creates a status byte out of the flags of the Status Register
// and returns it to the caller. It is used by the instuction PHP for example.
// The flag sequence is NVMXDIZC
func (s *StatReg) GetStatReg() byte {
	var sb byte

	n := s.FlagN << 7
	v := s.FlagV << 6
	m := s.FlagM << 5
	x := s.FlagX << 4
	d := s.FlagD << 3
	i := s.FlagI << 2
	z := s.FlagZ << 1
	c := s.FlagC

	sb = n + v + m + x + d + i + z + c

	return sb
}

// SetStatReg takes a byte and sets the flags of the Status Register
// accordingly. It is used by the instruction PLP for example.
func (s *StatReg) SetStatReg(b byte) {

	s.FlagN = bits.Reverse8(b & 0x80)
	s.FlagV = (b & 0x40) >> 6
	s.FlagM = (b & 0x20) >> 5
	s.FlagX = (b & 0x10) >> 4
	s.FlagD = (b & 0x08) >> 3
	s.FlagI = (b & 0x04) >> 2
	s.FlagZ = (b & 0x02) >> 1
	s.FlagC = (b & 0x01)
}

// StringStatReg returns the status register as an eight rune string with 1
// for set flags and 0 for cleared. The sequence is NVMXDIZC
func (s *StatReg) StringStatReg() string {
	var sb byte

	sb = s.GetStatReg()

	return fmt.Sprintf("%08b", sb)
}

// TestZ takes an int and sets the Z flag to true if the value is zero and to
// false otherwise
func (s *StatReg) TestAndSetZ(i int) {
	if i == 0 {
		s.FlagZ = SET
	} else {
		s.FlagZ = CLEAR
	}
}

// TestN takes a int and sets the N flag if highest bit is set, else clears it
func (s *StatReg) TestAndSetN(i int) {
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

	WidthA  int // Current width of Accumulator, either W8 or W16
	WidthXY int // Current width of X and Y registers, either W8 or W16

	IsHalted       bool // Signals if CPU stopped by Angua CLI
	IsWaiting      bool // CPU stopped by WAI instruction
	IsStopped      bool // CPU stopped by STP instruction
	SingleStepMode bool // Signals if we are in single step mode

	Mem *mem.Memory // Pointer to the memory we're working on

	StatReg
}

// getFullPC merges the Bank Byte and the Program Counter and returns the
// current address as a 24 bit value
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

	// Set Direct Page to 0000 (where the Zero Page is on the 6502)
	c.DP = 0

	// Set Stack high byte to 01 (where it is on the 6502)
	c.SP = 0x0100

	// Set Program Bank Register to 00
	c.PBR = 0

	// Set Data Bank Register to 00
	c.DBR = 0

	// Clear registers
	// TODO see if this is what really happens, if we have garbage in
	// the registers after a reset, we want to emulate that as well
	c.A8 = 0
	c.A16 = 0
	c.B = 0
	c.X8 = 0
	c.X16 = 0
	c.Y8 = 0
	c.Y16 = 0

	// Set register widths the fast way
	c.WidthA = W8
	c.WidthXY = W8

	// Status Register: M=1, X=1, D=0, I=1
	// TODO But what happens to the others?
	c.FlagM = SET   // Width of A is 8 bit
	c.FlagX = SET   // Width of XY is 8 bit
	c.FlagD = CLEAR // "Best practice" though we don't use it
	c.FlagI = SET   // Stop interrupts

	// Reset switches us to emulated mode, which means we'll have to get
	// out of it as soon as possible
	c.FlagE = SET

	// Get address at 0xFFFC (Reset Vector)
	rv, ok := c.Mem.FetchMore(resetAddr, 2)
	if !ok {
		log.Println("ERROR: Couldn't get RESET vector from", resetAddr)
		return
	}

	addr := common.Addr16(rv)

	// Make sure we have "magic number" at address: The first instructions
	// must always be CLC XCE which translates as 0xFB18 because of the
	// little-endian fetch
	bootInst, ok := c.Mem.FetchMore(common.Addr24(rv), 2)
	if !ok {
		log.Println("ERROR: Couldn't get instructions from Reset target", addr.HexString())
		return
	}

	if common.Data16(bootInst) != MAGICNUMBER {
		log.Println("ERROR: Reset address must start with 0xFB18 (CLC XCE), got", addr.HexString())
		return
	}

	// Everything is fine, we are good to go
	c.PC = addr
	c.IsHalted = false
	c.IsStopped = false
	c.IsWaiting = false
}
