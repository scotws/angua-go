// Test server for Angua
// Scot W. Stevenson
// First version: 26. Dec 2018
// This version: 26. Dec 2018

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
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening ...")
	doChat(conn)
	fmt.Println("... done")
}

func doChat(c net.Conn) {
	defer c.Close()
	r := bufio.NewReader(c)

	putChar(c, "Welcome to the Server\n")
	putChar(c, "Type 'quit' or use Ctrl-c to quit\n\n")

	for {
		n, msg := getChar(*r)
		ns := strconv.Itoa(n)

		fmt.Print("Got: ", msg)

		receipt := "You sent me " + ns + " bytes (including whitespace)\n"
		putChar(c, receipt)

		// getChar returns string with line feed
		if strings.TrimSpace(msg) == "quit" {
			putChar(c, "Closing connection.\n")
			break
		}
	}
}

// getChar gets a string from the client and returns it. String is not processed
// and includes line feeds etc
func getChar(r bufio.Reader) (int, string) {

	p := make([]byte, 256)

	n, err := r.Read(p)
	if err != nil {
		log.Fatal(err)
	}

	return n, string(p[:n])
}

// putChar takes a string and prints it at the remote connection. Note that the
// raw string is printed, without line feeds
func putChar(c net.Conn, s string) {
	_, err := io.WriteString(c, s)

	if err != nil {
		log.Print(err)
	}
}
