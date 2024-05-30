package main

import (
	"fmt"
	filehandler "nim-text-editor/file_handler"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run main.go <filename>\n")
		return
	}

	filename := os.Args[1]
	fmt.Printf("Opening file: %s...\n", filename)

	filehandler.InitHandler(filename)
}
