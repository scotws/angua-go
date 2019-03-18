// Test server version 2 for Angua
// Scot W. Stevenson
// First version: 26. Dec 2018
// This version: 17. Mar 2019

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
	"strconv"
	"strings"
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

func doChat(c net.Conn) {
	defer c.Close()
	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)

	putString(c, "Welcome to the Server\n")
	putString(c, "Type 'quit' or use Ctrl-c to quit\n\n")

	for {
		n, msg := getString(*r)
		ns := strconv.Itoa(n)

		fmt.Print("Got: ", msg)

		receipt := "You sent me " + ns + " bytes (including whitespace)\n"
		putString(c, receipt)

		err := putByte(0x61, w)
		if err != nil {
			fmt.Println("Error in putByte:", err)
		}

		// getString returns string with line feed
		if strings.TrimSpace(msg) == "quit" {
			putString(c, "Closing connection.\n")
			break
		}
	}
}

// getString gets a string from the client and returns it. String is not processed
// and includes line feeds etc
func getString(r bufio.Reader) (int, string) {

	p := make([]byte, 256)

	n, err := r.Read(p)
	if err != nil {
		log.Fatal(err)
	}

	return n, string(p[:n])
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
