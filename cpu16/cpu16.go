// Angua CPU System - Native Mode CPU
// Scot W. Stevenson
// First version: 06. Nov 2018
// First version: 06. Nov 2018

package cpu16

const (
	maxAddr = 1<<24 - 1
)

type reg8 uint8
type reg16 uint16

type addr8 uint8
type addr16 uint16
type addr24 uint32

type CpuNative struct {
	A reg16
	X reg16
	Y reg16

	PC reg16
}
