// Manual (built-in help) system Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 11. Nov 2018
// This version: 11. Nov 2018

// man is a package to provide an online manual system to look up instructions
// while using the Angua emulator for the 65816. As a spin off, it allows the
// export as a JSON file for other uses
package man

import (
	"fmt"
)

const (
	// Addressing modes
	ProgCountRel     = 1
	ProgCountRelLong = 2
	StackInterrupt   = 3

	// Strings
	mpu6502  = "6502"
	mpu65c02 = "65c02"
	mpu65816 = "65816" // also 65802
)

type ManEntry struct {
	Opcode        byte     // 65816 machine language opcode
	MneSAN        string   // Simpler Assembler Notation (SAN) mnemonic
	MneTrad       string   // Traditional notation mnemonic
	Name          string   // Official name from Eyes and Lichty
	AddrMode      int      // Addressing mode (see table)
	Bytes         int      // Size in bytes (including opcode)
	Cycles        int      // Minimum number of cycles
	Page          int      // Page number in Eyes and Lichty, (2007)
	MPUs          []string // MPUs it is present on
	Operand       bool     // Takes an operand?
	NativeExpand  bool     // Size expands in 65816 native mode
	FlagsAffected []string // Which flags are changed (upper case)
	Notes         string   // Other information
}

var (
	OpcodeDict = make(map[byte]ManEntry)   // Dictionary organized by Opcodes
	SANDict    = make(map[string]ManEntry) // Dictionary organized by SAN mnemonics
)

// generateDicts creates the Dictionaries on the fly during boot. This is
// intended to run a goroutine while booting
func GenerateDicts() {

	OpcodeDict[0x00] = ManEntry{0x00, "brk", "brk", "Software Break", StackInterrupt,
		2, 7, 338, []string{mpu6502, mpu65c02, mpu65816}, false, false,
		[]string{"B", "D", "I"}, "PC incremented by 2 for signature byte"}

	OpcodeDict[0x50] = ManEntry{0x50, "bvc", "bvc", "Branch if Overflow Clear", ProgCountRel,
		2, 2, 341, []string{mpu6502, mpu65c02, mpu65816}, false, false,
		[]string{""}, "Also set by Set Overflow Signal to chip"}

	OpcodeDict[0x82] = ManEntry{0x82, "bra.l", "brl", "Branch Always Long", ProgCountRelLong,
		3, 4, 340, []string{mpu65816}, true, false,
		[]string{""}, "Relocateable, but JMP absolute one cycle faster"}

	SANDict["brk"] = OpcodeDict[0x00]
	SANDict["bvc"] = OpcodeDict[0x50]
	SANDict["bra.l"] = OpcodeDict[0x82]
}

// PrintOpcodeInfo takes a byte and prints out a nicely formatted dump of
// information on the instruction such referenced.
// TODO move these series of prints to string/html template
func PrintOpcodeInfo(opc byte) {
	fmt.Println(OpcodeDict[opc].MneSAN, OpcodeDict[opc].MneTrad, OpcodeDict[opc].Name)
}

// ExportJSON exports the complete Dictionary as a JSON file to a file given
func ExportJSON(f string) {
	fmt.Println("DUMMY ExportJSON to", f)
}
