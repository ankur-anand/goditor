package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	var charValue byte
	reader := bufio.NewReader(os.Stdin)
	// Each Iteration, reader reads a byte of data from the
	// source and assign it to the charValue, until there are
	// no more bytes to read.
	for {
		var err error
		// read one byte
		charValue, err = reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				fmt.Println("END OF FILE")
			}
		}
		fmt.Println(charValue)

	}
}
