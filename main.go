package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"

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
func enableRawMode() (*unix.Termios, error) {

	cooked, err := unix.IoctlGetTermios(STDIN_FILENO, unix.TCGETS)
	if err != nil {
		return nil, err
	}
	state := *cooked
	// The ECHO feature causes each key you type to be printed to the terminal,
	// so you can see what you’re typing

	// ICANON flag  to turn off canonical mode, so that input can be read, byte
	// by byte instead of line-by-line.
	// ln line-by-line mode you have to press Enter for program to read input.

	// ISIG - To disable signal-generating characters (INTR, QUIT, SUSP)
	// Eg disable Ctrl-C(SIGINT) and Ctrl-Z(SIGTSTP) signals

	// IEXTERN - Disable Ctrl-V

	state.Lflag &^= (unix.ECHO | unix.ICANON | unix.ISIG | unix.IEXTEN)

	// IXON (enable start/stop output control)
	// Disable control characters that Ctrl-S and Ctrl-Q produce

	// ICRNL - Input is gathered into lines, terminated by one of the line-delimiter characters:
	// NL, EOL, EOL2 (if the IEXTEN flag is set), EOF (at anything other than the initial
	// position in the line), or CR (if the ICRNL flag is enabled).
	// Default setting ^M
	//  with the ICRNL (map CR to NL on input) flag set
	// (the default), this character is first converted to a newline (ASCII 10 decimal, ^J)
	// before being passed to the reading proces

	state.Iflag &^= (unix.IXON | unix.ICRNL)

	//TCSANOW is TCSETS
	err = unix.IoctlSetTermios(STDIN_FILENO, unix.TCSETS, &state)
	if err != nil {
		return nil, err
	}

	return cooked, nil
}

// Disable raw mode at exit
func disableRawMode(cookedState *unix.Termios) error {
	err := unix.IoctlSetTermios(STDIN_FILENO, unix.TCSETS, cookedState)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	cookedState, err := enableRawMode()
	if err != nil {
		log.Fatalln(err)
	}
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
			log.Fatalln(err)
		}

		if unicode.IsControl(rune(charValue)) == true {
			fmt.Println(charValue)
		} else {
			fmt.Println(string(charValue))
		}
		// press q to quit.
		if charValue == 'q' {
			// disable the RAW MODE
			err := disableRawMode(cookedState)
			if err != nil {
				log.Fatalln(err)
			}
			os.Exit(0)
		}

	}
}
