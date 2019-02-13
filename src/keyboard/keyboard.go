package keyboard

import (
	"bufio"
	"io"
	"os"
)

// all action related to keyboard

// Key just represent a kind of key pressed for action.
type Key int

const (
	// Quit Editor
	Quit           Key = 17
	escapeSeq      Key = '\x1b'
	escapeFollowed Key = '['
	// ArrowUp Cursor
	ArrowUp Key = iota + 1000
	// ArrowDown Cursor
	ArrowDown
	// ArrowRight Cursor
	ArrowRight
	// ArrowLeft Cursor
	ArrowLeft
	// PageUp Cursor
	PageUp
	// PageDown Cursor
	PageDown
	// DeleteKey Cursor
	DeleteKey
)

// cursorMap maps some of the control character
var cursorMap = map[byte]Key{
	'A': ArrowUp,
	'B': ArrowDown,
	'C': ArrowRight,
	'D': ArrowLeft,
}

// reader is used to read a byte from terminal
var reader = bufio.NewReader(os.Stdin)

// HandleKeyStroke handles each byte KeyStroke
func HandleKeyStroke() (Key, error) {
	return readByte(reader)
}

// to wait for one keypress, and return it
func readByte(reader io.ByteReader) (Key, error) {
	char, err := reader.ReadByte()
	if err != nil {
		return Key(char), err
	}
	// check if the char type that we have got is escapeSeq
	switch Key(char) {
	case escapeSeq:
		// wait for two other input.
		input1, _ := reader.ReadByte()
		if Key(input1) == escapeFollowed {
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
					return DeleteKey, nil
				case 5:
					return PageUp, nil
				case 6:
					return PageDown, nil

				}
			}
			return Key(input2), nil
		}

	}

	return Key(char), nil
}
