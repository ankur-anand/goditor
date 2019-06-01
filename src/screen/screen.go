package screen

import (
	"bytes"
	"fmt"
	"os"

	"github.com/ankur-anand/goditor/src/keyboard"
	"github.com/ankur-anand/goditor/src/terminal"
	"golang.org/x/sys/unix"
)

// get the current window size of terminal.

// State is to keep track of the cursor’s x and y position
// and Winsize
type State struct {
	curX    uint16
	curY    uint16
	winsize *unix.Winsize
}

// NewState initialize a new goditor state which manages the screen of editor.
func NewState() *State {
	// we initialize the both cursor at position
	// at X:1 and Y:1 where X is row and Y is Column.
	s := &State{curX: 1, curY: 1, winsize: terminal.WinSize()}
	return s
}

// WriteToStdOut writes n bytes to the os.Stdout file
func WriteToStdOut(value string) error {
	STDOUTFILE := os.Stdout
	_, err := fmt.Fprintf(STDOUTFILE, value)
	return err
}

// ClearScreen clear the screen and
// reposition the cursor when program exits
func ClearScreen() {
	var buffer bytes.Buffer
	// \x1b (27) Escape Character.
	// three bytes are [2J
	// J command (Erase In Display)
	// https://vt100.net/docs/vt100-ug/chapter3.html#ED
	buffer.WriteString("\x1b[2J")
	// https://vt100.net/docs/vt100-ug/chapter3.html#CUP
	// CUP – Cursor Position
	buffer.WriteString("\x1b[H")
	WriteToStdOut(buffer.String())
}

// ProcessKey waits for a keypress recieved, and then handles it
func (s *State) ProcessKey(key keyboard.Key) {

	switch key {

	case keyboard.ArrowUp:
		// Prevent moving the cursor values to go into the negatives
		// if goditorState.curRow != 1 {
		// 	goditorState.curRow--
		// }
		// goditorMoveCursor()
	case keyboard.ArrowDown:
		// if goditorState.curRow != goditorState.winsizeStruct.Row-1 {
		// 	goditorState.curRow++
		// }
		// goditorMoveCursor()
	case keyboard.ArrowLeft:
		// if goditorState.curCol != 1 {
		// 	goditorState.curCol--
		// }
		// goditorMoveCursor()
	case keyboard.ArrowRight:
		// if goditorState.curCol != goditorState.winsizeStruct.Col-1 {
		// 	goditorState.curCol++
		// }
		// goditorMoveCursor()
	default:
		insertChar(key)
	}

}

// insertCharAtScreen inserts a single character into an row
func insertCharAtScreen(char keyboard.Key) {
	WriteToStdOut(string(char))
}

func insertChar(char keyboard.Key) {
	insertCharAtScreen(char)
}
