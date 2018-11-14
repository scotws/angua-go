// Info (built-in help) system for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 11. Nov 2018
// This version: 14. Nov 2018

// info is a package to provide an online manual system to look up instructions
// while using the Angua emulator for the 65816. As a spin off, it allows the
// export as a JSON file for other uses
package info

import (
	"fmt"
)

const (
	// Addressing modes
	ProgCountRel     = 1
	ProgCountRelLong = 2
	StackInterrupt   = 3
	Implied          = 4
	Immediate        = 5
	Accumulator      = 6
	Absolute         = 7
	DirPageIndexX    = 8
	DirPage          = 9
	StackRel         = 10

	// Strings
	mpu6502  = "6502"
	mpu65c02 = "65c02"
	mpu65816 = "65802/65816"
	mpuAll   = "6502 65c02 65802/65816"
)

type InfoEntry struct {
	Opcode       byte     // 65816 machine language opcode
	MneSAN       string   // Simpler Assembler Notation (SAN) mnemonic
	MneTrad      string   // Traditional notation mnemonic
	Name         string   // Official name from Eyes and Lichty
	AddrMode     int      // Addressing mode (see table)
	Bytes        int      // Size in bytes (including opcode)
	Cycles       int      // Minimum number of cycles
	Page         int      // Page number in Eyes and Lichty, (2007)
	MPUs         []string // MPUs it is present on
	Operand      bool     // Takes an operand?
	NativeExpand bool     // Size expands in 65816 native mode
	Flags        []string // Which flags are changed (upper case)
	Note         string   // Other information
}

var (
	OpcodeDict = make(map[byte]InfoEntry)   // Dictionary organized by Opcodes
	SANDict    = make(map[string]InfoEntry) // Dictionary organized by SAN mnemonics
	ModeDict   = make(map[int]string)       // Dictionary to print Addressing Modes
)

// generateDicts creates the Dictionaries on the fly during boot. This is
// intended to run a goroutine while booting
func GenerateDicts() {

	OpcodeDict[0x00] = InfoEntry{0x00, "brk", "BRK", "Software Break", StackInterrupt,
		2, 7, 338, []string{mpu6502, mpu65c02, mpu65816}, false, false,
		[]string{"B", "D", "I"}, "PC incremented by 2 for signature byte"}

	OpcodeDict[0x18] = InfoEntry{0x18, "clc", "CLC", "Clear Carry Flag", Implied,
		1, 2, 343, []string{mpu6502, mpu65c02, mpu65816}, false, false,
		[]string{"C"}, "On 65816/65802, used to switch to Native Mode"}

	OpcodeDict[0x50] = InfoEntry{0x50, "bvc", "BVC", "Branch if Overflow Clear", ProgCountRel,
		2, 2, 341, []string{mpu6502, mpu65c02, mpu65816}, false, false,
		[]string{""}, "Also set by Set Overflow Signal to chip"}

	OpcodeDict[0x70] = InfoEntry{0x70, "bvs", "BVS", "Branch if Overflow Set", ProgCountRel,
		2, 2, 342, []string{mpu6502, mpu65c02, mpu65816}, false, false,
		[]string{""}, "Also set by Set Overflow Signal to chip"}

	OpcodeDict[0x3A] = InfoEntry{0x3A, "dec.a", "DEC A", "Decrement", Accumulator,
		1, 2, 352, []string{mpu6502, mpu65c02, mpu65816}, false, false,
		[]string{"N", "Z"}, "Does not affect carry flag C"}

	OpcodeDict[0x9C] = InfoEntry{0x9C, "stz", "STZ ????", "Store Zero to Memory",
		Absolute, 3, 4, 405, []string{mpu65c02, mpu65816}, true, false,
		[]string{""}, "Flags unaffected by store instructions"}

	OpcodeDict[0x74] = InfoEntry{0x74, "stz.dx", "STZ ??,X", "Store Zero to Memory",
		DirPageIndexX, 2, 4, 405, []string{mpu65c02, mpu65816}, true, false,
		[]string{""}, "Flags unaffected by store instructions"}

	OpcodeDict[0x82] = InfoEntry{0x82, "bra.l", "BRL", "Branch Always Long", ProgCountRelLong,
		3, 4, 340, []string{mpu65816}, true, false,
		[]string{""}, "Relocateable, but JMP Absolute one cycle faster"}

	OpcodeDict[0x8E] = InfoEntry{0x8E, "stx", "STX ????", "Store Index Register X to Memory",
		Absolute, 3, 4, 403, []string{mpuAll}, true, false,
		[]string{""}, "Flags unaffected by store instructions"}

	OpcodeDict[0xA3] = InfoEntry{0xA3, "lda.s", "LDA ??,S", "Load Accumulator from Memory",
		StackRel, 2, 4, 363, []string{mpu65816}, true, false,
		[]string{"N", "Z"}, "(none)"}

	OpcodeDict[0xA9] = InfoEntry{0xA9, "lda.#", "LDA #??", "Load Accumulator from Memory",
		Immediate, 2, 2, 363, []string{mpuAll}, true, true,
		[]string{"N", "Z"}, "On 65802/65816, 16 bit of data compared if flag X clear"}

	OpcodeDict[0xE0] = InfoEntry{0xE0, "cpx.#", "CPX #??", "Compare Index Register X with Memory",
		Immediate, 2, 2, 350, []string{mpu6502, mpu65c02, mpu65816}, true, true,
		[]string{"N", "Z", "C"}, "On 65802/65816, 16 bit of data compared if flag X clear"}

	OpcodeDict[0xEA] = InfoEntry{0xEA, "nop", "NOP", "No Operation",
		Implied, 1, 2, 369, []string{mpu6502, mpu65c02, mpu65816}, false, false,
		[]string{""}, "For more than one nop, see other instructions"}

	SANDict["bra.l"] = OpcodeDict[0x82]
	SANDict["brk"] = OpcodeDict[0x00]
	SANDict["bvc"] = OpcodeDict[0x50]
	SANDict["bvs"] = OpcodeDict[0x70]
	SANDict["clc"] = OpcodeDict[0x18]
	SANDict["cpx.#"] = OpcodeDict[0xE0]
	SANDict["dec.a"] = OpcodeDict[0x3A]
	SANDict["lda.#"] = OpcodeDict[0xA9]
	SANDict["lda.s"] = OpcodeDict[0xA3]
	SANDict["nop"] = OpcodeDict[0xEA]
	SANDict["stx"] = OpcodeDict[0x8E]
	SANDict["stz"] = OpcodeDict[0x9C]
	SANDict["stz.dx"] = OpcodeDict[0x74]

	ModeDict[Absolute] = "Absolute"
	ModeDict[Accumulator] = "Accumulator"
	ModeDict[DirPage] = "Direct Page"
	ModeDict[DirPageIndexX] = "Direct Page Index X"
	ModeDict[Immediate] = "Immediate"
	ModeDict[Implied] = "Implied"
	ModeDict[ProgCountRelLong] = "Program Counter Relative Long"
	ModeDict[ProgCountRel] = "Program Counter Relative"
	ModeDict[StackInterrupt] = "Stack/Interrupt"
	ModeDict[StackRel] = "Stack Relative"
}

// PrintOpcodeInfo takes a byte and prints out a nicely formatted dump of
// information on the instruction such referenced.
// TODO move these series of prints to string/html template
func PrintOpcodeInfo(opc byte) {
	fmt.Print(OpcodeDict[opc].MneSAN, " / ", OpcodeDict[opc].MneTrad,
		" -- ", OpcodeDict[opc].Name)
	fmt.Printf("  (Eyes & Lichty p. %d)\n", OpcodeDict[opc].Page)
	fmt.Println("Mode:", ModeDict[OpcodeDict[opc].AddrMode],
		" Flags:", OpcodeDict[opc].Flags)
	fmt.Println("Size:", OpcodeDict[opc].Bytes, "byte(s)", " Cycles:",
		OpcodeDict[opc].Cycles, " MPUs:", OpcodeDict[opc].MPUs)
	fmt.Println("Note:", OpcodeDict[opc].Note)
}

// ExportJSON exports the complete Dictionary as a JSON file to a file given
func ExportJSON(f string) {
	fmt.Println("DUMMY ExportJSON to", f)
}
