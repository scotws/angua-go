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

var (
	reqSwitchTo8  = make(chan struct{})
	reqSwitchTo16 = make(chan struct{})
	enable8       = make(chan struct{})
	enable16      = make(chan struct{})
)

// heartBeat is a testing routine to print a string every 10 seconds to
// show that the go routine is working. It is called as a go routine from the
// switcher for testing
func heartBeat() {
	for {
		fmt.Println("<Switcher is alive>")
		time.Sleep(2 * time.Minute)
	}
}

func Run(c8 *cpu8.Cpu8) {
	fmt.Println("Switcher: DUMMY: Run")
	enable8 <- struct{}{}
	return
}

// Run instructs the Switcher to start the CPU at the given value in the PC This
// is the main Switcher loop that runs as a loop as a go routine
// TODO add a verbose flag to log if switcher is alive and waiting
func Daemon(c8 *cpu8.Cpu8, c16 *cpu16.Cpu16, cmd <-chan int) {

	fmt.Println("Switcher: DUMMY: Daemon running")

	go heartBeat()

	go c8.Run(cmd, enable8, reqSwitchTo16)
	go c16.Run(cmd, enable16, reqSwitchTo8)

	for {
		select {

		case <-reqSwitchTo16:
			fmt.Println("Switcher: DUMMY: cpu8 requests switch to Native Mode")
			goEmulated(c8, c16)
			enable16 <- struct{}{}

		case <-reqSwitchTo8:
			fmt.Println("Switcher: DUMMY: cpu16 requests switch to Emulated Mode")
			goNative(c8, c16)
			enable8 <- struct{}{}
		}
	}

}

// goNative changes the CPU from 8 bit (emulated) to 16 bit (native)
// This is triggered by a request from the cpu8 via a channel
func goNative(c8 *cpu8.Cpu8, c16 *cpu16.Cpu16) {
	fmt.Println("Switcher: DUMMY: goEmulated, cpu16 -> cpu8")
}

// goEmulated changes the CPU from 16 bit (native) to 8 bit (emulated)
// This is triggered by a request from the cpu16 via a channel
func goEmulated(c8 *cpu8.Cpu8, c16 *cpu16.Cpu16) {
	fmt.Println("Switcher: DUMMY: goEmulated, cpu8 -> cpu16")
}
