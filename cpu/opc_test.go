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

	c := CPU{}

	var tests = []struct {
		name string
		init func(*CPU)
		flag byte
		want byte
	}{
		{"0x18 (clc)", Opc18, c.FlagC, CLEAR},
	}

	for _, test := range tests {
		test.init(&c)
		if test.flag != test.want {
			t.Errorf("TestOpcFlags for %v returns %X, wanted %X", test.name, c.FlagC, test.want)
		}
	}
}
