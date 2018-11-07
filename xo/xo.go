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

const (
	// Commands to the XO from Angua CLI via the cmd channel
	HALT   = 1
	STATUS = 2
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

	fmt.Println("XO: DUMMY: Init")
}

// heartBeat is a testing routine to print a string every 10 seconds to
// show that the go routine is working. It is called as a go routine from xo
// itself
func heartBeat() {
	for {
		fmt.Println("<XO is alive>")
		time.Sleep(10 * time.Second)
	}
}

// MakeItSo instructs the XO to start the CPU at the given value in the PC
// This is the main XO loop that runs as a loop as a go routine
func MakeItSo(c8 *cpu8.Cpu8, c16 *cpu16.Cpu16, in <-chan int) {

	fmt.Println("XO: DUMMY: MakeItSo")
	go heartBeat()

	// Run the Cpu8 first
	go cpuEmu.Run()

	// We wait for input from Angua over the in command channel
	for c := range in {

		switch c {
		case HALT:
			fmt.Println("XO: DUMMY: Received Halt cmd")
		case STATUS:
			fmt.Println("XO: DUMMY: Received Halt cmd")
		default:
			fmt.Println("XO: Unknown Signal:", c)
		}
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
