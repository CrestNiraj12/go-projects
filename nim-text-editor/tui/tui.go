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
	ef                    *editFile.EditFile
	startX, width, height int
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

  ef.Content = &editFile.ContentTable{}
	content := []rune(string(contentByte))
	ef.Content.Original = &content
	ef.Content.Pieces = append(ef.Content.Pieces, &editFile.PieceTable{Start: 0, Length: len(content), Source: editFile.ORIGINAL})
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
	_, totalLines := ef.GetContent()
	if totalLines == 0 {
		insertLineNum(0, 0)
		termbox.SetCell(0, 0, ' ', termbox.ColorWhite, termbox.ColorDarkGray)
		return
	} else {
		tui.width, tui.height = termbox.Size()

		for y := 0; y < tui.height && (cur.ScrollY+y) < totalLines; y++ {
			line, _ := ef.GetLine(y)
			insertLineNum(cur.ScrollY+y, y)

			for x := 0; x+cur.ScrollX < len(line) && x+constants.StartX < tui.width; x++ {
				termbox.SetCell(x+constants.StartX, y, rune(line[x+cur.ScrollX]), termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}

	termbox.SetCursor(cur.CursorX-cur.ScrollX, cur.CursorY-cur.ScrollY)
}

func (tui *TUI) displayMessage(message string) {
	ef := tui.ef
	_, height := termbox.Size()
	for _, line := range ef.SplitLines(string(message)) {
		for x, ch := range line {
			termbox.SetCell(x, height-1, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}

func (tui *TUI) HandleSave() (bytes []byte, ok bool) {
	ef := tui.ef
	content, _ := tui.ef.GetContent()
	for _, line := range ef.SplitLines(string(content)) {
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
	}
	tui.ef.IsModified = true
}
