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

func TestPutStatReg(t *testing.T) {

	var gotStatReg StatReg

	var tests = []struct {
		input byte
		want  StatReg
	}{
		{0xFF, StatReg{FlagN: 1, FlagV: 1, FlagM: 1, FlagX: 1, FlagD: 1, FlagI: 1, FlagZ: 1, FlagC: 1}},
		{0x00, StatReg{FlagN: 0, FlagV: 0, FlagM: 0, FlagX: 0, FlagD: 0, FlagI: 0, FlagZ: 0, FlagC: 0}},
		{0xF0, StatReg{FlagN: 1, FlagV: 1, FlagM: 1, FlagX: 1, FlagD: 0, FlagI: 0, FlagZ: 0, FlagC: 0}},
		{0x0F, StatReg{FlagN: 0, FlagV: 0, FlagM: 0, FlagX: 0, FlagD: 1, FlagI: 1, FlagZ: 1, FlagC: 1}},
		{0xAA, StatReg{FlagN: 1, FlagV: 0, FlagM: 1, FlagX: 0, FlagD: 1, FlagI: 0, FlagZ: 1, FlagC: 0}},
		{0x55, StatReg{FlagN: 0, FlagV: 1, FlagM: 0, FlagX: 1, FlagD: 0, FlagI: 1, FlagZ: 0, FlagC: 1}},

		{0x01, StatReg{FlagN: 0, FlagV: 0, FlagM: 0, FlagX: 0, FlagD: 0, FlagI: 0, FlagZ: 0, FlagC: 1}},
		{0x02, StatReg{FlagN: 0, FlagV: 0, FlagM: 0, FlagX: 0, FlagD: 0, FlagI: 0, FlagZ: 1, FlagC: 0}},
		{0x04, StatReg{FlagN: 0, FlagV: 0, FlagM: 0, FlagX: 0, FlagD: 0, FlagI: 1, FlagZ: 0, FlagC: 0}},
		{0x08, StatReg{FlagN: 0, FlagV: 0, FlagM: 0, FlagX: 0, FlagD: 1, FlagI: 0, FlagZ: 0, FlagC: 0}},
		{0x10, StatReg{FlagN: 0, FlagV: 0, FlagM: 0, FlagX: 1, FlagD: 0, FlagI: 0, FlagZ: 0, FlagC: 0}},
		{0x20, StatReg{FlagN: 0, FlagV: 0, FlagM: 1, FlagX: 0, FlagD: 0, FlagI: 0, FlagZ: 0, FlagC: 0}},
		{0x40, StatReg{FlagN: 0, FlagV: 1, FlagM: 0, FlagX: 0, FlagD: 0, FlagI: 0, FlagZ: 0, FlagC: 0}},
		{0x80, StatReg{FlagN: 1, FlagV: 0, FlagM: 0, FlagX: 0, FlagD: 0, FlagI: 0, FlagZ: 0, FlagC: 0}},
	}

	for _, test := range tests {
		gotStatReg.SetStatReg(test.input)
		if gotStatReg != test.want {
			t.Errorf("cpu.SetStatReg(%q) = %v", test.input, gotStatReg)
		}
	}
}
