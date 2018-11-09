// Common files and type for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// This version: 07. Nov 2018
// First version: 09. Nov 2018

// This package contains base definitions and helper functions for all
// parts of Angua

package common

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	maxAddr = 1<<24 - 1

	// Commands from the CLI to the CPU
	HALT   = 0 // Stop the CPU
	RESUME = 1 // Continue processing from current PC address
	RUN    = 1 // Synonym for RESUME
	STEP   = 2 // Run one step
	STATUS = 3 // Print status information to stdout

	// Interrupts, resets, power toggle. These are sent by the CLI to the
	// CPU which then pretends that this condition just happened
	BOOT  = 4 // Power on
	RESET = 5 // Reset line
	IRQ   = 6 // Maskable interrupt
	NMI   = 7 // Non-maskable interrupt
	ABORT = 8 // Abort signal to chip

	// Reserved for future use
	VERBOSE = 9
	LACONIC = 10 // Turns verbose off
	TRACE   = 11 // Print out trace information
	NOTRACE = 12 // Turn off trace information
)

type Data8 uint8
type Data16 uint16

type Addr8 uint8 // For Direct Page
type Addr16 uint16
type Addr24 uint32 // We don't have a 24 bit type

// Ensure24 is a method of Addr24 that makes sure the upper byte of the
// underlying uint32 is actually zero
func (a Addr24) Ensure24() Addr24 {
	return a & 0x00FFFFFF
}

// Lsb takes a variant of int and returns the Least Significant Byte (LSB), that
// is, the lowest 8 bits, as a byte
func Lsb(n interface{}) byte {

	var r byte

	switch n := n.(type) {
	case uint8:
		r = n
	case Addr8:
		r = byte(n)
	case Addr16:
		r = byte(n & 0xFF)
	case Addr24:
		r = byte(n & 0xFF)
	case Data8:
		r = byte(n)
	case Data16:
		r = byte(n & 0xFF)
	default:
		log.Fatalf("ERROR: Lsb: Type %T doesn't have a Lsb", n)
	}

	return r
}

// Msb takes an int and returns the Most Significant Byte (MSB), that
// is, the second 8 bits as a byte.
func Msb(n interface{}) byte {
	var r byte

	switch n := n.(type) {
	case uint8:
		log.Fatalf("ERROR: Type %T doesn't have a MSB", n)
	case Addr8:
		log.Fatalf("ERROR: Type %T doesn't have a MSB", n)
	case Addr16:
		r = byte((n & 0xFF00) >> 8)
	case Addr24:
		r = byte((n & 0x00FF00) >> 8)
	case Data8:
		log.Fatalf("ERROR: Type %T doesn't have a MSB", n)
	case Data16:
		r = byte((n & 0xff00) >> 8)
	default:
		log.Fatalf("ERROR: Type %T doesn't have a MSB", n)
	}

	return r
}

// Bank takes an int and returns the Bank Byte (Bank), that
// is, the highest 8 bits as a byte.
func Bank(n interface{}) byte {
	var r byte

	switch n := n.(type) {
	case uint8, Addr8, Addr16, Data8, Data16:
		log.Fatalf("ERROR: Type %T doesn't have a bank byte", n)
	case Addr24:
		r = byte((n & 0xFF0000) >> 16)
	default:
		log.Fatalf("ERROR: Type %T doesn't have a bank byte", n)
	}

	return r
}

// convNum Converts a number string -- hex, binary, or decimal -- to an uint.
// We accept ':' and '.' as delimiters, use $ or 0x for hex numbers, % for
// binary numbers, and nothing for decimal numbers.
func ConvNum(s string) uint {

	ss := StripDelimiters(s)

	switch {

	case strings.HasPrefix(ss, "$"):
		n, err := strconv.ParseInt(ss[1:], 16, 0)
		if err != nil {
			log.Fatal("ERROR: Can't convert ", ss, " as hex number")
		}
		return uint(n)

	case strings.HasPrefix(ss, "0x"):
		n, err := strconv.ParseInt(ss[2:], 16, 0)
		if err != nil {
			log.Fatal("ERROR: Can't convert ", ss, " as hex number")
		}
		return uint(n)

	case strings.HasPrefix(ss, "%"):
		n, err := strconv.ParseInt(ss[1:], 2, 0)
		if err != nil {
			log.Fatal("ERROR: Can't convert ", ss, " as binary number")
		}
		return uint(n)

	default:
		n, err := strconv.ParseInt(ss, 10, 0)
		if err != nil {
			log.Fatal("ERROR: Can't convert ", ss, " as decimal number")
		}
		return uint(n)
	}
}

// fmtAddr takes a 65816 24 bit address as a uint and returns a hex number
// string with a ':' between the bank byte and the rest of the address. Hex
// digits are capitalized. Assumes we are sure that the address is valid
func FmtAddr(a Addr24) string {
	s1 := fmt.Sprintf("%06X", a)
	s2 := s1[0:2] + ":" + s1[2:len(s1)]
	return s2
}

// stripDelimiters removes '.' and ':' which users can use as number delimiters.
// Also removes spaces
func StripDelimiters(s string) string {
	s1 := strings.Replace(s, ":", "", -1)
	s2 := strings.Replace(s1, ".", "", -1)
	return strings.TrimSpace(s2)
}
