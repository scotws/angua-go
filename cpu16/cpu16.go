// Angua CPU System - Native Mode CPU
// Scot W. Stevenson
// First version: 06. Nov 2018
// First version: 06. Nov 2018

// TODO this package will be worked on when cpu8 is working

package cpu16

type reg8 uint8
type reg16 uint16

type Cpu16 struct {
	A reg16
	X reg16
	Y reg16

	DP reg16 // Direct Page register, 16 bit, not 8
	SP reg16 // Stack Pointer, 8 bit

	P byte // Status Register

	DBR reg8 // Data Bank Register, available in emulated mode
	PBR reg8 // Program Bank Register, available in emulated mode

	PC reg16 // Program counter

}
