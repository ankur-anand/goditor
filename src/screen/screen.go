package screen

import (
	"bytes"
	"fmt"
	"os"
)

// get the current window size of terminal.

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
