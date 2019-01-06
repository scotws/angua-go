// Test file for Angua CPU opcodes
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 05. Jan 2019
// This version: 05. Jan 2019

package cpu

import (
	"testing"
)

// clc
func TestOpcFlags(t *testing.T) {
	var c = CPU{}
	var tests = []struct {
		name string
		init func(*CPU) error
		flag *byte
		want byte
	}{
		{"0x18 (clc)", Opc18, &c.FlagC, CLEAR},
		{"0x38 (sec)", Opc38, &c.FlagC, SET},
		{"0x58 (cli)", Opc58, &c.FlagI, CLEAR},
		{"0x78 (sei)", Opc78, &c.FlagI, SET},
		{"0xD8 (cld)", OpcD8, &c.FlagD, CLEAR},
		{"0xB8 (clv)", OpcB8, &c.FlagV, CLEAR},
	}

	for _, test := range tests {
		test.init(&c)

		if *test.flag != test.want {
			t.Errorf("TestOpcFlags for %v returns %X, wanted %X",
				test.name, *test.flag, test.want)
		}

	}
}
