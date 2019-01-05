// Test file for Angua Specials
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 05. Jan 2019
// This version: 05. Jan 2019

package specials

import (
	"testing"
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
