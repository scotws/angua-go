// Common files and type for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// This version: 07. Nov 2018
// First version: 06. Jan 2019

// This package contains base definitions and helper functions for all
// parts of Angua

package common

import (
	"fmt"
	"math/rand"
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
	//      5 // CURRENTLY UNUSED, was BOOT ("Power on")
	RESET = 6 // Reset line
	IRQ   = 7 // Maskable interrupt
	NMI   = 8 // Non-maskable interrupt
	ABORT = 9 // Abort signal to chip

	// Reserved for future use
	VERBOSE   = 10
	NOVERBOSE = 11 // Turns verbose off
	TRACE     = 12 // Print out trace information
	NOTRACE   = 13 // Turn off trace information

	// Vectors' addresses relevant for the 65816 in native mode
	AbortAddr Addr24 = 0xFFE8
	BRKAddr   Addr24 = 0xFFE6
	COPAddr   Addr24 = 0xFFE4
	IRQAddr   Addr24 = 0xFFEE
	NMIAddr   Addr24 = 0xFFEA
	ResetAddr Addr24 = 0xFFFC
)

// Interrupt vectors. Note the reset vector is only for emulated mode.
// See http://6502.org/tutorials/65c816interrupts.html and Eyes & Lichty
// p. 195 for details. We store these as 24 bit addresses because that
// is the way we'll use them during mem.FetchMore().
type Vector struct {
	Addr Addr24
	Name string
}

var (
	Vectors = []Vector{
		{AbortAddr, "Abort"},
		{BRKAddr, "BRK"},
		{COPAddr, "COP"},
		{IRQAddr, "IRQ"},
		{NMIAddr, "NMI"},
		{ResetAddr, "Reset"},
	}
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

// Ensure24 takes a 24 bit address and makes sure that the upper byte of the
// underlying uint32 is actually zero
func Ensure24(a Addr24) Addr24 {
	return a & 0x00FFFFFF
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

// HexString returns a string representation of the Addr16 address with the hex
// numbers converted to uppercase. There is no prefix such as "$" or "0x"
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

// ConvertNum Converts a number string -- hex, binary, or decimal -- to an uint.
// We accept ':' and '.' as delimiters, use $ or 0x for hex numbers, % for
// binary numbers, and nothing for decimal numbers. Note octal is not supported
func ConvertNum(s string) (uint, error) {
	var n int64
	var err error

	ss := StripDelimiters(s)

	switch {

	case strings.HasPrefix(ss, "$"):
		n, err = strconv.ParseInt(ss[1:], 16, 0)
	case strings.HasPrefix(ss, "0x"):
		n, err = strconv.ParseInt(ss[2:], 16, 0)
	case strings.HasPrefix(ss, "%"):
		n, err = strconv.ParseInt(ss[1:], 2, 0)
	default:
		n, err = strconv.ParseInt(ss, 10, 0)
	}

	if err != nil {
		return 0, fmt.Errorf("ConvertNum: Couldn't convert %s: %v", s, err)
	}

	return uint(n), nil
}

// stripDelimiters removes '.' and ':' which users can use as number delimiters.
// Also removes spaces
func StripDelimiters(s string) string {
	s1 := strings.Replace(s, ":", "", -1)
	s2 := strings.Replace(s1, ".", "", -1)
	return strings.TrimSpace(s2)
}

// UndefinedByte returns a randomized byte. It is used for situations where the
// state of part of the CPU is undefinded so the user doesn't learn to expect
// (say) a 00. See cpu.reset() for an example.
func UndefinedByte() byte {
	return byte(rand.Intn(255))
}

// UndefinedBit returns either 0 or 1. It is used for situations where the
// state of a flag is undefinded so the user doesn't learn to expect
// a certain flag, for instance after a RESET (see cpu.reset()).
func UndefinedBit() byte {
	return byte(rand.Intn(1))
}
