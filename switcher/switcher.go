// CPU Switcher for Angua
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 10. Nov 2018
// This version: 10. Nov 2018

package switcher

import (
	"fmt"
	"time"

	"angua/emulated"
	"angua/native"
)

var (
	reqSwitchToEmu = make(chan struct{})
	reqSwitchToNat = make(chan struct{})
	enableEmu      = make(chan struct{})
	enableNat      = make(chan struct{})
)

// heartBeat is a testing routine to print a string every 10 seconds to
// show that the go routine is working. It is called as a go routine from the
// switcher for testing
func heartBeat() {
	for {
		fmt.Println("<Switcher is alive>")
		time.Sleep(5 * time.Minute)
	}
}

func Run(cEmu *emulated.Emulated) {
	fmt.Println("Switcher: DUMMY: Run")
	enableEmu <- struct{}{}
	return
}

// Run instructs the Switcher to start the CPU at the given value in the PC This
// is the main Switcher loop that runs as a loop as a go routine
// TODO add a verbose flag to log if switcher is alive and waiting
func Daemon(cEmu *emulated.Emulated, cNat *native.Native, cmd <-chan int) {

	fmt.Println("Switcher: DUMMY: Daemon running")

	go heartBeat()

	go cEmu.Run(cmd, enableEmu, reqSwitchToNat)
	go cNat.Run(cmd, enableNat, reqSwitchToEmu)

	for {
		select {

		case <-reqSwitchToNat:
			fmt.Println("Switcher: DUMMY: emulated requests switch to Native Mode")
			goEmulated(cEmu, cNat)
			enableNat <- struct{}{}

		case <-reqSwitchToEmu:
			fmt.Println("Switcher: DUMMY: native requests switch to Emulated Mode")
			goNative(cEmu, cNat)
			enableEmu <- struct{}{}
		}
	}

}

// goNative changes the CPU from 8 bit (emulated) to 16 bit (native)
// This is triggered by a request from the emulated via a channel
func goNative(cEmu *emulated.Emulated, cNat *native.Native) {
	fmt.Println("Switcher: DUMMY: goEmulated, native -> emulated")
}

// goEmulated changes the CPU from 16 bit (native) to 8 bit (emulated)
// This is triggered by a request from the native via a channel
func goEmulated(cEmu *emulated.Emulated, cNat *native.Native) {
	fmt.Println("Switcher: DUMMY: goEmulated, emulated -> native")
}
