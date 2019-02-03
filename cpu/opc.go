// Opcodes for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 02. Jan 2019
// This version: 19. Jan 2019

// This package contains the opcodes and opcode data for the 65816 instructions.
// Note there is redundancy with the information in the info package. We keep
// these separate for the moment to allow more robustness while editing.

package cpu

import (
	"fmt"

	"angua/common"
)

// TODO see if we need the boolean
type OpcData struct {
	Size     int              // number of bytes including operand (1 to 4)
	Code     func(*CPU) error // function for actual code of instruction
	Expands  bool             // true if size affected by 8->16 bit register switch
	Mnemonic string           // SAN mnemonic (for the disassembler)
}

// In theory, we could use a giant switch statement instead of a map of
// functions. However,
// https://hashrocket.com/blog/posts/switch-vs-map-which-is-the-better-way-to-branch-in-go
// suggests that maps are faster
var (
	InsSet [256]OpcData
)

func init() {
	InsSet[0x00] = OpcData{2, (*CPU).Opc00, false, "brk"} // with signature byte
	InsSet[0x01] = OpcData{2, (*CPU).Opc01, false, "ora.dxi"}
	InsSet[0x02] = OpcData{2, (*CPU).Opc02, false, "cop"} // with signature byte
	// ...
	InsSet[0x18] = OpcData{1, (*CPU).Opc18, false, "clc"}
	// ...
	InsSet[0x20] = OpcData{3, (*CPU).Opc20, false, "jsr"}
	// ...
	InsSet[0x38] = OpcData{1, (*CPU).Opc38, false, "sec"}
	// ...
	InsSet[0x40] = OpcData{1, (*CPU).Opc40, false, "rti"}
	// ...
	InsSet[0x48] = OpcData{1, (*CPU).Opc48, false, "pha"}
	// ...
	InsSet[0x4B] = OpcData{1, (*CPU).Opc4B, false, "phk"} // pushes PBR
	InsSet[0x4C] = OpcData{3, (*CPU).Opc4C, false, "jmp"}
	// ...
	InsSet[0x58] = OpcData{1, (*CPU).Opc58, false, "cli"}
	// ...
	InsSet[0x60] = OpcData{1, (*CPU).Opc60, false, "rts"}
	// ...
	InsSet[0x68] = OpcData{1, (*CPU).Opc68, false, "pla"}
	// ...
	InsSet[0x78] = OpcData{1, (*CPU).Opc78, false, "sei"}
	// ...
	InsSet[0x7A] = OpcData{1, (*CPU).Opc7A, false, "ply"}
	// ...
	InsSet[0x80] = OpcData{2, (*CPU).Opc80, false, "bra"}
	// ...
	InsSet[0x85] = OpcData{2, (*CPU).Opc85, false, "sta.d"}
	// ...
	InsSet[0x8A] = OpcData{1, (*CPU).Opc8A, false, "txa"}
	InsSet[0x8B] = OpcData{1, (*CPU).Opc8B, false, "phb"} // pushes DBR
	// ...
	InsSet[0x8D] = OpcData{3, (*CPU).Opc8D, false, "sta"}
	// ...
	InsSet[0x9A] = OpcData{1, (*CPU).Opc9A, false, "txs"}
	InsSet[0x9B] = OpcData{1, (*CPU).Opc9B, false, "txy"}
	// ...
	InsSet[0xA8] = OpcData{1, (*CPU).OpcA8, false, "tay"}
	InsSet[0xA9] = OpcData{2, (*CPU).OpcA9, true, "lda.#"} // includes lda.8/lda.16
	// ...
	InsSet[0xAA] = OpcData{1, (*CPU).OpcAA, false, "tax"}
	// ...
	InsSet[0xAD] = OpcData{3, (*CPU).OpcAD, false, "lda"}
	// ...
	InsSet[0xB8] = OpcData{1, (*CPU).OpcB8, false, "clv"}
	// ...
	InsSet[0xBB] = OpcData{1, (*CPU).OpcBB, false, "tyx"}
	// ...
	InsSet[0xC2] = OpcData{2, (*CPU).OpcC2, false, "rep.#"}
	// ...
	InsSet[0xC8] = OpcData{1, (*CPU).OpcC8, false, "iny"}
	// ...
	InsSet[0xD8] = OpcData{1, (*CPU).OpcD8, false, "cld"}
	// ...
	InsSet[0xDB] = OpcData{1, (*CPU).OpcDB, false, "stp"}
	// ...
	InsSet[0xE2] = OpcData{2, (*CPU).OpcE2, false, "sep.#"}
	// ...
	InsSet[0xE8] = OpcData{1, (*CPU).OpcE8, false, "inx"}
	// ...
	InsSet[0xEA] = OpcData{1, (*CPU).OpcEA, false, "nop"}
	InsSet[0xEB] = OpcData{1, (*CPU).OpcEB, false, "xba"}
	// ...
	InsSet[0xF4] = OpcData{3, (*CPU).OpcF4, false, "phe.#"}
	// ...
	InsSet[0xFA] = OpcData{1, (*CPU).OpcFA, false, "plx"}
	InsSet[0xFB] = OpcData{1, (*CPU).OpcFB, false, "xce"}
}

// --- Instruction functions ---

// ---- 0000 ----

func (c *CPU) Opc00() error { // brk
	fmt.Println("OPC: DUMMY: Executing brk (00) at", c.PC.HexString())
	return fmt.Errorf("brk (0x00): shouldn't be here")
}

func (c *CPU) Opc01() error { // ora.dxi
	fmt.Println("OPC: DUMMY: Executing ora.dxi (02) at", c.PC.HexString())
	return fmt.Errorf("ora.dxi (0x01): shouldn't be here")
}

func (c *CPU) Opc02() error { // cop
	fmt.Println("OPC: DUMMY: Executing cop (03)")
	return fmt.Errorf("cop (0x02): shouldn't be here")
}

// ...

// ---- 1111 ----

func (c *CPU) Opc18() error { // clc
	c.FlagC = CLEAR
	c.PC++
	return nil
}

// ---- 2222 ----

func (c *CPU) Opc20() error { // jsr p. 362

	addr, err := c.modeAbsolute()
	if err != nil {
		return fmt.Errorf("jsr (0x20): couldn't fetch address from %s: %v",
			addr.HexString(), err)
	}

	// Push the PC to the stack: The PC plus the length of this instruction
	// minus one (in other words, plus 2). We push MSB first
	returnPC := c.PC + 2
	err = c.pushData16(common.Data16(returnPC))
	if err != nil {
		return fmt.Errorf("jsr (0x20): couldn't push address to stack: %v", err)
	}

	c.PC = common.Addr16(addr)

	return nil
}

// ---- 3333 ----

func (c *CPU) Opc38() error { // sec
	c.FlagC = SET
	c.PC++
	return nil
}

// ---- 4444 ----

func (c *CPU) Opc40() error { // rti p. 391 (note wrong drawings)
	// The full rti instruction pulls a different number of bytes depending
	// on if we are in native or emulated mode. Since we only work in native
	// mode, we throw an error if the emulated flag is set.

	if c.FlagE == SET {
		return fmt.Errorf("rti (0x40): emulation flag set, we only support native")
	}

	// We pull four bytes: The status register, the PC, and the PBR
	b, err := c.pullByte()
	if err != nil {
		return fmt.Errorf("rti (0x40): can't retrieve Status Register: %v", err)
	}

	c.SetStatReg(b)

	d16, err := c.pullData16()
	if err != nil {
		return fmt.Errorf("rti (0x40): can't retrieve PC: %v", err)
	}

	// In contrast to RTS, we don't need to adjust the PC, just store it as
	// it is
	c.PC = common.Addr16(d16)

	d8, err := c.pullData8()
	if err != nil {
		return fmt.Errorf("rti (0x40): can't retrieve PBR: %v", err)
	}

	c.PBR = d8

	return nil
}

func (c *CPU) Opc4B() error { // phk (Program Bank Register)
	err := c.pushData8(c.PBR)
	if err != nil {
		return fmt.Errorf("phk (0x4B): couldn't push byte to stack: %v", err)
	}

	c.PC++

	return nil
}

func (c *CPU) Opc4C() error { // jmp p. 360

	addr, err := c.getNextData16()
	if err != nil {
		return fmt.Errorf("jmp (4C): couldn't fetch address from %s: %v", addr.HexString(), err)
	}

	// We need to subtract three from the address because the loop will add
	// three
	c.PC = common.Addr16(addr)

	return nil
}

// ...

func (c *CPU) Opc48() error { // pha p. 375
	var err error

	switch c.WidthA {

	case W8:
		err = c.pushData8(c.A8)
		if err != nil {
			return fmt.Errorf("pha (0x48) 8 bit: couldn't push byte %s to stack: %v",
				c.A8.HexString(), err)
		}

	case W16:
		err = c.pushData16(c.A16)
		if err != nil {
			return fmt.Errorf("pha (0x48) 16 bit: couldn't push %s to stack: %v",
				c.A16.HexString(), err)
		}

	default:
		return fmt.Errorf("pha (0x48): illegal width for register A:%d", c.WidthA)
	}

	c.PC++

	return nil
}

// ---- 5555 ----

func (c *CPU) Opc58() error { // cli
	c.FlagI = CLEAR
	c.PC++
	return nil
}

// ---- 6666 ----

func (c *CPU) Opc60() error { // rts p. 394

	// The LSB ist stored on top of the stack
	lsb, err := c.pullByte()
	if err != nil {
		return fmt.Errorf("rts (0x60): couldn't get LSB off the stack: %v", err)
	}

	msb, err := c.pullByte()
	if err != nil {
		return fmt.Errorf("rts (0x60): couldn't get MSB off the stack: %v", err)
	}

	c.PC = (common.Addr16(msb) << 8) | common.Addr16(lsb)

	// Remember we saved one byte less
	c.PC++

	return nil
}

func (c *CPU) Opc68() error { // pla p. 382
	switch c.WidthA {

	case W8:
		d, err := c.pullData8()
		if err != nil {
			return fmt.Errorf("pla (0x68) A8: couldn't get byte off the stack: %v", err)
		}

		c.A8 = d
		c.TestNZ8(c.A8)

	case W16:
		d, err := c.pullData16()
		if err != nil {
			return fmt.Errorf("pla (0x68) A16: couldn't get word off the stack: %v", err)
		}

		c.A16 = d
		c.TestNZ16(c.A16)

	default:
		return fmt.Errorf("pla (0x68): illegal width for register A:%d", c.WidthA)

	}

	c.PC++
	return nil
}

// ---- 7777 ----

func (c *CPU) Opc78() error { // sei
	c.FlagI = SET
	c.PC++
	return nil
}

func (c *CPU) Opc7A() error { // ply
	switch c.WidthXY {

	case W8:
		d, err := c.pullData8()
		if err != nil {
			return fmt.Errorf("ply (0x7A) XY8: couldn't get byte off the stack: %v", err)
		}

		c.Y8 = d
		c.TestNZ8(c.Y8)

	case W16:
		d, err := c.pullData16()
		if err != nil {
			return fmt.Errorf("ply (0x7A) XY16: couldn't get word off the stack: %v", err)
		}

		c.Y16 = d
		c.TestNZ16(c.Y16)

	default:
		return fmt.Errorf("ply (0x7A): illegal width for register Y:%d", c.WidthXY)

	}

	c.PC++
	return nil
}

// ---- 8888 ----

func (c *CPU) Opc80() error { // bra
	b, err := c.getNextByte()
	if err != nil {
		return fmt.Errorf("bra (0x80): couldn't get offset: %v", err)
	}

	addr, err := c.modeBranch(b)
	if err != nil {
		return fmt.Errorf("bra (0x80): branch target wrong: %v", err)
	}

	c.PC = addr

	return nil
}

// ...

func (c *CPU) Opc85() error { // sta.d

	addr, err := c.modeDirectPage()
	if err != nil {
		return fmt.Errorf("sta.d (0x85): couldn't fetch address from %s: %v", addr.HexString(), err)
	}

	err = c.storeA(addr)
	if err != nil {
		return fmt.Errorf("sta.d (0x85): couldn't store A at address %s: %v", addr.HexString(), err)
	}

	c.PC += 2

	return nil
}

func (c *CPU) Opc8A() error { // txa
	txaFNS[c.WidthA][c.WidthXY](c)
	c.PC += 1
	return nil
}

func (c *CPU) Opc8B() error { // phb (Data Bank Register)
	err := c.pushData8(c.DBR)
	if err != nil {
		return fmt.Errorf("phb (0x8B): couldn't push byte to stack: %v", err)
	}

	c.PC++

	return nil
}

func (c *CPU) Opc8D() error { // sta

	addr, err := c.modeAbsolute()
	if err != nil {
		return fmt.Errorf("sta (8D): couldn't fetch address from %s: %v", addr.HexString(), err)
	}

	err = c.storeA(addr)
	if err != nil {
		return fmt.Errorf("sta (8D): couldn't store A at address %s: %v", addr.HexString(), err)
	}

	c.PC += 3

	return nil
}

// ---- 9999 ----

// txs (0x9a) transfers the X register to the Stack Pointer. Note that if we
// have 8-bit index registers, the MSB of the Stack Pointer is zeroed, not set
// to 01. See p. 416 for details. Flags are not affected by this operation. Note
// that internally, the Stack Pointer type is an an Addr16, while X is a
// Data8/Data16
func (c *CPU) Opc9A() error { // txs
	switch c.WidthXY {

	case W8:
		c.SP = common.Addr16(0x0000 + common.Data16(c.X8)) // emphasise MSB is zero

	case W16:
		c.SP = common.Addr16(c.X16)

	default: // paranoid
		return fmt.Errorf("txs (0x9A): Illegal width for register X:%d", c.WidthXY)
	}
	c.PC++

	return nil
}

func (c *CPU) Opc9B() error { // txy
	switch c.WidthXY {

	case W8:
		c.Y8 = c.X8
		c.TestNZ8(c.X8)

	case W16:
		c.Y16 = c.X16
		c.TestNZ16(c.X16)

	default: // paranoid
		return fmt.Errorf("txy (0x9B): Illegal width for register X or Y:%d", c.WidthXY)
	}
	c.PC++

	return nil
}

// ---- AAAA ----

func (c *CPU) OpcA8() error { // tay
	tayFNS[c.WidthA][c.WidthXY](c)
	c.PC += 1
	return nil
}

func (c *CPU) OpcA9() error { // lda.# (lda.8/lda.16)

	switch c.WidthA {
	case W8:
		operand, err := c.modeImmediate8()
		if err != nil {
			return fmt.Errorf("lda.8 (A9): Couldn't fetch data: %v", err)
		}

		c.A8 = operand
		c.TestNZ8(c.A8)
		c.PC += 2

	case W16:
		operand, err := c.modeImmediate16()
		if err != nil {
			return fmt.Errorf("lda.16 (A9): Couldn't fetch data: %v", err)
		}

		c.A16 = operand
		c.TestNZ16(c.A16)
		c.PC += 3

	default: // paranoid
		return fmt.Errorf("lda.# (0xA9): Illegal width for register A:%d", c.WidthA)
	}

	return nil
}

func (c *CPU) OpcAA() error { // tax
	taxFNS[c.WidthA][c.WidthXY](c)
	c.PC += 1
	return nil
}

func (c *CPU) OpcAD() error { // lda
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

	c.PC += 3
	return nil
}

// ...

// ---- BBBB ----

func (c *CPU) OpcB8() error { // clv
	c.FlagV = CLEAR
	c.PC++
	return nil
}

func (c *CPU) OpcBB() error { // tyx
	switch c.WidthXY {

	case W8:
		c.X8 = c.Y8
		c.TestNZ8(c.Y8)

	case W16:
		c.X16 = c.Y16
		c.TestNZ16(c.Y16)

	default: // paranoid
		return fmt.Errorf("tyx (0xBB): Illegal width for register X or Y:%d", c.WidthXY)
	}
	c.PC++

	return nil
}

// ...

// ---- CCCC ----

func (c *CPU) OpcC2() error { // rep.#
	rb, err := c.getNextByte()
	if err != nil {
		return fmt.Errorf("rep.# (0xC2): can't get next byte: %v:", err)
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

	c.PC += 2

	return nil
}

func (c *CPU) OpcC8() error { // iny
	switch c.WidthXY {

	case W8:
		c.Y8 += 1
		c.TestNZ8(c.Y8)

	case W16:
		c.Y16 += 1
		c.TestNZ16(c.Y16)

	default: // paranoid
		return fmt.Errorf("iny (0xC8): illegal WidthXY value: %v", c.WidthXY)
	}

	c.PC++
	return nil
}

// ...

// ---- DDDD ----

func (c *CPU) OpcD8() error { // cld
	c.FlagD = CLEAR
	c.PC++
	return nil
}

// ...

func (c *CPU) OpcDB() error { // stp
	// TODO make sure we actually add one to the PC here
	c.IsStopped = true
	fmt.Println("Machine stopped by stp (0xDB) in block", c.PBR.HexString(), "at address", c.PC.HexString())
	c.PC++
	return nil
}

// ...

// ---- EEEE ----

func (c *CPU) OpcE2() error { // sep.#

	rb, err := c.getNextByte()
	if err != nil {
		return fmt.Errorf("sep.# (0xE2): can't get next byte: %v:", err)
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

	c.PC += 2

	return nil
}

func (c *CPU) OpcE8() error { // inx
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

	c.PC++
	return nil
}

func (c *CPU) OpcEA() error { // nop
	c.PC++
	return nil
}

func (c *CPU) OpcEB() error { // xba p.422
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

	c.PC++
	return nil
}

// ---- FFFF ----

func (c *CPU) OpcF4() error { // phe.#
	d, err := c.getNextData16()
	if err != nil {
		return fmt.Errorf("phe.# (0xF4): couldn't get operand: %v", err)
	}

	err = c.pushData16(d)
	if err != nil {
		return fmt.Errorf("phe.# (0xF4): couldn't push operand: %v", err)
	}

	c.PC += 3
	return nil
}

func (c *CPU) OpcFA() error { // plx
	switch c.WidthXY {

	case W8:
		d, err := c.pullData8()
		if err != nil {
			return fmt.Errorf("plx (0xFA) XY8: couldn't get byte off the stack: %v", err)
		}

		c.X8 = d
		c.TestNZ8(c.X8)

	case W16:
		d, err := c.pullData16()
		if err != nil {
			return fmt.Errorf("plx (0xFA) XY16: couldn't get word off the stack: %v", err)
		}

		c.X16 = d
		c.TestNZ16(c.X16)

	default:
		return fmt.Errorf("plx (0xFA): illegal width for register X:%d", c.WidthXY)

	}

	c.PC++
	return nil
}

func (c *CPU) OpcFB() error { // xce
	tmp := c.FlagE
	c.FlagE = c.FlagC
	c.FlagC = tmp
	c.PC++
	return nil
}
