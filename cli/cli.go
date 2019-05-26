// Angua Command Line Interface (CLI)
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 24. May 2019
// This version: 26. May 2019

// The Command Line Interface (CLI) of Angua is started in another terminal
// window. To connect on Linux, use the nc program with the line
//
//	nc localhost 8000
//
// from another terminal. Quit with 'exit' in the shell.

// TODO Real error handling

package cli

import (
	// "bufio"
	"fmt"
	"io"
	"log"
	"net"

	// Use ishell as interface https://godoc.org/gopkg.in/abiosoft/ishell.v2
	"github.com/abiosoft/readline"
	"gopkg.in/abiosoft/ishell.v2"
)

const (
	banner string = `Welcome to Angua
An Emulatr for the 65816 in Native Mode
Version ALPHA 0.1  26. May 2019
Copyright (c) 2018-2019 Scot W. Stevenson
Angua comes with absolutely NO WARRANTY
Type 'help' for more information
`
	host   string = "localhost:8000"
	prompt string = "> "
)

// Start listens for a connection to the emulator and sets up the Command Line
// Interface (CLI) for Angua
func Start(cmd chan int, resp chan string, done chan struct{}) {

	// ----------------------------------------------------------
	// Create a listener for login from other terminal
	// See https://golang.org/pkg/net/#Listen
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Access accepted from ", conn.RemoteAddr().String())

	// ----------------------------------------------------------
	// Define Command Line Interface

	// ishell uses readline to set the configuration. We use it to redirect
	// the input and output of the shell to our external terminal. See
	// https://godoc.org/github.com/abiosoft/readline#Config for details
	termCfg := &readline.Config{
		Prompt:      prompt,
		Stdin:       io.ReadCloser(conn),
		Stdout:      io.Writer(conn),
		StdinWriter: io.Writer(conn),
		Stderr:      io.Writer(conn),
	}

	// Start interactive shell. Note that by default, this provides the
	// directives "exit", "help", and "clear".
	shell := ishell.NewWithConfig(termCfg)

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

	shell.Println(banner)
	shell.Run()

	// We reach this part when we use the 'exit' command
	shell.Close()

	// Signal emulator that we're all done
	fmt.Println("Closing shell.")
	done <- struct{}{}
}
