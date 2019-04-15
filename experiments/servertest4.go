// https://tutorialedge.net/golang/reading-console-input-golang/

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(char)

	switch char {
	case 'A':
		fmt.Println("Key 'A' pressed")
	case 'a':
		fmt.Println("Key 'a' pressed")
	}
}
