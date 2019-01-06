// List of special routines for memory-mapped actions
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 05. Jan 2019
// This version: 05. Jan 2019

/* To add an action to a certain memory address, create the function here
   either as read or write action. Use the line formats

   	reading from <ADDR> calls <FUNC>
	writing to <ADDR> calls <FUNC>

   in the configuration file to connect the address to the action. The code
   should always contain GetChar and PutChar.
*/

package specials

import (
	"fmt"
)

// GetChar returns a byte that comes from the user stitting at the interface and
// a bool to indicate success or failure. This routine should, but doesn't have
// to, always be present.
func GetChar() (byte, bool) {
	// TODO DUMMY: Return ASCII character "a"
	// TODO move this to actual input
	return 0x61, true
}

// GetCharBlocks returns a byte that comes from the user stitting at the
// interface and a bool to indicate success or failure. This routine should, but
// doesn't have to, always be present. This variant blocks until a character is
// received
func GetCharBlocks() (byte, bool) {
	// TODO DUMMY: Return ASCII character "a"
	// TODO move this to actual input
	return 0x61, true
}

// PutChar takes a byte and writes it to the standard output. This routine
// should, but doesn't have to, always be present
func PutChar(b byte) {
	// TODO DUMMY: Print as ASCII char to normal screen
	// TODO Move this to actual output
	fmt.Printf("%c", b)
	return
}
