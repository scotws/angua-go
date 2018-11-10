// CPU Switcher for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 10. Nov 2018
// This version: 10. Nov 2018

package switcher

import (
	"fmt"
	"time"

	"angua/cpu16"
	"angua/cpu8"
)

// heartBeat is a testing routine to print a string every 10 seconds to
// show that the go routine is working. It is called as a go routine from the
// switcher for testing
func heartBeat() {
	for {
		fmt.Println("<Switcher is alive>")
		time.Sleep(20 * time.Second)
	}
}

func Init(c8 *cpu8.Cpu8) {
	enable8 := make(chan<- struct{})
	enable8 <- struct{}{}
}

// Run instructs the Switcher to start the CPU at the given value in the PC This
// is the main Switcher loop that runs as a loop as a go routine
func Run(c8 *cpu8.Cpu8, c16 *cpu16.Cpu16) {

	fmt.Println("Switcher: DUMMY: Run")

	reqSwitch8 := make(<-chan struct{})
	reqSwitch16 := make(<-chan struct{})
	enable8 := make(chan<- struct{})
	enable16 := make(chan<- struct{})

	go heartBeat()

	select {

	case <-reqSwitch8:
		embiggen(c8, c16)
		enable16 <- struct{}{}

	case <-reqSwitch16:
		debiggen(c8, c16)
		enable8 <- struct{}{}
	}

}

// Embiggen changes the CPU from 8 bit (emulated) to 16 bit (native)
// This is triggered by a request from the cpu8 via a channel
func embiggen(c8 *cpu8.Cpu8, c16 *cpu16.Cpu16) {
	fmt.Println("Switcher: DUMMY: Embiggen")
}

// Debiggen changes the CPU from 16 bit (native) to 8 bit (emulated)
// This is triggered by a request from the cpu16 via a channel
func debiggen(c8 *cpu8.Cpu8, c16 *cpu16.Cpu16) {
	fmt.Println("Switcher: DUMMY: Debiggen")
}
