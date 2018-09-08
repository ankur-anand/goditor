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

	// IXON (enable start/stop output control)
	// Disable control characters that Ctrl-S and Ctrl-Q produce

	// ICRNL - Input is gathered into lines, terminated by one of the line-delimiter characters:
	// NL, EOL, EOL2 (if the IEXTEN flag is set), EOF (at anything other than the initial
	// position in the line), or CR (if the ICRNL flag is enabled).
	// Default setting ^M
	//  with the ICRNL (map CR to NL on input) flag set
	// (the default), this character is first converted to a newline (ASCII 10 decimal, ^J)
	// before being passed to the reading proces

	// BRKINT - Disable Signal interrupt (SIGINT) on BREAK condition
	state.Iflag &^= (unix.IXON | unix.ICRNL | unix.BRKINT)

	//  OPOST - to disables output postprocessing "\n" to "\r\n" and

	state.Oflag &^= (unix.OPOST)

	// The ECHO feature causes each key you type to be printed to the terminal,
	// so you can see what you’re typing

	// ICANON flag  to turn off canonical mode, so that input can be read, byte
	// by byte instead of line-by-line.
	// ln line-by-line mode you have to press Enter for program to read input.

	// ISIG - To disable signal-generating characters (INTR, QUIT, SUSP)
	// Eg disable Ctrl-C(SIGINT) and Ctrl-Z(SIGTSTP) signals

	// IEXTERN - Disable Ctrl-V

	state.Lflag &^= (unix.ECHO | unix.ICANON | unix.ISIG | unix.IEXTEN)

	// set timeout read
	// VMIN - the minimum number of bytes of input needed before read,
	//  Character-at-a-time input
	state.Cc[unix.VMIN] = 1
	// VTIME -the maximum amount of time to wait before read returns,
	// It is in tenths of a second, so we set it to 1/10 of a second,
	// or 100 milliseconds
	// times out, it will return 0
	state.Cc[unix.VTIME] = 1

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

// to wait for one keypress, and return it
func goditorReadKey(reader io.ByteReader) (byte, error) {
	char, err := reader.ReadByte()
	if err != nil {
		return char, err
	}
	return char, nil
}

// goditorActionKeypress waits for a keypress, and then handles it
func goditorActionKeypress(reader io.ByteReader) (int, error) {
	char, err := goditorReadKey(reader)
	if err != nil {
		return 0, err
	}

	if unicode.IsControl(rune(char)) == true {
		fmt.Printf("%b\r\n", char)
	} else {
		fmt.Printf("%s\r\n", string(char))
	}

	// mapping Ctrl + Q(17) to quit is
	switch char {
	case 17:
		return 1, nil
	}

	return 0, nil
}

// writeToTerminal writes n bytes to the os.Stdout file
func writeToTerminal(value string) error {
	STDOUTFILE := os.Stdout
	_, err := fmt.Fprintf(STDOUTFILE, value)
	return err
}

// goditorDrawRows() draws a tilde in each rowof the buffer.
// and each row is not part of the file.
func goditorDrawRows() {
	for i := 0; i < 24; i++ {
		writeToTerminal("~\r\n")
	}
}

// Clear the screen
func goditorRefreshScreen() error {
	// ED – Erase In Display
	// https://vt100.net/docs/vt100-ug/chapter3.html#ED
	err := writeToTerminal("\x1b[2J")
	if err != nil {
		return err
	}

	// CUP – Cursor Position
	// https://vt100.net/docs/vt100-ug/chapter3.html#CUP
	// position the cursor at the first row and first column,
	// not at the bottom.
	err = writeToTerminal("\x1b[H")
	if err != nil {
		return err
	}
	goditorDrawRows()

	// eposition the cursor back up at the top-left corner.
	err = writeToTerminal("\x1b[H")
	if err != nil {
		return err
	}
	return nil
}

// clearScreenOnExit clear the screen and
// reposition the cursor when program exits
func clearScreenOnExit() {
	writeToTerminal("\x1b[2J")
	writeToTerminal("\x1b[H")
}

func main() {
	cookedState, err := enableRawMode()
	if err != nil {
		clearScreenOnExit()
		log.Fatalln(err)
	}

	reader := bufio.NewReader(os.Stdin)
	goditorRefreshScreen()
	// Each Iteration, reader reads a byte of data from the
	// source and assign it to the charValue, until there are
	// no more bytes to read.
	for {

		read, err := goditorActionKeypress(reader)
		if err != nil {
			disableRawMode(cookedState)
			clearScreenOnExit()
			log.Fatalln(err)
		}

		// if read is 1 that means we need to end the loop.
		if read == 1 {
			clearScreenOnExit()
			err := disableRawMode(cookedState)
			if err != nil {
				log.Fatalln(err)
			}
			os.Exit(0)
		}

	}
}
