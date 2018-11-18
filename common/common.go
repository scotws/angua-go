// Common files and type for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// This version: 07. Nov 2018
// First version: 14. Nov 2018

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
	RUN    = 2 // Synonym for RESUME
	STEP   = 3 // Run one step
	STATUS = 4 // Print status information to stdout

	// Interrupts, resets, power toggle. These are sent by the CLI to the
	// CPU which then pretends that this condition just happened
	BOOT  = 5 // Power on
	RESET = 6 // Reset line
	IRQ   = 7 // Maskable interrupt
	NMI   = 8 // Non-maskable interrupt
	ABORT = 9 // Abort signal to chip

	// Reserved for future use
	VERBOSE = 10
	LACONIC = 11 // Turns verbose off
	TRACE   = 12 // Print out trace information
	NOTRACE = 13 // Turn off trace information
)

// ==== INTERFACES ===

type Lsber interface {
	Lsb() byte
}

type Msber interface {
	Msb() byte
}

type Banker interface {
	Bank() byte
}

// LilEnder is an interface for types that have a routine to convert numbers to
// little endian byte arrarys
type LilEnder interface {
	LilEnd() []byte
}

// ==== DATA TYPES ===

// We define special data types for addresses and data instead of using generic
// unit16 etc or even just int to make sure that the logic is right. At some
// point we might be forced to collapse Addr16 and Data16 together. Note that Go
// doesn't let us add methods to existing types such as byte or uint16

// --- Addr8 (byte) ---

type Addr8 uint8

func (a Addr8) Lsb() byte {
	return byte(a)
}

func (a Addr8) LilEnd() []byte {
	return []byte{byte(a)}
}

// HexString returns the value of the address as a byte in uppercase hex
// notation, but without a hex prefix such as "$" or "0x"
func (a Addr8) HexString() string {
	return fmt.Sprintf("%02X", uint8(a))
}

// --- Addr16 (double word) ---

type Addr16 uint16

func (a Addr16) Lsb() byte {
	return byte(a & 0x00FF)
}

func (a Addr16) Msb() byte {
	return byte((a & 0xFF00) >> 8)
}

func (a Addr16) LilEnd() []byte {
	return []byte{a.Lsb(), a.Msb()}
}

func (a Addr16) HexString() string {
	return fmt.Sprintf("%02X%02X", a.Msb(), a.Lsb())
}

// --- Addr24 (double word) ---

type Addr24 uint32 // We don't have a 24 bit type

func (a Addr24) Lsb() byte {
	return byte(a & 0x0000FF)
}

func (a Addr24) Msb() byte {
	return byte((a & 0x00FF00) >> 8)
}

func (a Addr24) Bank() byte {
	return byte((a & 0xFF0000) >> 16)
}

func (a Addr24) LilEnd() []byte {
	return []byte{a.Lsb(), a.Msb(), a.Bank()}
}

// HexString returns a string representation of the Addr24 address with the hex
// numbers converted to uppercase and a ":" delimiter between the bank byte and
// the rest of the address. There is no prefix such as "$" or "0x"
func (a Addr24) HexString() string {
	return fmt.Sprintf("%02X:%02X%02X", a.Bank(), a.Msb(), a.Lsb())
}

// Ensure24 is a method of Addr24 that makes sure the upper byte of the
// underlying uint32 is actually zero
func (a Addr24) Ensure24() Addr24 {
	return a & 0x00FFFFFF
}

// --- Data8 ---

// Data8 and Data16 are used for registers

type Data8 uint8

func (d Data8) Lsb() byte {
	return byte(d)
}

func (d Data8) LilEnd() []byte {
	return []byte{byte(d)}
}

func (d Data8) HexString() string {
	return fmt.Sprintf("%02X", uint8(d))
}

// --- Data16 ---

type Data16 uint16

func (d Data16) Lsb() byte {
	return byte(d & 0x00FF)
}

func (d Data16) Msb() byte {
	return byte((d & 0xFF00) >> 8)
}

func (d Data16) LilEnd() []byte {
	return []byte{d.Lsb(), d.Msb()}
}

func (d Data16) HexString() string {
	return fmt.Sprintf("%02X%02X", d.Msb(), d.Lsb())
}

// === General Helper Functions ===

// convNum Converts a number string -- hex, binary, or decimal -- to an uint.
// We accept ':' and '.' as delimiters, use $ or 0x for hex numbers, % for
// binary numbers, and nothing for decimal numbers. Note octal is not supported
// TODO We currently fail hard with a fatal log. When system is stable, replace
// by a different error scheme
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
// TODO See if we need this since we have the methods that do the same directly
// on the data
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
