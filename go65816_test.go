// Test file for go65816
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 15. Mar 2018
// This version: 21. Mar 2018

package main

import "testing"

func TestConvNumber(t *testing.T) {
	var tests = []struct {
		input string
		want  uint
	}{
		{"0", 0},
		{"100", 100},
		{"$400", 1024},
		{"0x400", 1024},
		{"%10", 2},
		{"%00001111", 0x0F},
		{"0x00:0400", 1024},
	}

	for _, test := range tests {
		got := convNum(test.input)

		if got != test.want {
			t.Errorf("convNum(%q) = %v", test.input, got)
		}
	}
}

func TestIsValidAddr(t *testing.T) {
	var tests = []struct {
		input uint
		want  bool
	}{
		{0, true},
		{1 << 24, false},
		{1<<24 - 1, true},
	}

	for _, test := range tests {
		got := isValidAddr(test.input)
		if got != test.want {
			t.Errorf("isValidAddr(%q) = %v", test.input, got)
		}
	}
}

func TestStripDelimiters(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"0000", "0000"},
		{"1000", "1000"},
		{"1000 ", "1000"},
		{"00:0000", "000000"},
		{"00.0000", "000000"},
		{"0x00:FFFF", "0x00FFFF"},
		{"0x00.FFFF", "0x00FFFF"},
		{"$00:FFFF", "$00FFFF"},
		{" $00:FFFF", "$00FFFF"},
	}

	for _, test := range tests {
		got := stripDelimiters(test.input)
		if got != test.want {
			t.Errorf("stripDelimiters(%q) = %v", test.input, got)
		}
	}
}

func TestFmtAddr(t *testing.T) {
	var tests = []struct {
		input uint
		want  string
	}{
		{0, "00:0000"},
		{1, "00:0001"},
		{1023, "00:03FF"},
		{0x01ffff, "01:FFFF"},
	}
	for _, test := range tests {
		got := fmtAddr(test.input)
		if got != test.want {
			t.Errorf("fmtAddr(%q) = %v", test.input, got)
		}
	}
}
