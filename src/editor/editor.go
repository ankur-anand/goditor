package editor

import (
	"fmt"

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

// StartGoditor starts a new Editor
func StartGoditor() {
	s := newState()
	term := terminal.CurrentTerm()
	term.EnableRawMode()
	fmt.Println(s)
}
