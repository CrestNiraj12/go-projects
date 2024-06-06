package file_handler

import (
	"errors"
	"fmt"
	"os"
)

// TODO
// Remember cursor position
// Implement better data structure - rope, tabulation, gap buffer, etc
// Undo / Redo
// Word wrap
// Refactor code

var (
	content    [][]rune
	isModified bool
)

func ReadFile(filename string) (content []byte, ok bool) {
	contentByte, err := os.ReadFile(filename)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return nil, false
	}
	return contentByte, true
}

func SaveFile(filename string, bytes *[]byte) (ok bool) {
	err := os.WriteFile(filename, *bytes, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		return false
	}
	return true
}
