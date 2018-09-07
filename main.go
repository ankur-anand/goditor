package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/unix"
)

const (
	STDIN_FILENO = 0
)

// enableRawMode is used to put the terminal into raw mode.
// There are number of attributes of the terminal, but in its defaultt state
// called "cooked mode" input is line buffered, characters are automatically
// echoed, ^C - raise SIGINT, ^D signals EOF, and so on.
// In raw mode none of thos things happens; you get every keystroke
// immediately without waiting for a new line, and it's not echoed.
// ^C doesn't cause SIGINT, and so on.
func enableRawMode() *unix.Termios {

	cooked, err := unix.IoctlGetTermios(STDIN_FILENO, unix.TCGETS)
	if err != nil {
		panic(err)
	}
	state := *cooked
	// The ECHO feature causes each key you type to be printed to the terminal,
	// so you can see what you’re typing
	state.Lflag &^= unix.ECHO

	//TCSANOW is TCSETS
	err = unix.IoctlSetTermios(STDIN_FILENO, unix.TCSETS, &state)
	if err != nil {
		panic(err)
	}

	return cooked
}

// Disable raw mode at exit
func disableRawMode(cookedState *unix.Termios) {
	err := unix.IoctlSetTermios(STDIN_FILENO, unix.TCSETS, cookedState)
	if err != nil {
		panic(err)
	}
}

func main() {
	cookedState := enableRawMode()

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
		// press q to quit.
		if charValue == 'q' {
			// disable the RAW MODE
			disableRawMode(cookedState)
			os.Exit(0)
		}

	}
}
