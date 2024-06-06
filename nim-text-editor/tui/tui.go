package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
	constants "nim-text-editor/constants"
	editFile "nim-text-editor/edit_file"
	fileHandler "nim-text-editor/file_handler"
)

type TUI struct {
	ef     *editFile.EditFile
	startX int
}

func NewTUI(filename string) *TUI {
	return &TUI{
		ef: &editFile.EditFile{
			Cursor:   &editFile.FileCursor{CursorX: constants.StartX},
			FileName: filename,
		},
		startX: constants.StartX,
	}
}

func (tui *TUI) Init() {
	ef := tui.ef
	contentByte, ok := fileHandler.ReadFile(ef.FileName)
	if !ok {
		return
	}

	for _, line := range SplitLines(string(contentByte)) {
		ef.Content = append(ef.Content, []rune(line))
	}

	tui.displayContent()
	tui.handleInput()
}

// n is line number, i is y position
// line number (n) increases with line while y position (i) is static on the screen
func insertLineNum(n int, i int) {
	lineFormat := fmt.Sprintf("%*d", constants.StartX-1, n+1)
	for j, r := range lineFormat {
		termbox.SetCell(j, i, r, termbox.ColorWhite, termbox.ColorDarkGray)
	}
}

func (tui *TUI) displayContent() {
	ef := tui.ef
	if err := termbox.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing TUI: %v\n", err)
		return
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	defer termbox.Flush()

	cur := *(ef.Cursor)
	if len(ef.Content) == 0 {
		insertLineNum(0, 0)
		termbox.SetCell(0, 0, ' ', termbox.ColorWhite, termbox.ColorDarkGray)
		return
	} else {
		width, height := termbox.Size()

		for y := 0; y < height && (cur.ScrollY+y) < len(ef.Content); y++ {
			line := ef.Content[cur.ScrollY+y]
			insertLineNum(cur.ScrollY+y, y)

			for x := 0; x+cur.ScrollX < len(line) && x+constants.StartX < width; x++ {
				termbox.SetCell(x+constants.StartX, y, rune(line[x+cur.ScrollX]), termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}

	termbox.SetCursor(cur.CursorX-cur.ScrollX, cur.CursorY-cur.ScrollY)
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

func (tui *TUI) Save() (bytes []byte, ok bool) {
	for _, line := range tui.ef.Content {
		bytes = append(bytes, []byte(strings.TrimSpace(string(line))+"\n")...)
	}

	ok = fileHandler.SaveFile(tui.ef.FileName, &bytes)
	return
}

func (tui *TUI) modifyWrapperFunc(handler interface{}, arg ...interface{}) {
	switch fn := handler.(type) {
	case func():
		fn()
	case func(rune):
		if len(arg) == 1 {
			if argValue, ok := arg[0].(rune); ok {
				fn(argValue)
			} else {
				fmt.Println("Error: Expected an rune argument.")
				return
			}
		} else {
			fmt.Println("Error: One argument expected for this function.")
			return
		}
  case func(int):
		if len(arg) == 1 {
			if argValue, ok := arg[0].(int); ok {
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
	tui.ef.IsModified = true
}
