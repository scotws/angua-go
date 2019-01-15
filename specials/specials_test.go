// Test file for Angua Specials
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 05. Jan 2019
// This version: 05. Jan 2019

package specials

import (
	"fmt"
	"testing"

	"angua/common"
)

func TestGetChar(t *testing.T) {

	var tests = []struct {
		want byte
	}{
		{0x61}, // TODO TESTING ASCII hex for "a"
	}

	for _, test := range tests {
		got, _ := GetChar()
		if got != test.want {
			t.Errorf("GetChar(%q) = 0x%X", test.want, got)
		}
	}
}

func TestPutChar(t *testing.T) {

	var tests = []struct {
		char common.Data8
	}{
		{0x61}, // ASCII hex for "a"
		{0x62}, // ASCII hex for "a"
		{0x63}, // ASCII hex for "a"
	}

	for _, test := range tests {
		fmt.Printf("Printing '%c': ", test.char)
		PutChar(test.char)
		fmt.Println()
	}
}

// Test of sleep function
func TestSleep(t *testing.T) {

	var tests = []struct {
		duration common.Data8
	}{
		{1},
		{5},
	}

	for _, test := range tests {
		fmt.Printf("Starting 8 in %d seconds ...", test.duration)
		Sleep8(test.duration)
		fmt.Println(" done.")
	}
}
