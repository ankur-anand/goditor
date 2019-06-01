package screen

import (
	"bytes"
	"fmt"
	"os"

	"github.com/ankur-anand/goditor/src/keyboard"
	"github.com/ankur-anand/goditor/src/terminal"
	"golang.org/x/sys/unix"
)

const welcomeMessage = "Goditor -- version 0.1"

// get the current window size of terminal.

// State is to keep track of the cursor’s x and y position
// and Winsize
type State struct {
	curX       uint16
	curY       uint16
	winsize    *unix.Winsize
	screenrows int
	screencols int
}

// NewState initialize a new goditor state which manages the screen of editor.
func NewState() *State {
	// we initialize the both cursor at position
	// at X:1 and Y:1 where X is row and Y is Column.
	s := &State{curX: 1, curY: 1, winsize: terminal.WinSize()}
	s.screenrows = int(s.winsize.Row)
	s.screencols = int(s.winsize.Col)
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
func (s *State) ClearScreen() {
	var buf bytes.Buffer
	buf.WriteString("\x1b[?25l") // vt set mode
	// \x1b (27) Escape Character.
	// three bytes are [2J
	// J command (Erase In Display)
	// https://vt100.net/docs/vt100-ug/chapter3.html#ED
	// buf.WriteString("\x1b[2J") instead use erase in line now
	//WriteToStdOut("\x1b[H")
	s.drawTildesRow(&buf)
	// https://vt100.net/docs/vt100-ug/chapter3.html#CUP
	// CUP – Cursor Position
	buf.WriteString("\x1b[H")
	buf.WriteString("\x1b[?25h") // vt reset mode
	WriteToStdOut(buf.String())
}

// func (s *State) ClearScreenOnExit() {
// 	WriteToStdOut("\x1b[J")
// 	WriteToStdOut("\x1b[H")
// }

// moves the cursor around the screen
func (s *State) moveCursor() {
	// CUP – Cursor Position
	// https://vt100.net/docs/vt100-ug/chapter3.html#CUP
	WriteToStdOut(fmt.Sprintf("\x1b[%d;%dH", s.curX, s.curY))
}

// ProcessKey waits for a keypress recieved, and then handles it
func (s *State) ProcessKey(key keyboard.Key) {

	switch key {

	case keyboard.ArrowUp:
		// Prevent moving the cursor values to go into the negatives
		if s.curX != 1 {
			s.curX--
		}
		s.moveCursor()
	case keyboard.ArrowDown:
		//if s.curY
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

// draw a tildes on the screen as in vim editor
func (s *State) drawTildesRow(buf *bytes.Buffer) {
	for colY := 0; colY < s.screenrows; colY++ {
		buf.WriteString("~")

		if colY == s.screenrows/3 {
			// add padding to welcome message
			pad := (s.screencols - len(welcomeMessage)) / 2
			for pad > 0 {
				buf.WriteString(" ")
				pad--
			}
			buf.WriteString(welcomeMessage)
		}
		buf.WriteString("\x1b[K") // erase in line
		if colY < s.screenrows-1 {
			buf.WriteString("\r\n")
		}
	}
}

// insertCharAtScreen inserts a single character into an row
func insertCharAtScreen(char keyboard.Key) {
	WriteToStdOut(string(char))
}

func insertChar(char keyboard.Key) {
	insertCharAtScreen(char)
}
