package main

import (
	"fmt"
	tui "nim-text-editor/tui"
	"os"
)
// TODO
// Fix cursor and space issue
// Remember cursor position
// Implement better data structure - rope, tabulation, gap buffer, etc
// Undo / Redo
// Word wrap
// Refactor code

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run main.go <filename>\n")
		return
	}

	filename := os.Args[1]
  tui.NewTUI(filename).Init()
}
