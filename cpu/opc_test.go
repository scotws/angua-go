// Test file for Angua CPU opcodes
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 05. Jan 2019
// This version: 19. Jan 2019

package cpu

import (
	"testing"

	"angua/common"
	"angua/mem"
)

// ===== MODE TESTS =====

func TestModeBranch(t *testing.T) {
	var c = CPU{}
	var m = mem.Memory{}
	nc, _ := mem.NewChunk(0000, 0xFFFF, "ram")
	m.Chunks = append(m.Chunks, nc)
	c.Mem = &m

	var tests = []struct {
		pc     common.Addr16
		offset byte
		addr   common.Addr16
	}{
		{0x0000, 0x00, 0x0002},
		{0x0002, 0x01, 0x0005},
		{0x0003, 0xFB, 0x0000},
		{0x0004, 0xFA, 0x0000},
		{0x0005, 0xFC, 0x0003},
	}

	for _, test := range tests {
		c.PC = test.pc
		got, _ := c.modeBranch(test.offset)
		if got != test.addr {
			t.Errorf("TestModeBranch for 0x%20X returns %X, wanted %X",
				test.offset, got, test.addr)
		}
	}
}

// ===== TESTS FOR INDIVIDUAL INSTRUCTIONS =====

// clc, cld, cli, clv, sec, sei, sev
func TestOpcFlags(t *testing.T) {
	var c = CPU{}
	var tests = []struct {
		name string
		init func() error
		flag *byte
		want byte
	}{
		{"0x18 (clc)", c.Opc18, &c.FlagC, CLEAR},
		{"0xD8 (cld)", c.OpcD8, &c.FlagD, CLEAR},
		{"0x58 (cli)", c.Opc58, &c.FlagI, CLEAR},
		{"0xB8 (clv)", c.OpcB8, &c.FlagV, CLEAR},

		{"0x38 (sec)", c.Opc38, &c.FlagC, SET},
		// {"0x?? (sed)", Opc??, &c.FlagD, SET}, TODO
		{"0x78 (sei)", c.Opc78, &c.FlagI, SET},
		// {"0x?? (sev)", Opc??, &c.FlagV, SET}, TODO
	}

	for _, test := range tests {
		test.init() // executes opcode

		if *test.flag != test.want {
			t.Errorf("TestOpcFlags for %v returns %X, wanted %X",
				test.name, *test.flag, test.want)
		}

	}
}

// txs. Does not affect flags
func TestOpc9A(t *testing.T) {
	var c = CPU{}

	// Test with 8 bit X

	var tests8 = []struct {
		init func() error
		have common.Data8
		want common.Addr16 // SP is Addr16
	}{
		{c.Opc9A, 0x01, 0x0001},
	}

	c.WidthXY = W8

	for _, test8 := range tests8 {
		c.X8 = test8.have
		_ = test8.init() // executes opcode, dump err

		if c.SP != test8.want {
			t.Errorf("TestOpc9A (txs) 8 returns %X for c.SP, wanted %X", c.SP, test8.want)
		}
	}

	// Test with 16 bit X

	var tests16 = []struct {
		init func() error
		have common.Data16
		want common.Addr16 // SP is Addr16
	}{
		{c.Opc9A, 0x01, 0x0001},
	}

	c.WidthXY = W16

	for _, test16 := range tests16 {
		c.X16 = test16.have
		_ = test16.init() // executes opcode, dump err

		if c.SP != test16.want {
			t.Errorf("TestOpc9A (txs) 16 returns %X for c.SP, wanted %X", c.SP, test16.want)
		}
	}
}

// xba
func TestOpcEB(t *testing.T) {
	var c = CPU{}

	// First step: Test with 8 bit A

	var tests8 = []struct {
		init func() error
		have common.Data8
		want common.Data8
	}{
		{c.OpcEB, 0xFF, 0xFF},
		{c.OpcEB, 0x00, 0x00},
	}

	c.WidthA = W8 // We start in native mode otherwise

	for _, test8 := range tests8 {
		c.A8 = test8.have
		_ = test8.init() // executes opcode, dump err

		if c.B != test8.want {
			t.Errorf("TestOpcEB (xba) A8 returns %X for B and %X for A, wanted %X", c.B, c.A8, test8.want)
		}
	}

	// Second step: Test with 16 bit A

	var tests16 = []struct {
		init func() error
		have common.Data16
		want common.Data16
	}{
		{c.OpcEB, 0xFF00, 0x00FF},
		{c.OpcEB, 0x0000, 0x0000},
	}

	c.WidthA = W16

	for _, test16 := range tests16 {
		c.A16 = test16.have
		_ = test16.init() // executes opcode

		if c.A16 != test16.want {
			t.Errorf("TestOpcEB (xba) A16 returns %X, wanted %X", c.A16, test16.want)
		}
	}

}

// tax (0xAA). The test currently doesn't check flags
func TestOpcAA(t *testing.T) {
	var c = CPU{}

	// First test with A8 and X8

	var tests8n8 = []struct {
		init func() error
		have common.Data8
		want common.Data8
	}{
		{c.OpcAA, 0x00, 0x00},
		{c.OpcAA, 0xAA, 0xAA},
		{c.OpcAA, 0xFF, 0xFF},
	}

	c.WidthA = W8
	c.WidthXY = W8

	for _, test8 := range tests8n8 {
		c.A8 = test8.have
		_ = test8.init() // executes opcode

		if c.X8 != test8.want {
			t.Errorf("TestOpcAA (tax) A8X8 with %X returns %X, wanted %X",
				test8.have, c.X8, test8.want)
		}
	}

	// Second test with A8 and X16

	var tests8n16 = []struct {
		init func() error
		have common.Data8
		want common.Data16
	}{
		{c.OpcAA, 0x00, 0x0000},
		{c.OpcAA, 0xAA, 0x00AA},
		{c.OpcAA, 0xFF, 0x00FF},
	}

	c.WidthA = W8
	c.WidthXY = W16

	for _, test8 := range tests8n16 {
		c.A8 = test8.have
		_ = test8.init() // executes opcode

		if c.X16 != test8.want {
			t.Errorf("TestOpcAA (tax) A8X16 with %X returns %X, wanted %X",
				test8.have, c.X16, test8.want)
		}
	}

	// Third test with A16 and X8

	var tests16n8 = []struct {
		init func() error
		have common.Data16
		want common.Data8
	}{
		{c.OpcAA, 0x0000, 0x00},
		{c.OpcAA, 0x00AA, 0xAA},
		{c.OpcAA, 0x00FF, 0xFF},
		{c.OpcAA, 0x0FF0, 0xF0},
		{c.OpcAA, 0xFF00, 0x00},
	}

	c.WidthA = W16
	c.WidthXY = W8

	for _, test8 := range tests16n8 {
		c.A16 = test8.have
		_ = test8.init() // executes opcode

		if c.X8 != test8.want {
			t.Errorf("TestOpcAA (tax) A8X16 with %X returns %X, wanted %X",
				test8.have, c.X8, test8.want)
		}
	}

	// Fourth test with A16 and X16

	var tests16n16 = []struct {
		init func() error
		have common.Data16
		want common.Data16
	}{
		{c.OpcAA, 0x0000, 0x0000},
		{c.OpcAA, 0x00AA, 0x00AA},
		{c.OpcAA, 0x00FF, 0x00FF},
		{c.OpcAA, 0x0FF0, 0x0FF0},
		{c.OpcAA, 0xFF00, 0xFF00},
	}

	c.WidthA = W16
	c.WidthXY = W16

	for _, test8 := range tests16n16 {
		c.A16 = test8.have
		_ = test8.init() // executes opcode

		if c.X16 != test8.want {
			t.Errorf("TestOpcAA (tax) A8X16 with %X returns %X, wanted %X",
				test8.have, c.X16, test8.want)
		}
	}

}

// txa (0x8A). The test currently doesn't check flags
func TestOpc8A(t *testing.T) {
	var c = CPU{}

	// First test with X8 and A8

	var tests8n8 = []struct {
		init func() error
		have common.Data8
		want common.Data8
	}{
		{c.Opc8A, 0x00, 0x00},
		{c.Opc8A, 0xAA, 0xAA},
		{c.Opc8A, 0xFF, 0xFF},
	}

	c.WidthXY = W8
	c.WidthA = W8

	for _, test8 := range tests8n8 {
		c.X8 = test8.have
		_ = test8.init() // executes opcode

		if c.A8 != test8.want {
			t.Errorf("TestOpc8A (txa) X8A8 with %X returns %X, wanted %X",
				test8.have, c.A8, test8.want)
		}
	}

	// Second test with X8 and A16

	var tests8n16 = []struct {
		init func() error
		have common.Data8
		want common.Data16
	}{
		{c.Opc8A, 0x00, 0x0000},
		{c.Opc8A, 0xAA, 0x00AA},
		{c.Opc8A, 0xFF, 0x00FF},
	}

	c.WidthXY = W8
	c.WidthA = W16

	for _, test8 := range tests8n16 {
		c.X8 = test8.have
		_ = test8.init() // executes opcode

		if c.A16 != test8.want {
			t.Errorf("TestOpc8A (txa) X8A16 with %X returns %X, wanted %X",
				test8.have, c.A16, test8.want)
		}
	}

	// Third test with X16 and A8

	var tests16n8 = []struct {
		init func() error
		have common.Data16
		want common.Data8
	}{
		{c.Opc8A, 0x0000, 0x00},
		{c.Opc8A, 0x00AA, 0xAA},
		{c.Opc8A, 0x00FF, 0xFF},
		{c.Opc8A, 0x0FF0, 0xF0},
		{c.Opc8A, 0xFF00, 0x00},
	}

	c.WidthXY = W16
	c.WidthA = W8

	for _, test8 := range tests16n8 {
		c.X16 = test8.have
		_ = test8.init() // executes opcode

		if c.A8 != test8.want {
			t.Errorf("TestOpc8A (txa) X8A16 with %X returns %X, wanted %X",
				test8.have, c.A8, test8.want)
		}
	}

	// Fourth test with X16 and A16

	var tests16n16 = []struct {
		init func() error
		have common.Data16
		want common.Data16
	}{
		{c.Opc8A, 0x0000, 0x0000},
		{c.Opc8A, 0x00AA, 0x00AA},
		{c.Opc8A, 0x00FF, 0x00FF},
		{c.Opc8A, 0x0FF0, 0x0FF0},
		{c.Opc8A, 0xFF00, 0xFF00},
	}

	c.WidthXY = W16
	c.WidthA = W16

	for _, test8 := range tests16n16 {
		c.X16 = test8.have
		_ = test8.init() // executes opcode

		if c.A16 != test8.want {
			t.Errorf("TestOpc8A (txa) X8A16 with %X returns %X, wanted %X",
				test8.have, c.A16, test8.want)
		}
	}

}
