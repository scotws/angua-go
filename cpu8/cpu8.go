// Angua CPU System - Emulated Mode (8-bit) CPU
// Scot W. Stevenson
// First version: 06. Nov 2018
// First version: 06. Nov 2018

package cpu8

import (
	"fmt"
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
	A reg8
	B reg8 // Special hidden register of the 65816
	X reg8
	Y reg8

	DP reg16 // Direct Page register, 16 bit, not 8
	SP reg8  // Stack Pointer, 8 bit

	P byte // Status Register

	DBR reg8 // Data Bank Register, available in emulated mode
	PBR reg8 // Program Bank Register, available in emulated mode

	PC reg16 // Program counter

	StatReg
}

// TODO test version to Execute one opcode
func (c *Cpu8) Execute(b byte) {
	opcodes8[b](c)
}
