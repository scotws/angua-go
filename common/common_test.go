// Test file for Angua Common Files
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 15. Mar 2018
// This version: 01. Jan 2019

package common

import (
	"bytes"
	"testing"
)

// --- Tests for Addr8 ---

func TestAddr8Lsb(t *testing.T) {
	var tests = []struct {
		input Addr8
		want  byte
	}{
		{0x00, 0x00},
		{0x01, 0x01},
		{0xAA, 0xAA},
	}

	for _, test := range tests {
		got := test.input.Lsb()

		if got != test.want {
			t.Errorf("Addr8.Lsb(%q) = %v", test.input, got)
		}
	}
}

func TestAddr8LilEnd(t *testing.T) {
	var tests = []struct {
		input Addr8
		want  []byte
	}{
		{0x00, []byte{0x00}},
		{0x01, []byte{0x01}},
		{0xAA, []byte{0xAA}},
	}

	for _, test := range tests {
		got := test.input.LilEnd()

		// Can't do normal comparisons with byte slices
		if !bytes.Equal(got, test.want) {
			t.Errorf("Addr8.LilEnd(%q) = %v", test.input, got)
		}
	}
}

func TestAddr8HexString(t *testing.T) {
	var tests = []struct {
		input Addr8
		want  string
	}{
		{0x00, "00"},
		{0x01, "01"},
		{0xAA, "AA"},
	}

	for _, test := range tests {
		got := test.input.HexString()

		if got != test.want {
			t.Errorf("Addr8.HexString(%q) = %v", test.input, got)
		}
	}
}

// --- Tests for Addr16 ---

func TestAddr16Lsb(t *testing.T) {
	var tests = []struct {
		input Addr16
		want  byte
	}{
		{0x0000, 0x00},
		{0x0201, 0x01},
		{0xFFAA, 0xAA},
	}

	for _, test := range tests {
		got := test.input.Lsb()

		if got != test.want {
			t.Errorf("Addr16.Lsb(%q) = %v", test.input, got)
		}
	}
}

func TestAddr16Msb(t *testing.T) {
	var tests = []struct {
		input Addr16
		want  byte
	}{
		{0x0000, 0x00},
		{0x0201, 0x02},
		{0xFFAA, 0xFF},
	}

	for _, test := range tests {
		got := test.input.Msb()

		if got != test.want {
			t.Errorf("Addr16.Msb(%q) = %v", test.input, got)
		}
	}
}

func TestAddr16LilEnd(t *testing.T) {
	var tests = []struct {
		input Addr16
		want  []byte
	}{
		{0x0000, []byte{0x00, 0x00}},
		{0x0201, []byte{0x01, 0x02}},
		{0xAAFF, []byte{0xFF, 0xAA}},
	}

	for _, test := range tests {
		got := test.input.LilEnd()

		// Can't do normal comparisons with byte slices
		if !bytes.Equal(got, test.want) {
			t.Errorf("Addr16.LilEnd(%q) = %v", test.input, got)
		}
	}
}

func TestAddr16HexString(t *testing.T) {
	var tests = []struct {
		input Addr16
		want  string
	}{
		{0x0000, "0000"},
		{0x0102, "0102"},
		{0xAAFF, "AAFF"},
	}

	for _, test := range tests {
		got := test.input.HexString()

		if got != test.want {
			t.Errorf("Addr16.HexString(%q) = %v", test.input, got)
		}
	}
}

// --- Tests for Addr24 ---

func TestAddr24Lsb(t *testing.T) {
	var tests = []struct {
		input Addr24
		want  byte
	}{
		{0x000000, 0x00},
		{0x030201, 0x01},
		{0xAABBCC, 0xCC},
	}

	for _, test := range tests {
		got := test.input.Lsb()

		if got != test.want {
			t.Errorf("Addr24.Lsb(%q) = %v", test.input, got)
		}
	}
}

func TestAddr24Msb(t *testing.T) {
	var tests = []struct {
		input Addr24
		want  byte
	}{
		{0x000000, 0x00},
		{0x030201, 0x02},
		{0xBBFFAA, 0xFF},
	}

	for _, test := range tests {
		got := test.input.Msb()

		if got != test.want {
			t.Errorf("Addr24.Msb(%q) = %v", test.input, got)
		}
	}
}

func TestAddr24Bank(t *testing.T) {
	var tests = []struct {
		input Addr24
		want  byte
	}{
		{0x000000, 0x00},
		{0x030201, 0x03},
		{0xBBFFAA, 0xBB},
	}

	for _, test := range tests {
		got := test.input.Bank()

		if got != test.want {
			t.Errorf("Addr24.Bank(%q) = %v", test.input, got)
		}
	}
}

func TestAddr24LilEnd(t *testing.T) {
	var tests = []struct {
		input Addr24
		want  []byte
	}{
		{0x000000, []byte{0x00, 0x00, 0x00}},
		{0x030201, []byte{0x01, 0x02, 0x03}},
		{0xAABBCC, []byte{0xCC, 0xBB, 0xAA}},
	}

	for _, test := range tests {
		got := test.input.LilEnd()

		// Can't do normal comparisons with byte slices
		if !bytes.Equal(got, test.want) {
			t.Errorf("Addr24.LilEnd(%q) = %v", test.input, got)
		}
	}
}

func TestAddr24HexString(t *testing.T) {
	var tests = []struct {
		input Addr24
		want  string
	}{
		{0x000000, "00:0000"},
		{0x010203, "01:0203"},
		{0xAABBCC, "AA:BBCC"},
	}

	for _, test := range tests {
		got := test.input.HexString()

		if got != test.want {
			t.Errorf("Addr24.HexString(%q) = %v", test.input, got)
		}
	}
}

func TestEnsure24(t *testing.T) {
	var tests = []struct {
		input Addr24
		want  Addr24
	}{
		{0x00000000, 0x000000},
		{0x04030201, 0x030201},
		{0xAABBCCDD, 0x00BBCCDD},
	}

	for _, test := range tests {
		got := Ensure24(test.input)

		if got != test.want {
			t.Errorf("Ensure24(%q) = %v", test.input, got)
		}
	}
}

// --- Tests for Data8 ---

func TestData8Lsb(t *testing.T) {
	var tests = []struct {
		input Data8
		want  byte
	}{
		{0x00, 0x00},
		{0x01, 0x01},
		{0xAA, 0xAA},
	}

	for _, test := range tests {
		got := test.input.Lsb()

		if got != test.want {
			t.Errorf("Data8.Lsb(%q) = %v", test.input, got)
		}
	}
}

func TestData8LilEnd(t *testing.T) {
	var tests = []struct {
		input Data8
		want  []byte
	}{
		{0x00, []byte{0x00}},
		{0x01, []byte{0x01}},
		{0xAA, []byte{0xAA}},
	}

	for _, test := range tests {
		got := test.input.LilEnd()

		// Can't do normal comparisons with byte slices
		if !bytes.Equal(got, test.want) {
			t.Errorf("Data8.LilEnd(%q) = %v", test.input, got)
		}
	}
}

func TestData8HexString(t *testing.T) {
	var tests = []struct {
		input Data8
		want  string
	}{
		{0x00, "00"},
		{0x01, "01"},
		{0xAA, "AA"},
	}

	for _, test := range tests {
		got := test.input.HexString()

		if got != test.want {
			t.Errorf("Data8.HexString(%q) = %v", test.input, got)
		}
	}
}

// --- Tests for Data16 ---

func TestData16Lsb(t *testing.T) {
	var tests = []struct {
		input Data16
		want  byte
	}{
		{0x0000, 0x00},
		{0x0201, 0x01},
		{0xFFAA, 0xAA},
	}

	for _, test := range tests {
		got := test.input.Lsb()

		if got != test.want {
			t.Errorf("Data16.Lsb(%q) = %v", test.input, got)
		}
	}
}

func TestData16Msb(t *testing.T) {
	var tests = []struct {
		input Data16
		want  byte
	}{
		{0x0000, 0x00},
		{0x0201, 0x02},
		{0xFFAA, 0xFF},
	}

	for _, test := range tests {
		got := test.input.Msb()

		if got != test.want {
			t.Errorf("Data16.Msb(%q) = %v", test.input, got)
		}
	}
}

func TestData16LilEnd(t *testing.T) {
	var tests = []struct {
		input Data16
		want  []byte
	}{
		{0x0000, []byte{0x00, 0x00}},
		{0x0201, []byte{0x01, 0x02}},
		{0xAAFF, []byte{0xFF, 0xAA}},
	}

	for _, test := range tests {
		got := test.input.LilEnd()

		// Can't do normal comparisons with byte slices
		if !bytes.Equal(got, test.want) {
			t.Errorf("Data16.LilEnd(%q) = %v", test.input, got)
		}
	}
}

func TestData16HexString(t *testing.T) {
	var tests = []struct {
		input Data16
		want  string
	}{
		{0x0000, "0000"},
		{0x0102, "0102"},
		{0xAAFF, "AAFF"},
	}

	for _, test := range tests {
		got := test.input.HexString()

		if got != test.want {
			t.Errorf("Data16.HexString(%q) = %v", test.input, got)
		}
	}
}

// --- Other helpers ---

func TestConvertNumber(t *testing.T) {
	type conv struct {
		numb uint
		ok   bool
	}

	var tests = []struct {
		input string
		want  conv
	}{
		{"0", conv{0, true}},
		{"100", conv{100, true}},
		{"$400", conv{1024, true}},
		{"0x400", conv{1024, true}},
		{"%10", conv{2, true}},
		{"%00001111", conv{0x0F, true}},
		{"%0000:1111", conv{0x0F, true}},
		{"%0000.1111", conv{0x0F, true}},
		{"0x00:0400", conv{1024, true}},
		{"00::0400", conv{400, true}}, // gracefully handle typos

		{"foobar", conv{0, false}},
		{"&0001", conv{0, false}},
	}

	for _, test := range tests {
		got, ok := ConvNum(test.input)
		res := conv{got, ok}
		if res != test.want {
			t.Errorf("convNum(%q) = %v", test.input, res)
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
