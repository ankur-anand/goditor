package editor

import (
	"log"
	"os"

	"github.com/ankur-anand/goditor/src/keyboard"
	"github.com/ankur-anand/goditor/src/screen"
	"github.com/ankur-anand/goditor/src/terminal"
)

const (
	stdInFileNo    = 0
	stdoutFileNo   = 1
	goditorVersion = 0.1
)

type editor struct {
	term   *terminal.Term
	screen *screen.State
}

// StartGoditor starts a new Editor
func StartGoditor() {
	e := editor{
		term:   terminal.CurrentTerm(),
		screen: screen.NewState(),
	}

	// enable raw mode for the terminal
	e.term.EnableRawMode()
	for {
		e.screen.ClearScreen()
		readValue := e.keyboardRead()
		e.screen.ProcessKey(readValue)
	}
}

// keyboardRead reads a byte of data from the keyboard
func (e editor) keyboardRead() keyboard.Key {
	keyedIn, err := keyboard.HandleKeyStroke()
	if err != nil {
		e.term.DisableRawMode()
		e.screen.ClearScreen()
		log.Fatalln(err)
	}

	// quit Editor case
	if keyedIn == keyboard.Quit {
		e.term.DisableRawMode()
		e.screen.ClearScreen()
		os.Exit(0)
	}
	return keyedIn
}