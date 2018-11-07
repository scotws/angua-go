// Test file for Angua Common Files
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 15. Mar 2018
// This version: 07. Nov 2018

package common

import "testing"

func TestEnsure24(t *testing.T) {
	var tests = []struct {
		input Addr24
		want  Addr24
	}{
		{0xFFAABBCC, 0xAABBCC},
		{0x00AABBCC, 0xAABBCC},
	}

	for _, test := range tests {
		got := Addr24.Ensure24(test.input)

		if got != test.want {
			t.Errorf("Ensure24(%q) = %v", test.input, got)
		}
	}
}

// --- LSB ---

func TestLsbByte(t *testing.T) {
	var tests = []struct {
		input byte
		want  byte
	}{
		{0x00, 0},
		{0x80, 0x80},
		{0xFF, 0xFF},
	}

	for _, test := range tests {
		got := Lsb(test.input)

		if got != test.want {
			t.Errorf("Lsb(%q) = %v", test.input, got)
		}
	}
}

func TestLsbAddr24(t *testing.T) {
	var tests = []struct {
		input Addr24
		want  byte
	}{
		{0x000000, 0},
		{0xEE0080, 0x80},
		{0xFFFFFF, 0xFF},
	}

	for _, test := range tests {
		got := Lsb(test.input)

		if got != test.want {
			t.Errorf("Lsb(%q) = %v", test.input, got)
		}
	}
}

// --- MSB ---

func TestMsbAddr16(t *testing.T) {
	var tests = []struct {
		input Addr16
		want  byte
	}{
		{0x0000, 0},
		{0x8000, 0x80},
		{0xAABB, 0xAA},
	}

	for _, test := range tests {
		got := Msb(test.input)

		if got != test.want {
			t.Errorf("Msb(%q) = %v", test.input, got)
		}
	}
}

func TestMsbAddr24(t *testing.T) {
	var tests = []struct {
		input Addr24
		want  byte
	}{
		{0x00EE00, 0xEE},
		{0xEE0080, 0x00},
		{0xFFFFFF, 0xFF},
	}

	for _, test := range tests {
		got := Msb(test.input)

		if got != test.want {
			t.Errorf("Msb(%q) = %v", test.input, got)
		}
	}
}

func TestBankAddr24(t *testing.T) {
	var tests = []struct {
		input Addr24
		want  byte
	}{
		{0xAABBCC, 0xAA},
	}

	for _, test := range tests {
		got := Bank(test.input)

		if got != test.want {
			t.Errorf("Bank(%q) = %v", test.input, got)
		}
	}
}

// --- Other helpers ---

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
		got := ConvNum(test.input)

		if got != test.want {
			t.Errorf("convNum(%q) = %v", test.input, got)
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
		got := StripDelimiters(test.input)
		if got != test.want {
			t.Errorf("stripDelimiters(%q) = %v", test.input, got)
		}
	}
}

func TestFmtAddr(t *testing.T) {
	var tests = []struct {
		input Addr24
		want  string
	}{
		{0, "00:0000"},
		{1, "00:0001"},
		{1023, "00:03FF"},
		{0x01ffff, "01:FFFF"},
	}
	for _, test := range tests {
		got := FmtAddr(test.input)
		if got != test.want {
			t.Errorf("fmtAddr(%q) = %v", test.input, got)
		}
	}
}
