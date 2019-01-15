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

// TODO This is a temporary test
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

	for _, test := range tests {
		fmt.Printf("Starting 16 in %d seconds ...", test.duration)
		Sleep16(common.Data16(test.duration))
		fmt.Println(" done.")
	}
}
