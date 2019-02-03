// Mode routines for the Angua opcodes
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 03. Feb 2019 (Superbowl LIII)
// This version: 03. Feb 2019 (Superbowl LIII)

// This package contains the functions for the various modes of the angua
// opcodes

package cpu

import (
	"fmt"

	"angua/common"
)

// modeAbsolute returns the address stored in the next two bytes after the
// opcode and an error code.
func (c *CPU) modeAbsolute() (common.Addr24, error) {
	operandAddr := c.getFullPC() + 1
	addrUint, err := c.Mem.FetchMore(operandAddr, 2)
	if err != nil {
		return 0,
			fmt.Errorf("absolute mode: couldn't fetch address from %s: %v",
				common.Addr24(operandAddr).HexString(), err)
	}

	return common.Addr24(addrUint), nil
}

// modeBranch takes a byte as a signed int and returns the address created by an
// offset to the PC. Note the return address is common.Addr16, not
// common.Addr24
func (c *CPU) modeBranch(b byte) (common.Addr16, error) {

	// Convert byte offset to int8 first to preserve the sign
	offset := int8(b)
	addr := int(c.PC)

	// Now we need to calculate it all in int
	newAddr := common.Addr16(addr+int(offset)) + 2

	if !c.Mem.Contains(common.Addr24(newAddr)) {
		return 0, fmt.Errorf("modeBranch: address %s illegal", newAddr.HexString())
	}

	return newAddr, nil
}

// modeDirectPage returns the address stored on the Direct Page with the LSB as
// given in the byte after the opcode
func (c *CPU) modeDirectPage() (common.Addr24, error) {
	operandAddr := c.getFullPC() + 1
	dpOffset, err := c.Mem.Fetch(operandAddr)
	if err != nil {
		return 0,
			fmt.Errorf("direct page mode: couldn't fetch address from %s: %v",
				common.Addr24(operandAddr).HexString(), err)
	}

	addr := common.Addr24(c.DP) + common.Addr24(dpOffset)

	return addr, nil
}

// modeImmediate8 returns the byte stored in the address after the opcode and an
// error. This is a variant of getNextByte, except that we return a common.Data8
// instead of a byte. Keep these routines separate to allow modifications.
func (c *CPU) modeImmediate8() (common.Data8, error) {
	operandAddr := c.getFullPC() + 1
	operand, err := c.Mem.Fetch(operandAddr)
	if err != nil {
		return 0,
			fmt.Errorf("immediate 8 mode: couldn't fetch data from %s: %v",
				common.Addr24(operandAddr).HexString(), err)
	}

	return common.Data8(operand), nil
}

// getNextData8 is a synonym for modeImmediate8
func (c *CPU) getNextData8() (common.Data8, error) {
	return c.modeImmediate8()
}

// modeImmediate16 returns the word stored in the address after the opcode and an
// error.
func (c *CPU) modeImmediate16() (common.Data16, error) {
	operandAddr := c.getFullPC() + 1
	ui, err := c.Mem.FetchMore(operandAddr, 2)
	if err != nil {
		return 0,
			fmt.Errorf("immediate 16 mode: couldn't fetch data from %s: %v",
				common.Addr24(operandAddr).HexString(), err)
	}

	return common.Data16(ui), nil
}

// modeLong returns the address stored in the next three bytes after the
// opcode and an error code.
func (c *CPU) modeLong() (common.Addr24, error) {
	operandAddr := c.getFullPC() + 1
	addrUint, err := c.Mem.FetchMore(operandAddr, 3)
	if err != nil {
		return 0,
			fmt.Errorf("long mode: couldn't fetch address from %s: %v",
				common.Addr24(operandAddr).HexString(), err)
	}

	return common.Addr24(addrUint), nil
}

// getNextData16 is a synonym for modeImmediate16
func (c *CPU) getNextData16() (common.Data16, error) {
	return c.modeImmediate16()
}

// getNextByte takes a pointer to the CPU and returns the next byte - usually
// the byte after the opcode - and an error message. This is a slight variation
// in modeImmediate8, except we return a byte and not common.Data8. Keep them
// separate so we can modify them if required
func (c *CPU) getNextByte() (byte, error) {
	byteAddr := c.getFullPC() + 1
	b, err := c.Mem.Fetch(byteAddr)
	if err != nil {
		return 0,
			fmt.Errorf("getNextByte: couldn't fetch data from %s: %v",
				common.Addr24(byteAddr).HexString(), err)
	}

	return b, nil
}
