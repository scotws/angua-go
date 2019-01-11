// Test file for Angua CPU opcodes
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 05. Jan 2019
// This version: 11. Jan 2019

package cpu

import (
	"testing"

	"angua/common"
)

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
		_ = test8.init(&c) // executes opcode

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
