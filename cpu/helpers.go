// Helper functions for opcodes in Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 03. Feb 2019 (Superbowl LIII)
// This version: 24. Feb 2019 (Oscars)

// This package contains helper functions for the opcodes in opcodes.go

package cpu

import (
	"fmt"

	"angua/common"
)

// --- Load routines ---

// --- Store routines ---

/*
   storeA, storeX, and storeY are variants on the same theme and could be
   combined to one routine, passing the register involved as the parameter.
   For the moment, we leave them separate until we are sure everything works.
*/

// storeA takes a 24-bit address and stores the A register there, either as one
// byte (if A is 8 bit) or two bytes in little-endian (if A is 16 bit). An error
// is returned.
func (c *CPU) storeA(addr common.Addr24) error {
	var err error

	switch c.WidthA {
	case W8:
		err = c.Mem.Store(addr, byte(c.A8))
		if err != nil {
			return fmt.Errorf("storeA: couldn't store A8: %v", err)
		}

	case W16:
		err = c.Mem.StoreMore(addr, uint(c.A16), 2)
		if err != nil {
			return fmt.Errorf("storeA: couldn't store A16: %v", err)
		}

	default: // paranoid
		return fmt.Errorf("storeA: illegal width for register A:%d", c.WidthA)
	}

	return nil
}

// storeX takes a 24-bit address and stores the X register there, either as one
// byte (if X is 8 bit) or two bytes in little-endian (if X is 16 bit). An error
// is returned.
func (c *CPU) storeX(addr common.Addr24) error {
	var err error

	switch c.WidthXY {
	case W8:
		err = c.Mem.Store(addr, byte(c.X8))
		if err != nil {
			return fmt.Errorf("storeX: couldn't store X8: %v", err)
		}

	case W16:
		err = c.Mem.StoreMore(addr, uint(c.X16), 2)
		if err != nil {
			return fmt.Errorf("storeX: couldn't store X16: %v", err)
		}

	default: // paranoid
		return fmt.Errorf("storeX: illegal width for register X:%d", c.WidthXY)
	}

	return nil
}

// storeY takes a 24-bit address and stores the Y register there, either as one
// byte (if Y is 8 bit) or two bytes in little-endian (if Y is 16 bit). An error
// is returned.
func (c *CPU) storeY(addr common.Addr24) error {
	var err error

	switch c.WidthXY {
	case W8:
		err = c.Mem.Store(addr, byte(c.Y8))
		if err != nil {
			return fmt.Errorf("storeY: couldn't store Y8: %v", err)
		}

	case W16:
		err = c.Mem.StoreMore(addr, uint(c.Y16), 2)
		if err != nil {
			return fmt.Errorf("storeY: couldn't store Y16: %v", err)
		}

	default: // paranoid
		return fmt.Errorf("storeY: illegal width for register Y:%d", c.WidthXY)
	}

	return nil
}

// --- Stack routines ---

// pushByte pushes a byte defined to the stack as defined by the stack pointer,
// which it then adusts. This internal routine is used by all other stack push
// instructions such as pushData8 and pushData16.
func (c *CPU) pushByte(b byte) error {
	addr := common.Addr24(c.SP) // c.SP is defined as common.Addr16

	err := c.Mem.Store(addr, b)
	if err != nil {
		return fmt.Errorf("pushByte: couldn't push byte %X to stack at %s: %v",
			b, addr.HexString(), err)
	}

	// Since we don't support emulation mode, we don't have to care about
	// the weird wrapping behavior, see p. 278
	c.SP--

	return nil
}

// pushData8 is a wrapper function for pushByte that takes a common.Data8
// parameter as defined by our registers
func (c *CPU) pushData8(d common.Data8) error {
	b := byte(d)
	err := c.pushByte(b)

	if err != nil {
		return fmt.Errorf("pushData8: couldn't push %X to stack: %v", d.HexString(), err)
	}

	return nil
}

// pushData16 is a wrapper function for pushByte that takes a common.Data16
// parameter as defined by our registers. Remember the MSB is pushed first
func (c *CPU) pushData16(d common.Data16) error {
	msb := d.Msb()

	err := c.pushByte(msb)
	if err != nil {
		return fmt.Errorf("pushData16: couldn't push %X to stack: %v", msb, err)
	}

	lsb := d.Lsb()

	err = c.pushByte(lsb)
	if err != nil {
		return fmt.Errorf("pushData16: couldn't push %X to stack: %v", lsb, err)
	}

	return nil
}

// pullByte is the basic function for pulling a byte of the stack and then
// imcrementing the sack pointer
func (c *CPU) pullByte() (byte, error) {

	// We need to increment the stack pointer first
	c.SP++

	addr := common.Addr24(c.SP) // c.SP is defined as common.Addr16
	b, err := c.Mem.Fetch(addr)
	if err != nil {
		return 0, fmt.Errorf("pullByte: couldn't get byte from stack: %v", err)
	}

	return b, err
}

// pullData8 is a wrapper function to get a byte off the stack and return it as
// a common.Data8 that registers use
func (c *CPU) pullData8() (common.Data8, error) {

	b, err := c.pullByte()
	if err != nil {
		return 0, fmt.Errorf("pullData8: couldn't get byte from stack: %v", err)
	}

	return common.Data8(b), nil
}

// pullData16 is a wrapper function to get a word off the stack and return it as
// a common.Data16 that registers use
func (c *CPU) pullData16() (common.Data16, error) {

	// LSB is pulled first
	lsb, err := c.pullByte()
	if err != nil {
		return 0, fmt.Errorf("pullData16: couldn't get LSB from stack: %v", err)
	}

	// MSB is next
	msb, err := c.pullByte()
	if err != nil {
		return 0, fmt.Errorf("pullData16: couldn't get MSB from stack: %v", err)
	}

	d := (common.Data16(msb) << 8) | common.Data16(lsb)

	return d, nil
}

// --- Helper functions for various instruction groups

// Helper functions for txa (0x8A)

func (c *CPU) txaX8A8() {
	c.A8 = c.X8
	c.TestNZ8(c.A8)
}

func (c *CPU) txaX8A16() {
	c.A16 = common.Data16(c.X8)
	c.TestNZ16(c.A16)
}

func (c *CPU) txaX16A8() {
	c.A8 = common.Data8(c.X16.Lsb())
	c.TestNZ8(c.A8)
}

func (c *CPU) txaX16A16() {
	c.A16 = c.X16
	c.TestNZ16(c.A16)
}

var txaFNS [2][2]func(c *CPU)

func init() {
	txaFNS[W8][W8] = (*CPU).txaX8A8
	txaFNS[W16][W8] = (*CPU).txaX8A16
	txaFNS[W8][W16] = (*CPU).txaX16A8
	txaFNS[W16][W16] = (*CPU).txaX16A16
}

// Helper functions for TYA

func (c *CPU) tyaY8A8() {
	c.A8 = c.Y8
	c.TestNZ8(c.A8)
}

func (c *CPU) tyaY8A16() {
	c.A16 = common.Data16(c.Y8)
	c.TestNZ16(c.A16)
}

func (c *CPU) tyaY16A8() {
	c.A8 = common.Data8(c.Y16.Lsb())
	c.TestNZ8(c.A8)
}

func (c *CPU) tyaY16A16() {
	c.A16 = c.Y16
	c.TestNZ16(c.A16)
}

var tyaFNS [2][2]func(c *CPU)

func init() {
	tyaFNS[W8][W8] = (*CPU).tyaY8A8
	tyaFNS[W16][W8] = (*CPU).tyaY8A16
	tyaFNS[W8][W16] = (*CPU).tyaY16A8
	tyaFNS[W16][W16] = (*CPU).tyaY16A16
}

// Helper functions for tay (0xA8)

func (c *CPU) tayA8Y8() {
	c.Y8 = c.A8
	c.TestNZ8(c.Y8)
}

func (c *CPU) tayA8Y16() {
	MSB := common.Data16(c.B) << 8
	c.Y16 = MSB | common.Data16(c.A8)
	c.TestNZ16(c.Y16)
}

func (c *CPU) tayA16Y8() {
	c.Y8 = common.Data8(c.A16.Lsb())
	c.TestNZ8(c.Y8)
}

func (c *CPU) tayA16Y16() {
	c.Y16 = c.A16
	c.TestNZ16(c.Y16)
}

var tayFNS [2][2]func(c *CPU)

func init() {
	tayFNS[W8][W8] = (*CPU).tayA8Y8
	tayFNS[W8][W16] = (*CPU).tayA8Y16
	tayFNS[W16][W8] = (*CPU).tayA16Y8
	tayFNS[W16][W16] = (*CPU).tayA16Y16
}

// Helper functions for tax (0xAA)

func (c *CPU) taxA8X8() {
	c.X8 = c.A8
	c.TestNZ8(c.X8)
}

func (c *CPU) taxA8X16() {
	MSB := common.Data16(c.B) << 8
	c.X16 = MSB | common.Data16(c.A8)
	c.TestNZ16(c.X16)
}

func (c *CPU) taxA16X8() {
	c.X8 = common.Data8(c.A16.Lsb())
	c.TestNZ8(c.X8)
}

func (c *CPU) taxA16X16() {
	c.X16 = c.A16
	c.TestNZ16(c.X16)
}

var taxFNS [2][2]func(c *CPU)

func init() {
	taxFNS[W8][W8] = (*CPU).taxA8X8
	taxFNS[W8][W16] = (*CPU).taxA8X16
	taxFNS[W16][W8] = (*CPU).taxA16X8
	taxFNS[W16][W16] = (*CPU).taxA16X16
}
