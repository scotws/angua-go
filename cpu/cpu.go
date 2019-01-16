// Angua CPU System
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 06. Nov 2018
// This version: 16. Jan 2019

package cpu

import (
	"fmt"
	"math/bits"

	"angua/common"
	"angua/mem"
)

const (
	// Width of accumulator and registers. We follow the flag convention for
	// M and X flags: 0 (clear) is 16 bits, 1 (set) 8 bit. If these are
	// changed, the function immediateOffset further below must be changed
	// as well
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

// Step executes a single instruction from PC.
func (c *CPU) Step() error {

	// Get byte at PC
	ins, err := c.Mem.Fetch(c.getFullPC())
	if err != nil {
		return fmt.Errorf("Step: can't get instruction at %s", c.getFullPC().HexString())
	}

	// Execute the instruction by accessing the entry in the Instruction
	// Jump table. We pass a pointer to the CPU struct. The instructions are
	// responsible for updating the PC
	err = InsSet[ins].Code(c)
	if err != nil {
		return fmt.Errorf("Step: instruction returned error: %v", err)
	}

	return nil
}

// Run is the main loop of the CPU.
func (c *CPU) Run(cmd chan int) {
	var err error
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
				err = c.Abort()
				if err != nil {
					fmt.Printf("Abort returned error: %v", err)
				}

			case common.DESTROY:
				// Stop the go routine
				return

			case common.HALT:
				c.IsHalted = true

			case common.IRQ:
				err = c.IRQ()
				if err != nil {
					fmt.Printf("IRQ returned error: %v", err)
				}

			case common.NOTRACE:
				fmt.Println("CPU: DUMMY: Received cmd NOTRACE")
				trace = false

			case common.NOVERBOSE:
				fmt.Println("CPU: DUMMY: Received cmd NOVERBOSE")
				verbose = false

			case common.NMI:
				err = c.NMI()
				if err != nil {
					fmt.Printf("NMI returned error: %v", err)
				}

			case common.RESET: // Also used for cold boot
				err = c.Reset()
				if err != nil {
					fmt.Printf("CPU: Reset returned error: %v", err)
				}

				c.IsHalted = false
				c.SingleStepMode = false

			case common.RESUME:
				c.IsHalted = false

			case common.STEP:
				fmt.Println("CPU: DUMMY: Received cmd STEP")
				c.SingleStepMode = true

			case common.TRACE:
				trace = true

			case common.VERBOSE:
				verbose = true

				// No default clause because we have the CLI check the
				// signals that we send
			}

		default:

			// This is where the CPU actually runs an instruction.
			if !c.IsHalted && !c.IsStopped {
				err = c.Step()
				if err != nil {
					fmt.Printf("CPU: Execution error: %v", err)
					c.IsHalted = true
				}

				if c.SingleStepMode {
					<-cmd
				}

			} else {
				lock := <-cmd

				switch lock {

				case common.RESET:
					c.Reset()
				case common.RESUME:
					c.IsHalted = false

				}

			}

		}
	}

}
