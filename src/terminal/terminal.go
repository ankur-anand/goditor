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
	// a process can use the ioctl() TIOCGWINSZ operation to
	// find out the current size of the terminal window
	winsizeS, err := unix.IoctlGetWinsize(stdoutFileNo, unix.TIOCGWINSZ)
	if err != nil {
		log.Fatal(err)
	}

	return winsizeS
}
