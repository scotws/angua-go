// Opcodes for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 02. Jan 2019
// This version: 11. Jan 2019

// This package contains the opcodes and opcode data for the 65816 instructions.
// Note there is redundancy with the information in the info package. We keep
// these separate for the moment to allow more robustness while editing.

package cpu

import (
	"fmt"

	"angua/common"
)

// TODO see if we even need these modes
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
	InsData[0xC2] = OpcData{2, IMMEDIATE, false} // rep
	// ...
	InsData[0xD8] = OpcData{1, IMPLIED, false} // cld
	// ...
	InsData[0xDB] = OpcData{1, IMPLIED, false} // stp
	// ...
	InsData[0xE2] = OpcData{2, IMPLIED, false} // sep
	// ...
	InsData[0xE8] = OpcData{1, IMPLIED, true} // inx
	// ...
	InsData[0xEA] = OpcData{1, IMPLIED, false} // nop
	InsData[0xEB] = OpcData{1, IMPLIED, false} // xba
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
	InsJump[0xC2] = OpcC2 // rep
	// ...
	InsJump[0xD8] = OpcD8 // cld
	// ...
	InsJump[0xDB] = OpcDB // stp
	// ...
	InsJump[0xE2] = OpcE2 // sep
	// ...
	InsJump[0xE8] = OpcE8 // inx
	// ...
	InsJump[0xEA] = OpcEA // nop
	InsJump[0xEB] = OpcEB // xba
	// ...
	InsJump[0xFB] = OpcFB // xce
}

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

// modeImmediate8 returns the byte stored in the address after the opcode and an
// error. This is a variant of getNextByte, except that we return a common.Data8
// instead of a byte. Keep these routines separate to allow modifications. This
// could also be named getNextData8()
func (c *CPU) modeImmediate8() (common.Data8, error) {
	operandAddr := c.getFullPC() + 1
	operand, err := c.Mem.Fetch(operandAddr)
	if err != nil {
		return 0, fmt.Errorf("immediate 8 mode: couldn't fetch data from %s: %v", common.Addr24(operandAddr).HexString(), err)
	}

	return common.Data8(operand), nil
}

// modeImmediate16 returns the byte stored in the address after the opcode and an
// error
func (c *CPU) modeImmediate16() (common.Data16, error) {
	operandAddr := c.getFullPC() + 1
	operand, err := c.Mem.Fetch(operandAddr)
	if err != nil {
		return 0, fmt.Errorf("immediate 16 mode: couldn't fetch data from %s: %v", common.Addr24(operandAddr).HexString(), err)
	}

	return common.Data16(operand), nil
}

// --- Low-level helper functions ---

// getNextByte takes a pointer to the CPU and returns the next byte - usually
// the byte after the opcode - and an error message. This is a slight variation
// in modeImmediate8, except we return a byte and not common.Data8. Keep them
// separate so we can modify them if required
func getNextByte(c *CPU) (byte, error) {
	byteAddr := c.getFullPC() + 1
	b, err := c.Mem.Fetch(byteAddr)
	if err != nil {
		return 0, fmt.Errorf("getNextByte: couldn't fetch data from %s: %v", common.Addr24(byteAddr).HexString(), err)
	}

	return b, nil
}

// --- Instruction functions ---

// The opcodes get the next bytes but do not change the PC, this is left for the
// main loop

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

// ---- AAAA ----

func OpcA9(c *CPU) error { // lda.# (lda.8/lda.16)

	switch c.WidthA {
	case W8:
		operand, err := c.modeImmediate8()
		if err != nil {
			return fmt.Errorf("lda.8 (A9): Couldn't fetch data: %v", err)
		}

		c.A8 = operand
		c.TestNZ8(c.A8)

	case W16:
		operand, err := c.modeImmediate16()
		if err != nil {
			return fmt.Errorf("lda.16 (A9): Couldn't fetch data: %v", err)
		}

		c.A16 = operand
		c.TestNZ16(c.A16)

	default: // paranoid
		return fmt.Errorf("lda.# (0xA9): Illegal width for register A:%d", c.WidthA)
	}

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
		b, err := c.Mem.FetchMore(addr, 2)
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

// ---- BBBB ----

func OpcB8(c *CPU) error { // clv
	c.FlagV = CLEAR
	return nil
}

// ...

// ---- CCCC ----

func OpcC2(c *CPU) error { // rep
	rb, err := getNextByte(c)
	if err != nil {
		return fmt.Errorf("rep (0xC2): can't get next byte: %v:", err)
	}

	// The sequence is NVMXDIZE. All bits which are set turn the equivalent
	// flag to zero. Most imporant are the M X flags because they trigger
	// the changes to the width of the A and XY registers.
	oldFlagM := c.FlagM
	oldFlagX := c.FlagX

	sb := c.GetStatReg()
	nb := ^rb & sb // complement (inverse) is ^ not ~ in Go
	c.SetStatReg(nb)

	// Now that we've gotten the formal part out of the way, we need to see
	// if we've reset the M or X flags. We need to do something if the flag
	// was SET before and CLEAR now

	// Change A size from 8 bit to 16 bit, see p. 51
	if (oldFlagM == SET) && (c.FlagM == CLEAR) {
		c.WidthA = W16
		c.A16 = (common.Data16(c.B) << 8) + common.Data16(c.A8)
	}

	// Change XY size from 8 bit to 16 bit, see p. 51
	if (oldFlagX == SET) && (c.FlagX == CLEAR) {
		c.WidthXY = W16
		c.X16 = 0x000 + common.Data16(c.X8) // emphasise: MSB is zero
		c.Y16 = 0x000 + common.Data16(c.Y8) // emphasise: MSB is zero
	}

	return nil
}

// ...

// ---- DDDD ----

func OpcD8(c *CPU) error { // cld
	c.FlagD = CLEAR
	return nil
}

// ...

func OpcDB(c *CPU) error { // stp
	// TODO print Addr24
	c.IsStopped = true
	fmt.Println("Machine stopped by STP (0xDB) in block", c.PBR.HexString(), "at address", c.PC.HexString())
	return nil
}

// ...

// ---- EEEE ----

func OpcE2(c *CPU) error { // sep

	rb, err := getNextByte(c)
	if err != nil {
		return fmt.Errorf("sep (0xE2): can't get next byte: %v:", err)
	}

	// The sequence is NVMXDIZE. All bits which are set also set the
	// corresponding flag. Most imporant are the M X flags because they
	// trigger the changes to the width of the A and XY registers.
	oldFlagM := c.FlagM
	oldFlagX := c.FlagX

	sb := c.GetStatReg()
	nb := rb | sb
	c.SetStatReg(nb)

	// Now that we've gotten the formal part out of the way, we need to see
	// if we've reset the M or X flags. We need to do something if the flag
	// was SET before and CLEAR now

	// Change A size from 16 bit to 8 bit, see p. 51
	if (oldFlagM == CLEAR) && (c.FlagM == SET) {
		c.WidthA = W8
		c.A8 = common.Data8(c.A16 & 0x00FF)
		c.B = common.Data8((c.A16 >> 8) & 0x00FF)
	}

	// Change XY size from 16 bit to 8 bit, see p. 51
	if (oldFlagX == CLEAR) && (c.FlagX == SET) {
		c.WidthXY = W8
		c.X8 = common.Data8(c.X16 & 0x00FF)
		c.Y8 = common.Data8(c.Y16 & 0x00FF)
	}

	return nil
}

func OpcE8(c *CPU) error { // inx
	switch c.WidthXY {

	case W8:
		c.X8 += 1
		c.TestNZ8(c.X8)

	case W16:
		c.X16 += 1
		c.TestNZ16(c.X16)

	default: // paranoid
		return fmt.Errorf("inx (0xE8): illegal WidthXY value: %v", c.WidthXY)
	}

	return nil
}

func OpcEA(c *CPU) error { // nop
	// We return the execution of a 'nop' as an error and let the higher-ups
	// decide what to do with it
	return fmt.Errorf("OpcEA: executed 'nop' (0xEA) at %s:%s", c.PBR.HexString(), c.PC.HexString())
}

func OpcEB(c *CPU) error { // xba p.422
	switch c.WidthA {

	case W8:
		// Exchange A8 and B
		tmp := c.B
		c.B = c.A8
		c.A8 = tmp
		c.TestNZ8(c.A8)

	case W16:
		// Swap LSB and MSB
		tmp := (c.A16 >> 8) & 0x00FF
		c.A16 = (c.A16 << 8) & 0xFF00
		c.A16 = c.A16 | tmp
		c.TestNZ8(common.Data8(tmp)) // XBA tests the lower byte always

	default: // paranoid
		return fmt.Errorf("xba (0xEB): illegal WidthA value: %v", c.WidthA)
	}

	return nil
}

// ---- FFFF ----

func OpcFB(c *CPU) error { // xce
	tmp := c.FlagE
	c.FlagE = c.FlagC
	c.FlagC = tmp
	return nil
}
