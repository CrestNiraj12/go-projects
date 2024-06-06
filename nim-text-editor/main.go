package main

import (
	"fmt"
	tui "nim-text-editor/tui"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run main.go <filename>\n")
		return
	}

	filename := os.Args[1]
  tui.NewTUI(filename).Init()
}
