// Executive Officer (XO) for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 07. Nov 2018
// This version: 07. Nov 2018

package xo

import (
	"fmt"
	"time"

	"angua/cpu16"
	"angua/cpu8"
)

var (
	cpuEmu *cpu8.Cpu8
	cpuNat *cpu16.Cpu16
)

// Init takes pointers to a cpu8 a cpu16 from the CLI and gets them ready
// to run
func Init(c8 *cpu8.Cpu8, c16 *cpu16.Cpu16) {
	cpuEmu = c8
	cpuNat = c16
}

// MakeItSo instructs the XO to start the CPU at the given value in the PC
// This is the main XO loop that runs as a loop as a go routine
func MakeItSo(c8 *cpu8.Cpu8, c16 *cpu16.Cpu16) {
	fmt.Println("XO: DUMMY: MakeItSo")

	// Run the Cpu8 first
	go cpuEmu.Run()

	// TODO This is test routine
	for {
		fmt.Println("<XO is on the job!>")
		time.Sleep(10 * time.Second)
	}
}

// Embiggen changes the CPU from 8 bit (emulated) to 16 bit (native)
// This is triggered by a request from the cpu8 via a channel
func Embiggen(c8 *cpu8.Cpu8, c16 *cpu16.Cpu16) {
	fmt.Println("XO: DUMMY: Embiggen")
}

// Debiggen changes the CPU from 16 bit (native) to 8 bit (emulated)
// This is triggered by a request from the cpu16 via a channel
func Debiggen(c8 *cpu8.Cpu8, c16 *cpu16.Cpu16) {
	fmt.Println("XO: DUMMY: Debiggen")
}
