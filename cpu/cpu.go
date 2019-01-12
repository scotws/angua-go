// Angua CPU System
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 06. Nov 2018
// This version: 11. Jan 2019

package cpu

import (
	"fmt"
	"math/bits"

	"angua/common"
	"angua/mem"
)

const (
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

// TestNZ8 takes a 8-bit data type such as a register and sets the N or Z flag
func (s *StatReg) TestNZ8(d common.Data8) {
	s.FlagN = byte((d >> 7) & 0x01)

	if d == 0 {
		s.FlagZ = SET
	} else {
		s.FlagZ = CLEAR
	}
}

// TestNZ16 takes a 16-bit data type such as a register and sets the N or Z flag
func (s *StatReg) TestNZ16(d common.Data16) {
	s.FlagN = byte((d >> 15) & 0x0001)

	if d == 0 {
		s.FlagZ = SET
	} else {
		s.FlagZ = CLEAR
	}
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
// TODO this is pretty much all fake at the moment
func (c *CPU) Step() {

	// Get byte at PC
	ins, err := c.Mem.Fetch(c.getFullPC())
	if err != nil {
		fmt.Errorf("Step: can't get instruction at %s", c.getFullPC().HexString())
		return
	}

	// Execute the instruction by accessing the entry in the Instruction
	// Jump table. We pass a pointer to the CPU struct.
	InsSet[ins].Code(c)
	c.PC += common.Addr16(InsSet[ins].Size)
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

				if c.SingleStepMode {
					// TODO print machine status
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
// 201 of Eyes & Lichty; see http://6502.org/tutorials/65c816interrupts.html for
// a more detailed discussion
func (c *CPU) reset() error {

	// For future reference: If the internal clock was stopped by STP or
	// WAI, it will be restarted here.

	// Set Direct Page to 0000 (where the Zero Page is on the 6502)
	c.DP = 0

	// Set Program and Data Bank Registers to 00
	c.PBR = 0
	c.DBR = 0

	// The LSB of the SP is in an undefined state after a reset, while the
	// MSB is set to 01. We simulate this to avoid users getting used to a
	// certain value
	c.SP = 0x0100 | common.Addr16(common.UndefinedByte())

	// Clear registers
	// TODO see if this is what really happens. If we have garbage in
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
	c.FlagM = SET   // Width of A is 8 bit
	c.FlagX = SET   // Width of XY is 8 bit
	c.FlagD = CLEAR // "Best practice" though we don't use it
	c.FlagI = SET   // Stop interrupts

	// Reset switches us to emulated mode, which means we'll have to get
	// out of it as soon as possible
	c.FlagE = SET

	// The flags CNVZ are undefined after a reset. We simulated this to make
	// sure the user doesn't get used to this.
	c.FlagC = common.UndefinedBit()
	c.FlagN = common.UndefinedBit()
	c.FlagV = common.UndefinedBit()
	c.FlagZ = common.UndefinedBit()

	// Get address at 0xFFFC (Reset Vector)
	rv, err := c.Mem.FetchMore(common.ResetAddr, 2)
	if err != nil {
		return fmt.Errorf("Reset: Couldn't get Reset vector from %s", string(common.ResetAddr))
	}

	addr := common.Addr16(rv)

	// Make sure we have "magic number" at address: The first instructions
	// must always be CLC XCE which translates as 0xFB18 because of the
	// little-endian fetch
	bootInst, err := c.Mem.FetchMore(common.Addr24(rv), 2)
	if err != nil {
		return fmt.Errorf("Reset: Couldn't get instructions from Reset target %s", addr.HexString())
	}

	if common.Data16(bootInst) != MAGICNUMBER {
		return fmt.Errorf("Reset: Code after reset must start with 0xFB18 (clc xce), got: %s", addr.HexString())
	}

	// Everything is fine, we are good to go
	c.PC = addr
	c.IsHalted = false
	c.IsStopped = false
	c.IsWaiting = false

	return nil
}
