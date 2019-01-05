// Test file for Angua CPU
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 05. Jan 2019
// This version: 05. Jan 2019

package cpu

import (
	"testing"
)

func TestGetStatReg(t *testing.T) {

	var tests = []struct {
		input StatReg
		want  byte
	}{
		{StatReg{FlagN: 1, FlagV: 1, FlagM: 1, FlagX: 1, FlagD: 1, FlagI: 1, FlagZ: 1, FlagC: 1}, 0xFF},
		{StatReg{FlagN: 0, FlagV: 0, FlagM: 0, FlagX: 0, FlagD: 0, FlagI: 0, FlagZ: 0, FlagC: 0}, 0x00},
		{StatReg{FlagN: 1, FlagV: 1, FlagM: 1, FlagX: 1, FlagD: 0, FlagI: 0, FlagZ: 0, FlagC: 0}, 0xF0},
		{StatReg{FlagN: 0, FlagV: 0, FlagM: 0, FlagX: 0, FlagD: 1, FlagI: 1, FlagZ: 1, FlagC: 1}, 0x0F},
		{StatReg{FlagN: 1, FlagV: 0, FlagM: 1, FlagX: 0, FlagD: 1, FlagI: 0, FlagZ: 1, FlagC: 0}, 0xAA},
		{StatReg{FlagN: 0, FlagV: 1, FlagM: 0, FlagX: 1, FlagD: 0, FlagI: 1, FlagZ: 0, FlagC: 1}, 0x55},
	}

	for _, test := range tests {
		got := test.input.GetStatReg()
		if got != test.want {
			t.Errorf("cpu.GetStatReg(%q) = %v", test.input, got)
		}
	}
}
