// Opcodes for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// This version: 02. Jan 2019
// First version: 02. Jan 2019

// This package contains the opcodes and opcode data for the 65816 instructions.
// Note there is redundancy with the information in the info package. We keep
// these separate to allow more robustness while editing.

package opc

import (
	"fmt" // TODO This is used for testing only
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

var (
	InsJump [256]func()  // Instruction jump table
	InsData [256]OpcData // Instruction data
)

func init() {
	// Instruction jumps
	InsData[0x00] = OpcData{2, STACK, false}    // brk with signature byte
	InsData[0x01] = OpcData{2, DP_IND_X, false} // ora.dxi
	InsData[0x02] = OpcData{2, STACK, false}    // cop
	// ...
	InsData[0x18] = OpcData{1, IMPLIED, false} // clc
	// ...
	InsData[0x85] = OpcData{2, DP, false} // sta.d
	// ...
	InsData[0xA9] = OpcData{2, IMMEDIATE, true} // lda.# (lda.8/lda.16)
	// ...
	InsData[0xAD] = OpcData{3, ABSOLUTE, false} // lda
	// ...
	InsData[0xEA] = OpcData{1, IMPLIED, false} // nop
	// ...
	InsData[0xDB] = OpcData{1, IMPLIED, false} // stp
	// ...
	InsData[0xFB] = OpcData{1, IMPLIED, false} // xce

	// Instruction data
	InsJump[0x00] = Opc00 // brk
	InsJump[0x01] = Opc01 // ora.dxi
	InsJump[0x02] = Opc02 // cop
	// ...
	InsJump[0x18] = Opc18 // clc
	// ...
	InsJump[0xA9] = OpcA9 // lda.# (lda.8/lda.16)
	// ...
	InsJump[0x85] = Opc85 // sta.d
	// ...
	InsJump[0xAD] = OpcAD // lda
	// ...
	InsJump[0xDB] = OpcDB // stp
	// ...
	InsJump[0xEA] = OpcEA // nop
	// ...
	InsJump[0xFB] = OpcFB // xce
}

// --- Instruction Functions ---

func Opc00() { // brk
	fmt.Println("OPC: DUMMY: Executing brk (00)")
}

func Opc01() { // ora.dxi
	fmt.Println("OPC: DUMMY: Executing ora.dxi (02)")
}

func Opc02() { // cop
	fmt.Println("OPC: DUMMY: Executing cop (03)")
}

func Opc18() { // nop
	fmt.Println("OPC: DUMMY: Executing clc (18) ")
}

// ...

func Opc85() { // sta.d
	fmt.Println("OPC: DUMMY: Executing sta.d (85) ")
}

// ...

func OpcA9() { // lda.# (lda.8/lda.16)
	fmt.Println("OPC: DUMMY: Executing lda.# (A9) ")
}

// ...

func OpcAD() { // lda
	fmt.Println("OPC: DUMMY: Executing lda (AD) ")
}

// ...

func OpcDB() { // stp
	fmt.Println("OPC: DUMMY: Executing stp (DB) ")
}

// ...

func OpcEA() { // nop
	fmt.Println("OPC: DUMMY: Executing nop (EA) ")
}

// ...

func OpcFB() { // nop
	fmt.Println("OPC: DUMMY: Executing xce (FB) ")
}
