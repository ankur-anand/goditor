package editor

import (
	"log"
	"os"

	"github.com/ankur-anand/goditor/src/keyboard"
	"github.com/ankur-anand/goditor/src/screen"
	"github.com/ankur-anand/goditor/src/terminal"
	"golang.org/x/sys/unix"
)

const (
	stdInFileNo    = 0
	stdoutFileNo   = 1
	goditorVersion = 0.1
)

// state is to keep track of the cursor’s x and y position
// and Winsize
type state struct {
	curX    uint16
	curY    uint16
	winsize *unix.Winsize
}

// initialize a new goditor state.
func newState() state {
	// we initialize the both cursor at position
	// at X:1 and Y:1 where X is row and Y is Column.
	s := state{curX: 1, curY: 1, winsize: terminal.WinSize()}
	return s
}

type editor struct {
	term *terminal.Term
}

// StartGoditor starts a new Editor
func StartGoditor() {
	//s := newState()
	e := editor{
		term: terminal.CurrentTerm(),
	}
	// enable raw mode for the terminal
	e.term.EnableRawMode()
	for {
		readValue := e.keyboardRead()
		screen.ProcessKey(readValue)
	}
}

// keyboardRead reads a byte of data from the keyboard
func (e editor) keyboardRead() keyboard.Key {
	keyedIn, err := keyboard.HandleKeyStroke()
	if err != nil {
		e.term.DisableRawMode()
		screen.ClearScreen()
		log.Fatalln(err)
	}

	// quit Editor case
	if keyedIn == keyboard.Quit {
		e.term.DisableRawMode()
		screen.ClearScreen()
		os.Exit(0)
	}
	return keyedIn
}
