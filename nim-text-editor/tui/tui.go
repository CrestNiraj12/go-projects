package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
	fileHandler "nim-text-editor/file_handler"
)

const (
	startX = 5
)

type EditFile struct {
	cursorX, cursorY, scrollY, scrollX int
	content                            [][]rune
	isModified                         bool
	fileName                           string
}

func Init(filename string) {
	editFile := &EditFile{
		cursorX:  startX,
		fileName: filename,
	}

	contentByte, ok := fileHandler.ReadFile(filename)
	if !ok {
		return
	}

	for _, line := range SplitLines(string(contentByte)) {
		editFile.content = append(editFile.content, []rune(line))
	}

	editFile.displayContent()
	editFile.handleInput()
}

// n is line number, i is y position
// line number (n) increases with line while y position (i) is static on the screen
func insertLineNum(n int, i int) {
	lineFormat := fmt.Sprintf("%*d", startX-1, n+1)
	for j, r := range lineFormat {
		termbox.SetCell(j, i, r, termbox.ColorWhite, termbox.ColorDarkGray)
	}
}

func (ef *EditFile) displayContent() {
	if err := termbox.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing TUI: %v\n", err)
		return
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	defer termbox.Flush()

	if len(ef.content) == 0 {
		insertLineNum(0, 0)
		termbox.SetCell(0, 0, ' ', termbox.ColorWhite, termbox.ColorDarkGray)
		return
	} else {
		width, height := termbox.Size()

		for y := 0; y < height && (ef.scrollY+y) < len(ef.content); y++ {
			line := ef.content[ef.scrollY+y]
			insertLineNum(ef.scrollY+y, y)

			for x := 0; x+ef.scrollX < len(line) && x+startX < width; x++ {
				termbox.SetCell(x+startX, y, rune(line[x+ef.scrollX]), termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}

	termbox.SetCursor(ef.cursorX-ef.scrollX, ef.cursorY-ef.scrollY)
}

func displayMessage(message string) {
	_, height := termbox.Size()
	for _, line := range SplitLines(message) {
		for x, ch := range line {
			termbox.SetCell(x, height-1, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}

func SplitLines(content string) []string {
	return strings.Split(content, "\n")
}

func (ef *EditFile) Save() (bytes []byte, ok bool) {
	for _, line := range ef.content {
		bytes = append(bytes, []byte(strings.TrimSpace(string(line))+"\n")...)
	}

	ok = fileHandler.SaveFile(ef.fileName, &bytes)
	return
}

func (ef *EditFile) modifyWrapperFunc(handler interface{}, arg ...interface{}) {
	switch fn := handler.(type) {
	case func():
		fn()
	case func(rune):
		if len(arg) == 1 {
			if argValue, ok := arg[0].(rune); ok {
				fn(argValue)
			} else {
				fmt.Println("Error: Expected an int argument.")
				return
			}
		} else {
			fmt.Println("Error: One argument expected for this function.")
			return
		}
	}
	ef.isModified = true
}
