// Opcodes for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 02. Jan 2019
// This version: 12. Jan 2019

// This package contains the opcodes and opcode data for the 65816 instructions.
// Note there is redundancy with the information in the info package. We keep
// these separate for the moment to allow more robustness while editing.

package cpu

import (
	"fmt"

	"angua/common"
)

type OpcData struct {
	Size int              // number of bytes including operand (1 to 4)
	Code func(*CPU) error // function for actual code of instruction
}

// In theory, we could use a giant switch statement instead of a map of
// functions. However,
// https://hashrocket.com/blog/posts/switch-vs-map-which-is-the-better-way-to-branch-in-go
// suggests that maps are faster
var (
	InsSet [256]OpcData
)

func init() {
	InsSet[0x00] = OpcData{2, Opc00} // brk with signature byte
	InsSet[0x01] = OpcData{2, Opc01} // ora.dxi
	InsSet[0x02] = OpcData{2, Opc02} // cop
	// ...
	InsSet[0x18] = OpcData{1, Opc18} // clc
	// ...
	InsSet[0x38] = OpcData{1, Opc38} // sec
	// ...
	InsSet[0x48] = OpcData{1, Opc48} // pha
	// ...
	InsSet[0x58] = OpcData{1, Opc58} // cli
	// ...
	InsSet[0x78] = OpcData{1, Opc78} // sei
	// ...
	InsSet[0x85] = OpcData{2, Opc85} // sta.d
	// ...
	InsSet[0x8D] = OpcData{3, Opc8D} // sta
	// ...
	InsSet[0x9A] = OpcData{1, Opc9A} // txs
	// ...
	InsSet[0xA9] = OpcData{2, OpcA9} // lda.# (lda.8/lda.16)
	// ...
	InsSet[0xAD] = OpcData{3, OpcAD} // lda
	// ...
	InsSet[0xB8] = OpcData{1, OpcB8} // clv
	// ...
	InsSet[0xC2] = OpcData{2, OpcC2} // rep.#
	// ...
	InsSet[0xD8] = OpcData{1, OpcD8} // cld
	// ...
	InsSet[0xDB] = OpcData{1, OpcDB} // stp
	// ...
	InsSet[0xE2] = OpcData{2, OpcE2} // sep.#
	// ...
	InsSet[0xE8] = OpcData{1, OpcE8} // inx
	// ...
	InsSet[0xEA] = OpcData{1, OpcEA} // nop
	InsSet[0xEB] = OpcData{1, OpcEB} // xba
	// ...
	InsSet[0xF4] = OpcData{3, OpcF4} // pha.#
	// ...
	InsSet[0xFB] = OpcData{1, OpcFB} // xce
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

func (c *CPU) getNextData8() (common.Data8, error) {
	return c.modeImmediate8()
}

// modeImmediate16 returns the word stored in the address after the opcode and an
// error. This could also be named getNextData16()
func (c *CPU) modeImmediate16() (common.Data16, error) {
	operandAddr := c.getFullPC() + 1
	ui, err := c.Mem.FetchMore(operandAddr, 2)
	if err != nil {
		return 0, fmt.Errorf("immediate 16 mode: couldn't fetch data from %s: %v", common.Addr24(operandAddr).HexString(), err)
	}

	d := common.Data16(ui)

	return d, nil
}

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

// ---- 3333 ----

func Opc38(c *CPU) error { // sec
	c.FlagC = SET
	return nil
}

// ---- 4444 ----

func Opc48(c *CPU) error { // pha p. 375
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

	return nil
}

// ---- 5555 ----

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

// ---- 9999 ----

// txs (0x9a) transfers the X register to the Stack Pointer. Note that if we
// have 8-bit index registers, the MSB of the Stack Pointer is zeroed, not set
// to 01. See p. 416 for details. Flags are not affected by this operation. Note
// that internally, the Stack Pointer type is an an Addr16, while X is a
// Data8/Data16
func Opc9A(c *CPU) error { // txs
	switch c.WidthXY {

	case W8:
		c.SP = common.Addr16(0x0000 + common.Data16(c.X8)) // emphasise MSB is zero

	case W16:
		c.SP = common.Addr16(c.X16)

	default: // paranoid
		return fmt.Errorf("txs (0x9A): Illegal width for register X:%d", c.WidthXY)
	}

	return nil
}

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

func OpcC2(c *CPU) error { // rep.#
	rb, err := getNextByte(c)
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
	fmt.Println("Machine stopped by stp (0xDB) in block", c.PBR.HexString(), "at address", c.PC.HexString())
	return nil
}

// ...

// ---- EEEE ----

func OpcE2(c *CPU) error { // sep.#

	rb, err := getNextByte(c)
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

func OpcF4(c *CPU) error { // phe.#
	d, err := c.modeImmediate16()
	if err != nil {
		return fmt.Errorf("phe.# (0xF4): couldn't get operand: %v", err)
	}

	err = c.pushData16(d)
	if err != nil {
		return fmt.Errorf("phe.# (0xF4): couldn't push operand: %v", err)
	}

	return nil
}

func OpcFB(c *CPU) error { // xce
	tmp := c.FlagE
	c.FlagE = c.FlagC
	c.FlagC = tmp
	return nil
}
