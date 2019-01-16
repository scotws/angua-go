// Interrupt Handling for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 15. Jan 2019
// This version: 16. Jan 2019

// This package contains the hardware interrupt handling routines for Angua, the
// software interupts (BRK and COP) are handled in their own routines.
// Background information for this section:
// - http://6502.org/tutorials/65c816interrupts.html
// - Eyes & Lichty pp. 192

package cpu

import (
	"fmt"

	"angua/common"
)

// commonIntr handles the steps of the interrupt routines which they all share
func (c *CPU) commonIntr() {
	c.pushByte(byte(c.PBR))           // Push PBR to the stack
	c.PBR = 0                         // Load PBR to bank zero
	c.pushData16(common.Data16(c.PC)) // MSB gets pushed first
	b := c.GetStatReg()               // Push status register
	c.pushByte(b)
	c.FlagI = SET   // Prevent other interrupts
	c.FlagD = CLEAR // We go through the motions
}

// The Abort single jumps to the routine at 00:FFE8
// Note that Abort on Angua currently doesn't follow the manual in that the
// current instruction is finished, but no results are stored
// See http://6502.org/tutorials/65c816interrupts.html#toc:interrupt_abort
func (c *CPU) Abort() error {

	// TODO TESTING
	fmt.Println("CPU: DUMMY: Abort interrupt routine")

	c.commonIntr()

	// Load PC with contents of Abort vector from 00:FFE8
	rv, err := c.Mem.FetchMore(common.AbortAddr, 2)
	if err != nil {
		return fmt.Errorf("Reset: couldn't get Abort vector from %s", string(common.AbortAddr))
	}

	c.PC = common.Addr16(rv)

	return nil
}

// The IRQ (maskable interrupt request) jumps to the routine at 00:FFEE
// See http://6502.org/tutorials/65c816interrupts.html#toc:interrupt_irq
func (c *CPU) IRQ() error {

	// TODO TESTING
	fmt.Println("CPU: DUMMY: IRQ interrupt routine")

	if c.FlagI == SET {
		return nil
	}

	c.commonIntr()

	// Load PC with contents of IRQ vector from 00:FFEE
	rv, err := c.Mem.FetchMore(common.IRQAddr, 2)
	if err != nil {
		return fmt.Errorf("Reset: couldn't get IRQ vector from %s", string(common.IRQAddr))
	}

	c.PC = common.Addr16(rv)

	return nil
}

// The NMI (non-maskable interrupt) jumps to the routine at 00:FFEA
// See http://6502.org/tutorials/65c816interrupts.html#toc:interrupt_nmi
func (c *CPU) NMI() error {

	// TODO TESTING
	fmt.Println("CPU: DUMMY: NMI interrupt routine")

	c.commonIntr()

	// Load PC with contents of NMI vector from 00:FFEA
	rv, err := c.Mem.FetchMore(common.NMIAddr, 2)
	if err != nil {
		return fmt.Errorf("Reset: couldn't get NMI vector from %s", string(common.NMIAddr))
	}

	c.PC = common.Addr16(rv)

	return nil
}

// The Reset signal jumps to the routine at 00:FFFC
// See http://6502.org/tutorials/65c816interrupts.html#toc:interrupt_reset
func (c *CPU) Reset() error {

	// For future reference: If the internal clock was stopped by STP or
	// WAI, it will be restarted here.

	// Set Direct Page to 00 (where the Zero Page is on the 6502)
	c.DP = 00

	// Set Program and Data Bank Registers to 00
	c.PBR = 00
	c.DBR = 00

	// The LSB of the SP is in an undefined state after a reset, while the
	// MSB is set to 01. We simulate this to avoid users getting used to a
	// certain value
	c.SP = 0x0100 | common.Addr16(common.UndefinedByte())

	// Clear registers
	// TODO see if this is what really happens. If we have garbage in
	// the registers after a reset, we want to emulate that as well
	c.A8 = 00
	c.A16 = 0000
	c.B = 00
	c.X8 = 00
	c.X16 = 0000
	c.Y8 = 00
	c.Y16 = 0000

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
		return fmt.Errorf("Reset: couldn't get Reset vector from %s", string(common.ResetAddr))
	}

	addr := common.Addr16(rv)

	// Make sure we have "magic number" at address: The first instructions
	// must always be CLC XCE which translates as 0xFB18 because of the
	// little-endian fetch
	bootInst, err := c.Mem.FetchMore(common.Addr24(rv), 2)
	if err != nil {
		return fmt.Errorf("Reset: couldn't get instructions from Reset target %s", addr.HexString())
	}

	if common.Data16(bootInst) != MAGICNUMBER {
		return fmt.Errorf("Reset: code after reset must start with 0xFB18 (clc xce), got: %s", addr.HexString())
	}

	// Everything is fine, we are good to go
	c.PC = addr

	c.IsHalted = false
	c.IsStopped = false // STP
	c.IsWaiting = false // WAI

	return nil
}
