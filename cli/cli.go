// Angua Command Line Interface (CLI)
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 24. May 2019
// This version: 25. May 2019

package cli

import (
	"fmt"

	"gopkg.in/abiosoft/ishell.v2" // https://godoc.org/gopkg.in/abiosoft/ishell.v2
)

const (
	shellBanner string = `Welcome to Angua
An Emulator for the 65816 in Native Mode
Version ALPHA 0.1  24. May 2019
Copyright (c) 2018-2019 Scot W. Stevenson
Angua comes with absolutely NO WARRANTY
Type 'help' for more information
`
)

// Start creates the actual CLI shell

func Start(cmd chan int, resp chan string, done chan struct{}) {

	// Start interactive shell. Note that by default, this provides the
	// directives "exit", "help", and "clear".
	shell := ishell.New()
	shell.Println(shellBanner)
	shell.SetPrompt("> ")

	// We create a history file. This currently only works for Linux
	// TODO point this out in the documentation
	shell.SetHomeHistoryPath(".angua_shell_history")

	// Individual shell comands. Normal short help is lower case with no
	// punctuation
	shell.AddCmd(&ishell.Cmd{
		Name: "beep",
		Help: "make a beeping noise",
		Func: func(c *ishell.Context) {
			c.Println("\a")
		},
	})

	// TODO Testing
	shell.AddCmd(&ishell.Cmd{
		Name: "command",
		Help: "send a command int",
		Func: func(*ishell.Context) {
			cmd <- 1
		},
	})

	shell.Run()

	// We reach this part when we use the 'exit' command
	shell.Close()

	// Signal emulator that we're all done
	fmt.Println("Closing shell.")
	done <- struct{}{}
}
