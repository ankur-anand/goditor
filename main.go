package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"

	"golang.org/x/sys/unix"
)

type key int

const (
	stdInFileNo        = 0
	stdoutFileNo       = 1
	goditorVersion     = 0.1
	quit           key = 17
	escapeSeq      key = '\x1b'
	escapeFollowed key = '['
	arrowUp        key = iota + 1000
	arrowDown
	arrowRight
	arrowLeft
	pageUp
	pageDown
	deleteKey
)

// cursorMap maps some of the control character
var cursorMap = map[byte]key{
	'A': arrowUp,
	'B': arrowDown,
	'C': arrowRight,
	'D': arrowLeft,
}

// storeGoditorRow for storing a row of text in our editor
type storeGoditorRow struct {
	text bytes.Buffer
}

// goditorStateT is to keep track of the cursor’s x and y position
// and Winsize
type goditorStateT struct {
	curX          uint16
	curY          uint16
	winsizeStruct *unix.Winsize
	row           []storeGoditorRow // slice of storeGoditorRow to store all row.
	numrows       int
	rowat         int // row of the file the user is currently scrolled to.
}

// global state.
var goditorState goditorStateT

// enableRawMode is used to put the terminal into raw mode.
// There are number of attributes of the terminal, but in its defaultt state
// called "cooked mode" input is line buffered, characters are automatically
// echoed, ^C - raise SIGINT, ^D signals EOF, and so on.
// In raw mode none of thos things happens; you get every keystroke
// immediately without waiting for a new line, and it's not echoed.
// ^C doesn't cause SIGINT, and so on.
func enableRawMode() (*unix.Termios, error) {

	cooked, err := unix.IoctlGetTermios(stdInFileNo, unix.TCGETS)
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

	// ISIG - To disable signal-generating characters (INTR, quit, SUSP)
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
	err = unix.IoctlSetTermios(stdInFileNo, unix.TCSETS, &state)
	if err != nil {
		return nil, err
	}

	return cooked, nil
}

// Disable raw mode at exit
func disableRawMode(cookedState *unix.Termios) error {
	err := unix.IoctlSetTermios(stdInFileNo, unix.TCSETS, cookedState)
	if err != nil {
		return err
	}
	return nil
}

// goditorMoveCursor move the cursor around
func goditorMoveCursor() {
	// CUP – Cursor Position
	// https://vt100.net/docs/vt100-ug/chapter3.html#CUP
	writeToTerminal(fmt.Sprintf("\x1b[%d;%dH", goditorState.curX, goditorState.curY))
}

// to wait for one keypress, and return it
func goditorReadKey(reader io.ByteReader) (key, error) {
	char, err := reader.ReadByte()
	if err != nil {
		return key(char), err
	}
	// check if the char type that we have got is escapeSeq
	switch key(char) {
	case escapeSeq:
		// wait for two other input.
		input1, _ := reader.ReadByte()
		if key(input1) == escapeFollowed {
			input2, _ := reader.ReadByte()
			// check if arrowLeft or arrowRight or
			// arrowUp or arrowDown key has been pressed.
			if val, ok := cursorMap[input2]; ok {
				return val, nil
			}

			// checkIf the PageUp or PageDown key has been
			// pressed
			// PageUp <esc>[5~
			// PageDown <esc>[6~
			if input2 >= 0 && input2 <= 9 {
				input3, _ := reader.ReadByte()
				switch input3 {
				case 3:
					return deleteKey, nil
				case 5:
					return pageUp, nil
				case 6:
					return pageDown, nil

				}
			}
			return key(input2), nil
		}

	}

	return key(char), nil
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
	//fmt.Println(goditorState, *goditorState.winsizeStruct)
	// mapping Ctrl + Q(17) to quit is
	switch char {
	case quit:
		return 1, nil
	case arrowUp:
		// Prevent moving the cursor values to go into the negatives
		if goditorState.curX != 1 {
			goditorState.curX--
		}
		goditorMoveCursor()
	case arrowDown:
		if goditorState.curX != goditorState.winsizeStruct.Row-1 {
			goditorState.curX++
		}
		goditorMoveCursor()
	case arrowLeft:
		if goditorState.curY != 1 {
			goditorState.curY--
		}
		goditorMoveCursor()
	case arrowRight:
		if goditorState.curY != goditorState.winsizeStruct.Col-1 {
			goditorState.curY++
		}
		goditorMoveCursor()
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
func goditorDrawRows() string {
	// get the size of the terminal
	// to know how many rows to draw
	editorConfig := *goditorState.winsizeStruct
	// Number of rows
	erow := editorConfig.Row
	crowat := goditorState.rowat

	var yaxis uint16
	// buffer is to avoid make a whole bunch of small
	// writeToTerminal every time of the loop.
	var buffer bytes.Buffer
	// ToDo: golang 1.10 has  strings.Builder type,
	// which achieves this even more efficiently

	for yaxis = 0; yaxis < erow; yaxis++ {
		editorrow := int(yaxis) + crowat
		// printing \r\n will cause the terminal to scroll
		// for a new blank line
		// so we should not printing the \r\n to last line
		// display the name of our editor and a version number too.
		if editorrow >= goditorState.numrows {
			// if numrows is zero, then only display the Goditor version
			// as in othercase we've opened a file
			if goditorState.numrows == 0 && yaxis == erow/2 {
				// and position it to the center of the screen
				message := "Goditor: v0.1"
				ecol := editorConfig.Col
				leftPad := (int(ecol) - len(message)) / 2
				if leftPad > 1 {
					buffer.WriteString("~")
					leftPad--
				}
				for {
					leftPad--
					if leftPad < 1 {
						break
					}
					buffer.WriteString(" ")
				}
				buffer.WriteString(message)
			} else {
				buffer.WriteString("~")
			}
		} else {

			text := goditorState.row[editorrow].text.String()
			buffer.WriteString(text)

		}
		// Erase In Line - the part of the line to the right of the cursor
		buffer.WriteString("\x1b[K")

		if yaxis < erow-1 {
			buffer.WriteString("\r\n")
		}

	}

	return buffer.String()
}

// Clear the screen
func goditorRefreshScreen() error {

	// write only once
	var buffer bytes.Buffer

	// terminal is drawning, cursor might be displayed.
	// hide the cursor for the flicker moments
	// https://vt100.net/docs/vt510-rm/DECTCEM.html
	// set Mode
	buffer.WriteString("\x1b[?25l")

	// CUP – Cursor Position
	// https://vt100.net/docs/vt100-ug/chapter3.html#CUP
	// position the cursor at the first row and first column,
	// not at the bottom.
	// get the position from cursorState
	buffer.WriteString(fmt.Sprintf("\x1b[%d;%dH", goditorState.curX, goditorState.curY))

	editor := goditorDrawRows()
	buffer.WriteString(editor)
	// reposition the cursor back up at the top-left corner.
	buffer.WriteString("\x1b[H")
	// cursor reset mode
	buffer.WriteString("\x1b[?25h")

	// Once write to Screen
	err := writeToTerminal(buffer.String())
	if err != nil {
		return err
	}
	return nil
}

// clearScreenOnExit clear the screen and
// reposition the cursor when program exits
func clearScreenOnExit() {
	var buffer bytes.Buffer
	buffer.WriteString("\x1b[2J")
	buffer.WriteString("\x1b[H")
	writeToTerminal(buffer.String())
}

// goditorOpen reads the contents from the file, and display
// the file content.
func goditorOpen(file *os.File) {

	// create a tmp slice
	rowvalue := make([]storeGoditorRow, 10)
	br := bufio.NewReader(file)
	// ReadLine does not include the line end ("\r\n" or "\n")
	i := 0
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		rowvalue[i].text.WriteString(string(line))
		i++
		goditorState.numrows++
	}

	// assign the goditor with the row value
	goditorState.row = rowvalue
}

func main() {
	// a process can use the ioctl() TIOCGWINSZ operation to
	// find out the current size of the terminal window
	winsizeS, err := unix.IoctlGetWinsize(stdoutFileNo, unix.TIOCGWINSZ)
	if err != nil {
		log.Fatal(err)
	}
	// Initialize the initial Cursor position.
	goditorState = goditorStateT{curX: 1, curY: 1, winsizeStruct: winsizeS, numrows: 0, rowat: 0}
	// if args is the fileName open the file content
	args := os.Args
	if len(args) == 2 {
		// open the file, for reading.
		file, err := os.Open("test.txt")
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		goditorOpen(file)
	}
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
