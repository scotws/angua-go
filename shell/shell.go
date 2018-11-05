// Test for ishell
// Scot W. Stevenson
// First version 30. June 2018
// This version 30. June 2018

// See https://github.com/abiosoft/ishell

package main

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

func main() {

	var fBatch = false

	shell := ishell.New()
	shell.Println("Sample shell")

	shell.AddCmd(&ishell.Cmd{
		Name: "hi",
		Help: "Say hello to user",
		Func: func(c *ishell.Context) {
			c.Println("Hello!")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "beep",
		Help: "Print beeping noise",
		Func: func(c *ishell.Context) {
			c.Println("\a")
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

	// Decide if echo should take parameter in quote
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

	// This doesn't work with the batch command
	shell.AddCmd(&ishell.Cmd{
		Name: "mode",
		Help: "Set mode of MPU (native or emulated)",
		Func: func(c *ishell.Context) {
			if !fBatch {
				choice := c.MultiChoice([]string{
					"native",
					"emulated",
				}, "Pick new mode")
				if choice == 0 {
					c.Println("(NATIVE MODE)")
				} else {
					c.Println("(EMULATED)")
				}

			} else {
				c.Println("ERROR: 'mode' not available in batch mode")
			}
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
