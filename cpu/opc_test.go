// Test file for Angua CPU opcodes
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 05. Jan 2019
// This version: 15. Jan 2019

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
		init func(*CPU) error
		flag *byte
		want byte
	}{
		{"0x18 (clc)", Opc18, &c.FlagC, CLEAR},
		{"0xD8 (cld)", OpcD8, &c.FlagD, CLEAR},
		{"0x58 (cli)", Opc58, &c.FlagI, CLEAR},
		{"0xB8 (clv)", OpcB8, &c.FlagV, CLEAR},

		{"0x38 (sec)", Opc38, &c.FlagC, SET},
		// {"0x?? (sed)", Opc??, &c.FlagD, SET}, TODO
		{"0x78 (sei)", Opc78, &c.FlagI, SET},
		// {"0x?? (sev)", Opc??, &c.FlagV, SET}, TODO
	}

	for _, test := range tests {
		test.init(&c) // executes opcode

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
		init func(*CPU) error
		have common.Data8
		want common.Addr16 // SP is Addr16
	}{
		{Opc9A, 0x01, 0x0001},
	}

	c.WidthXY = W8

	for _, test8 := range tests8 {
		c.X8 = test8.have
		_ = test8.init(&c) // executes opcode, dump err

		if c.SP != test8.want {
			t.Errorf("TestOpc9A (txs) 8 returns %X for c.SP, wanted %X", c.SP, test8.want)
		}
	}

	// Test with 16 bit X

	var tests16 = []struct {
		init func(*CPU) error
		have common.Data16
		want common.Addr16 // SP is Addr16
	}{
		{Opc9A, 0x01, 0x0001},
	}

	c.WidthXY = W16

	for _, test16 := range tests16 {
		c.X16 = test16.have
		_ = test16.init(&c) // executes opcode, dump err

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
		init func(*CPU) error
		have common.Data8
		want common.Data8
	}{
		{OpcEB, 0xFF, 0xFF},
		{OpcEB, 0x00, 0x00},
	}

	c.WidthA = W8 // We start in native mode otherwise

	for _, test8 := range tests8 {
		c.A8 = test8.have
		_ = test8.init(&c) // executes opcode, dump err

		if c.B != test8.want {
			t.Errorf("TestOpcEB (xba) A8 returns %X for B and %X for A, wanted %X", c.B, c.A8, test8.want)
		}
	}

	// Second step: Test with 16 bit A

	var tests16 = []struct {
		init func(*CPU) error
		have common.Data16
		want common.Data16
	}{
		{OpcEB, 0xFF00, 0x00FF},
		{OpcEB, 0x0000, 0x0000},
	}

	c.WidthA = W16

	for _, test16 := range tests16 {
		c.A16 = test16.have
		_ = test16.init(&c) // executes opcode

		if c.A16 != test16.want {
			t.Errorf("TestOpcEB (xba) A16 returns %X, wanted %X", c.A16, test16.want)
		}
	}

}
