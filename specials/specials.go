// List of special routines for memory-mapped actions
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 05. Jan 2019
// This version: 15. Jan 2019

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

	"angua/common"
	"time"
)

// ---- READING ----

// GetChar returns a byte that comes from the user stitting at the interface and
// a bool to indicate success or failure. This routine should, but doesn't have
// to, always be present.
func GetChar() (byte, error) {
	// TODO DUMMY: Return ASCII character "a"
	return 0x61, nil
}

// GetCharBlocks returns a byte that comes from the user stitting at the
// interface and a bool to indicate success or failure. This routine should, but
// doesn't have to, always be present. This variant blocks until a character is
// received
func GetCharBlocks() (byte, error) {
	// TODO DUMMY: Return ASCII character "a"
	return 0x61, nil
}

// ---- WRITING ----

// PutChar takes a common.Data8 as present in an 8-bit register and writes it to
// the standard output. This routine should, but doesn't have to, always be
// present
func PutChar(c common.Data8) {
	// TODO DUMMY: Print as ASCII char to normal screen
	// TODO Move this to actual output
	fmt.Printf("%c", byte(c))
	return
}

// Sleep8 takes the value of A in 8 bit width and uses this as the number of
// seconds to sleep.
func Sleep8(a common.Data8) {
	sec := time.Duration(a)
	time.Sleep(sec * time.Second)
}
