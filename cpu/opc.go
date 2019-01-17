// Opcodes for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 02. Jan 2019
// This version: 16. Jan 2019

// This package contains the opcodes and opcode data for the 65816 instructions.
// Note there is redundancy with the information in the info package. We keep
// these separate for the moment to allow more robustness while editing.

package cpu

import (
	"fmt"

	"angua/common"
)

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
	InsSet[0x8B] = OpcData{1, (*CPU).Opc8B, false, "phb"} // pushes DBR
	// ...
	InsSet[0x8D] = OpcData{3, (*CPU).Opc8D, false, "sta"}
	// ...
	InsSet[0x9A] = OpcData{1, (*CPU).Opc9A, false, "txs"}
	InsSet[0x9B] = OpcData{1, (*CPU).Opc9B, false, "txy"}
	// ...
	InsSet[0xA9] = OpcData{2, (*CPU).OpcA9, true, "lda.#"} // includes lda.8/lda.16
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
		return 0, fmt.Errorf("direct page mode: couldn't fetch address from %s: %v", common.Addr24(operandAddr).HexString(), err)
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
		return 0, fmt.Errorf("immediate 8 mode: couldn't fetch data from %s: %v", common.Addr24(operandAddr).HexString(), err)
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
		return 0, fmt.Errorf("immediate 16 mode: couldn't fetch data from %s: %v", common.Addr24(operandAddr).HexString(), err)
	}

	return common.Data16(ui), nil
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
		return 0, fmt.Errorf("getNextByte: couldn't fetch data from %s: %v", common.Addr24(byteAddr).HexString(), err)
	}

	return b, nil
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

// ...

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
