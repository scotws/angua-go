// Test file for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 15. Mar 2018
// This version: 01. Jan 2019

package main

import (
	"testing"

	"angua/common"
)

func TestParseAddressRange(t *testing.T) {

	type result struct {
		addr1 common.Addr24
		addr2 common.Addr24
		ok    bool
	}

	var tests = []struct {
		input []string
		want  result
	}{
		{[]string{"too", "few"}, result{0, 0, false}},

		{[]string{"bank", "00"}, result{0x000000, 0x00FFFF, true}},
		{[]string{"bank", "$01"}, result{0x010000, 0x01FFFF, true}},
		{[]string{"bank", "0x0A"}, result{0x0A0000, 0x0AFFFF, true}},
		{[]string{"bank", "10"}, result{0x0A0000, 0x0AFFFF, true}},

		{[]string{"0000", "to", "1000"}, result{0, 1000, true}},
		{[]string{"00:0000", "to", "00:1000"}, result{0, 1000, true}},
		{[]string{"00:0000", "to", "$00:FFFF"}, result{0, 0xFFFF, true}},
		{[]string{"$FF:0000", "to", "$FF:FFFF"}, result{0xFF0000, 0xFFFFFF, true}},

		{[]string{"0x2000", "$8000"}, result{0x2000, 32768, true}},
		{[]string{"0x002000", "$00:8000"}, result{0x2000, 32768, true}},
		{[]string{"0x002000", "0x00:8000"}, result{0x2000, 32768, true}},
		{[]string{"0x002000", "00:8000"}, result{0x2000, 8000, true}},
		{[]string{"$FF:0000", "$FF:FFFF"}, result{0xFF0000, 0xFFFFFF, true}},
	}

	for _, test := range tests {
		got1, got2, ok := parseAddressRange(test.input)
		if got1 != test.want.addr1 ||
			got2 != test.want.addr2 ||
			ok != test.want.ok {
			t.Errorf("ParseAddressRange(%q) = %v", test.input, result{got1, got2, ok})
		}
	}
}
