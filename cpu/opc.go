// Opcodes for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 02. Jan 2019
// This version: 10. Jan 2019

// This package contains the opcodes and opcode data for the 65816 instructions.
// Note there is redundancy with the information in the info package. We keep
// these separate for the moment to allow more robustness while editing.

package cpu

import (
	"fmt"

	"angua/common"
)

const (
	// Modes
	ABSOLUTE          = 1  // Absolute                  lda $1000
	ABSOLUTE_X        = 2  // Absolute X indexed        lda.x $1000
	ABSOLUTE_Y        = 3  // Absolute Y indexed        lda.y $1000
	ABSOLUTE_IND      = 3  // Absolute indirect         jmp.i $1000
	ABSOLUTE_IND_LONG = 4  // Absolute indirect long    jmp.il $1000
	ABSOLUTE_LONG     = 5  // Absolute long             jmp.l $101000
	ABSOLUTE_LONG_X   = 6  // Absolute long X indexed   jmp.lx $101000
	ACCUMULATOR       = 7  // Accumulator               inc.a
	BLOCK_MOVE        = 8  // Block move		    mvp
	DP                = 9  // Direct page (DP)          lda.d $10
	DP_IND            = 10 // Direct page indirect      lda.di $10
	DP_IND_X          = 11 // DP indirect X indexed     lda.dxi $10
	DP_IND_Y          = 12 // DP indirect Y indexed     lda.diy $10
	DP_IND_LONG       = 13 // DP indirect long          lda.dil $10
	DP_IND_LONG_Y     = 14 // DP indirect long Y index  lda.dily $10
	DP_X              = 15 // Direct page X indexed     lda.dx $10
	DP_Y              = 16 // Direct page Y indexed     ldx.dy $10
	IMMEDIATE         = 17 // Immediate                 lda.# $00
	IMPLIED           = 18 // Implied                   dex
	INDEX_IND         = 19 // Indexed indirect          jmp.xi $1000
	RELATIVE          = 20 // PC Relative               bra <LABEL>
	RELATIVE_LONG     = 21 // PC Relative long          bra.l <LABEL>
	STACK             = 22 // Stack                     pha
	STACK_REL_IND_Y   = 23 // Stack rel ind Y indexed   lda.siy 3
	STACK_REL         = 24 // Stack relative            lda.s 3
)

type OpcData struct {
	Size    int  // number of bytes including operand (1 to 4)
	Mode    int  // coded instruction modes, see above
	Expands bool // add one byte if register is 16 bit?
}

// In theory, we could use a giant switch statement instead of a map of
// functions. However,
// https://hashrocket.com/blog/posts/switch-vs-map-which-is-the-better-way-to-branch-in-go
// suggests that maps are by far faster
var (
	InsJump [256]func(*CPU) error // Instruction jump table
	InsData [256]OpcData          // Instruction data
)

func init() {
	// Instruction jumps
	InsData[0x00] = OpcData{2, STACK, false}    // brk with signature byte
	InsData[0x01] = OpcData{2, DP_IND_X, false} // ora.dxi
	InsData[0x02] = OpcData{2, STACK, false}    // cop
	// ...
	InsData[0x18] = OpcData{1, IMPLIED, false} // clc
	// ...
	InsData[0x38] = OpcData{1, IMPLIED, false} // sec
	// ...
	InsData[0x58] = OpcData{1, IMPLIED, false} // cli
	// ...
	InsData[0x78] = OpcData{1, IMPLIED, false} // sei
	// ...
	InsData[0x85] = OpcData{2, DP, false} // sta.d
	// ...
	InsData[0x8D] = OpcData{3, ABSOLUTE, false} // sta
	// ...
	InsData[0xA9] = OpcData{2, IMMEDIATE, true} // lda.# (lda.8/lda.16)
	// ...
	InsData[0xAD] = OpcData{3, ABSOLUTE, false} // lda
	// ...
	InsData[0xB8] = OpcData{1, IMPLIED, false} // clv
	// ...
	InsData[0xD8] = OpcData{1, IMPLIED, false} // cld
	// ...
	InsData[0xDB] = OpcData{1, IMPLIED, false} // stp
	// ...
	InsData[0xEA] = OpcData{1, IMPLIED, false} // nop
	// ...
	InsData[0xFB] = OpcData{1, IMPLIED, false} // xce

	// Instruction data
	InsJump[0x00] = Opc00 // brk
	InsJump[0x01] = Opc01 // ora.dxi
	InsJump[0x02] = Opc02 // cop
	// ...
	InsJump[0x18] = Opc18 // clc
	// ...
	InsJump[0x38] = Opc38 // sec
	// ...
	InsJump[0x58] = Opc58 // cli
	// ...
	InsJump[0x78] = Opc78 // sei
	// ...
	InsJump[0xA9] = OpcA9 // lda.# (lda.8/lda.16)
	// ...
	InsJump[0x85] = Opc85 // sta.d
	// ...
	InsJump[0x8D] = Opc8D // sta
	// ...
	InsJump[0xAD] = OpcAD // lda
	// ...
	InsJump[0xB8] = OpcB8 // clv
	// ...
	InsJump[0xD8] = OpcD8 // cld
	// ...
	InsJump[0xDB] = OpcDB // stp
	// ...
	InsJump[0xEA] = OpcEA // nop
	// ...
	InsJump[0xFB] = OpcFB // xce
}

// --- Store Routines ---

// storeA takes a 24-bit address and stores the A register there, either as one
// byte (if A is 8 bit) or two bytes in little-endian (if A is 16 bit). An error
// is returned
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

// --- Mode routines ---

// modeAbsolute returns the address stored in the next two bytes after the
// opcode and an error code.
func (c *CPU) modeAbsolute() (common.Addr24, error) {
	operandAddr := c.getFullPC() + 1
	addrUint, err := c.Mem.FetchMore(operandAddr, 2)
	if err != nil {
		return 0, fmt.Errorf("absolute mode: couldn't fetch address from %s: %v", common.Addr24(addrUint).HexString(), err)
	}

	return common.Addr24(addrUint), nil
}

// modeDirectPage returns the address stored on the Direct Page with the LSB as
// given in the byte after the opcode
func (c *CPU) modeDirectPage() (common.Addr24, error) {
	operandAddr := c.getFullPC() + 1
	dpOffset, err := c.Mem.Fetch(operandAddr)
	if err != nil {
		return 0, fmt.Errorf("direct page mode: couldn't fetch address from %s: %v", common.Addr24(operandAddr).HexString(), err)
	}

	addr := common.Addr24(c.DP) + common.Addr24(dpOffset)

	return addr, nil
}

// --- Instruction Functions ---

// The opcodes get the next bytes but do not change the PC

func Opc00(c *CPU) error { // brk
	fmt.Println("OPC: DUMMY: Executing brk (00)")
	return nil
}

func Opc01(c *CPU) error { // ora.dxi
	fmt.Println("OPC: DUMMY: Executing ora.dxi (02)")
	return nil
}

func Opc02(c *CPU) error { // cop
	fmt.Println("OPC: DUMMY: Executing cop (03)")
	return nil
}

// ...

func Opc18(c *CPU) error { // clc
	c.FlagC = CLEAR
	return nil
}

// ...

func Opc38(c *CPU) error { // sec
	c.FlagC = SET
	return nil
}

// ...

func Opc58(c *CPU) error { // cli
	c.FlagI = CLEAR
	return nil
}

// ...

func Opc78(c *CPU) error { // sei
	c.FlagI = SET
	return nil
}

// ...

func Opc85(c *CPU) error { // sta.d

	addr, err := c.modeDirectPage()
	if err != nil {
		return fmt.Errorf("sta.d (85): couldn't fetch address from %s: %v", addr.HexString(), err)
	}

	err = c.storeA(addr)
	if err != nil {
		return fmt.Errorf("sta.d (85): couldn't store A at address %s: %v", addr.HexString(), err)
	}

	return nil
}

// ...

func Opc8D(c *CPU) error { // sta

	addr, err := c.modeAbsolute()
	if err != nil {
		return fmt.Errorf("sta (8D): couldn't fetch address from %s: %v", addr.HexString(), err)
	}

	err = c.storeA(addr)
	if err != nil {
		return fmt.Errorf("sta (8D): couldn't store A at address %s: %v", addr.HexString(), err)
	}

	return nil
}

// ...

func OpcA9(c *CPU) error { // lda.# (lda.8/lda.16)
	fmt.Println("OPC: DUMMY: Executing lda.# (A9) ")
	return nil
}

// ...

func OpcAD(c *CPU) error { // lda

	addr, err := c.modeAbsolute()
	if err != nil {
		return fmt.Errorf("lda (AD): Couldn't fetch address from %s: %v", addr.HexString(), err)
	}

	// TODO generalize this in a routine for all LDA
	switch c.WidthA {
	case W8:
		b, err := c.Mem.Fetch(addr)
		if err != nil {
			return fmt.Errorf("lda (0xAD): Couldn't fetch byte from %s: %v", addr.HexString(), err)
		}

		c.A8 = common.Data8(b)
		c.TestNZ8(c.A8)

	case W16:
		b, err := c.Mem.FetchMore(addr, 29)
		if err != nil {
			return fmt.Errorf("lda (0xAD): Couldn't fetch two bytes from %s: %v", addr.HexString(), err)
		}

		c.A16 = common.Data16(b)
		c.TestNZ16(c.A16)

	default: // paranoid
		return fmt.Errorf("lda (0xAD): Illegal width for register A:%d", c.WidthA)
	}

	return nil
}

// ...

func OpcB8(c *CPU) error { // clv
	c.FlagV = CLEAR
	return nil
}

// ...

func OpcD8(c *CPU) error { // cld
	c.FlagD = CLEAR
	return nil
}

// ...

func OpcDB(c *CPU) error { // stp
	c.IsStopped = true
	// TODO print Addr24
	fmt.Println("Machine stopped by STP (0xDB) in block", c.PBR.HexString(), "at address", c.PC.HexString())
	return nil
}

// ...

func OpcEA(c *CPU) error { // nop
	// We return the execution of a 'nop' as an error and let the higher-ups
	// decide what to do with it
	return fmt.Errorf("OpcEA: Executed 'nop' (0xEA) at %s:%s", c.PBR.HexString(), c.PC.HexString())
}

// ...

func OpcFB(c *CPU) error { // xce
	tmp := c.FlagE
	c.FlagE = c.FlagC
	c.FlagC = tmp
	return nil
}
