// Test server version 2 for Angua
// Scot W. Stevenson
// First version: 26. Dec 2018
// This version: 14. April 2019

// Starts a little server that listens for a connection on localhost:8000 and
// returns information on what is typed. To connect, use the nc program

//		nc localhost 8000

// from another terminal. Quit with 'quit' or Ctrl-c (which sends EOF)

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

const host string = "localhost:8000" // Typed constant

func main() {
	// Create a listener, see https://golang.org/pkg/net/#Listen
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Access accepted from ", conn.RemoteAddr().String())
	// Not as a go routine because it would end everything immediately
	doChat(conn)
	fmt.Println("... done")
}

// doChat is the top-level loop for handling input and output
func doChat(c net.Conn) {
	defer c.Close()
	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)

	putString(c, "Welcome to the Server\n")
	putString(c, "Type 'q' or use Ctrl-c to quit\n\n")

	for {
		rc := getChar(*r)
		char := string(rc)

		fmt.Print("Got: ", char)

		receipt := "You sent me " + char + " \n"
		putString(c, receipt)

		err := putByte(0x61, w)
		if err != nil {
			fmt.Println("Error in putByte:", err)
		}

		if char == "q" {
			putString(c, "Closing connection.\n")
			break
		}
	}
}

// See https://tutorialedge.net/golang/reading-console-input-golang/
func getChar(r bufio.Reader) rune {

	char, _, err := r.ReadRune()
	if err != nil {
		log.Fatal(err)
	}

	return char
}

// putByte takes a byte and prints it to the given writer
func putByte(b byte, w *bufio.Writer) error {
	var errSum error = nil

	err1 := w.WriteByte(b)
	err2 := w.Flush()

	if err1 != nil || err2 != nil {
		errSum = fmt.Errorf("putByte Error: %v %v", err1, err2)
	}

	return errSum
}

// putString takes a string and prints it at the remote connection. Note that the
// raw string is printed, without line feeds
func putString(c net.Conn, s string) {
	_, err := io.WriteString(c, s)

	if err != nil {
		log.Print(err)
	}
}
