package terminal

import (
	"log"

	"golang.org/x/sys/unix"
)

const (
	stdInFileNo  = 0
	stdoutFileNo = 1
)

// WinSize get the current window size of the terminal.
func WinSize() *unix.Winsize {
	// a process on unix can use the ioctl() TIOCGWINSZ operation to
	// find out the current size of the terminal window
	winsizeS, err := unix.IoctlGetWinsize(stdoutFileNo, unix.TIOCGWINSZ)
	if err != nil {
		log.Fatal(err)
	}

	return winsizeS
}

// Term is a Simple Terminal State Holder
type Term struct {
	cookedState *unix.Termios
}

// CurrentTerm returns a current terminal with it's state
// default In cooked mode data is preprocessed before being given to a program
func CurrentTerm() *Term {
	cooked, err := unix.IoctlGetTermios(stdInFileNo, unix.TCGETS)
	if err != nil {
		log.Fatalln(err)
	}

	t := &Term{
		cookedState: cooked,
	}

	return t
}

// EnableRawMode is used to put the terminal into raw mode.
// There are number of attributes of the terminal, but in its defaultt state
// called "cooked mode" input is line buffered, characters are automatically
// echoed, ^C - raise SIGINT, ^D signals EOF, and so on.
// In raw mode none of thos things happens; you get every keystroke
// immediately without waiting for a new line, and it's not echoed.
// ^C doesn't cause SIGINT, and so on.
func (t *Term) EnableRawMode() {

	state := *t.cookedState

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
	err := unix.IoctlSetTermios(stdInFileNo, unix.TCSETS, &state)
	if err != nil {
		log.Fatalln(err)
	}

}

// DisableRawMode Disable raw mode at exit
func (t *Term) DisableRawMode() {
	err := unix.IoctlSetTermios(stdInFileNo, unix.TCSETS, t.cookedState)
	if err != nil {
		log.Fatalln(err)
	}
}
