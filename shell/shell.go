// Angua Interactive Shell
// Scot W. Stevenson
// First version 30. June 2018
// This version 06. Nov 2018

// See https://github.com/abiosoft/ishell

package shell

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/abiosoft/ishell.v2"
)

// procRange takes a list of arguments as a slice of strings an parses them
func procRange(s []string) (string, bool) {

	var r string
	var f bool

	if len(s) != 3 {
		return "Wrong number of arguments", false
	}

	if s[1] != "to" {
		return "Format must be <ADDR> 'to' <ADDR>", false
	}

	r = s[0] + "-" + s[2]
	f = true

	return r, f
}

func NewShell() {

	var fBatch = false

	// Create a new shell. By default, this includes the commands
	// "exit", "help", and "clear"

	shell := ishell.New()
	shell.Println("Angua Shell")

	shell.AddCmd(&ishell.Cmd{
		Name: "abort",
		Help: "Trigger the abort vector",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY trigger abort vector)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "beep",
		Help: "Print a beeping noise",
		Func: func(c *ishell.Context) {
			c.Println("\a")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "boot",
		Help: "Boot the machine. Same effect as turning on the power",
		Func: func(c *ishell.Context) {
			c.Println("(DUMMY boot the machine)")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "dump",
		Help: "Print hex dump of range",
		Func: func(c *ishell.Context) {
			r, ok := procRange(c.Args)
			if !ok {
				c.Print("ERROR: ", r, "\n")
			} else {
				c.Print("Got: ", r, "\n")
			}
		},
	})

	// TODO Decide if echo should take parameter in quotes
	shell.AddCmd(&ishell.Cmd{
		Name: "echo",
		Help: "Print following text to end of line",
		Func: func(c *ishell.Context) {
			c.Println(strings.Join(c.Args, " "))
		},
	})

	// To run with batch, just add "batch yellow frog"
	shell.AddCmd(&ishell.Cmd{
		Name: "yellow",
		Help: "Print word in yellow",
		Func: func(c *ishell.Context) {
			yellow := color.New(color.FgYellow).SprintFunc()
			c.Println(yellow(c.Args[0]))
		},
	})

	// Non-interactive execution
	if len(os.Args) > 1 && os.Args[1] == "batch" {
		fBatch = true
		shell.Process(os.Args[2:]...)
	} else {
		fBatch = false
		shell.Run()
		shell.Close()
	}
}
