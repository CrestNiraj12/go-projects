package file_handler

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
)

// TODO
// Remember cursor position
// Implement better data structure - rope, tabulation, gap buffer, etc
// Undo / Redo
// Word wrap
// Refactor code
const (
	startX = 5
)

var (
	content          [][]rune
	cursorY          int
	cursorX          = startX
	fileName         string
	scrollY, scrollX int
	width, height    int
	isModified       bool
)

func InitHandler(filename string) {
	fileName = filename
	openAndReadFile()
	displayContent()
	handleInput()
}

func openAndReadFile() {
	contentByte, err := os.ReadFile(fileName)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		}
		return
	}

	for _, line := range strings.Split(string(contentByte), "\n") {
		content = append(content, []rune(line))
	}
}

// n is line number, i is y position
// line number (n) increases with line while y position (i) is static on the screen
func insertLineNum(n int, i int) {
	lineFormat := fmt.Sprintf("%*d", startX-1, n+1)
	for j, r := range lineFormat {
		termbox.SetCell(j, i, r, termbox.ColorWhite, termbox.ColorDarkGray)
	}
}

func displayContent() {
	if err := termbox.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing TUI: %v\n", err)
		return
	}
	width, height = termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	defer termbox.Flush()

	if len(content) == 0 {
		insertLineNum(0, 0)
		termbox.SetCell(0, 0, ' ', termbox.ColorWhite, termbox.ColorDarkGray)
	} else {
		width, height := termbox.Size()

		for y := 0; y < height && (scrollY+y) < len(content); y++ {
			line := content[scrollY+y]
			insertLineNum(scrollY+y, y)

			for x := 0; x+scrollX < len(line) && x+startX < width; x++ {
				termbox.SetCell(x+startX, y, rune(line[x+scrollX]), termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}

	termbox.SetCursor(cursorX-scrollX, cursorY-scrollY)
}

func displayMessage(message string) {
	for _, line := range splitLines(message) {
		for x, ch := range line {
			termbox.SetCell(x, height-1, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}

func splitLines(content string) []string {
	return strings.Split(content, "\n")
}

func saveFile() {
	var bytes []byte
	for _, line := range content {
		bytes = append(bytes, []byte(strings.TrimSpace(string(line))+"\n")...)
	}

	err := os.WriteFile(fileName, bytes, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		return
	}
	isModified = false
	message := fmt.Sprintf("Written %d bytes to file. Press any key to continue", len(bytes))
	displayMessage(message)
	termbox.PollEvent()
}

func closeFile() {
	if isModified {
		if doSave := promptSave(); doSave {
			saveFile()
		}
	}
}
