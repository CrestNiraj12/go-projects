package tui

import (
	"fmt"
	constants "nim-text-editor/constants"
	"strings"
	"unicode"

	"github.com/nsf/termbox-go"
)

func promptSave() (prompt bool) {
	displayMessage("You have unsaved changes. Save before exiting? (y/n)")
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		prompt = unicode.ToLower(ev.Ch) == 'y'
	}
	return
}

func (tui *TUI) handleSave(prompt bool) {
	ef := tui.ef
	if !ef.IsModified {
		return
	}
	var doSave = true
	if prompt {
		doSave = promptSave()
	}
	if !doSave {
		return
	}
	if bytes, ok := tui.Save(); ok {
		ef.IsModified = false
		message := fmt.Sprintf("Written %d bytes to file. Press any key to continue", len(bytes))
		displayMessage(message)
		termbox.PollEvent()
	}
}

func (tui *TUI) handleInput() {
	ef := tui.ef
	cur := ef.Cursor
	startX := tui.startX

	width, height := termbox.Size()
	_, lineLength := ef.GetLine()

inputLoop:
	for {

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlQ:
				tui.handleSave(true)
				break inputLoop
			case termbox.KeyArrowUp:
				tui.onVerticalArrow(termbox.KeyArrowUp, height)
			case termbox.KeyArrowDown:
				tui.onVerticalArrow(termbox.KeyArrowDown, height)
			case termbox.KeyArrowLeft:
				if cur.CursorX > startX {
					cur.ChangeX(cur.CursorX - 1)
					if cur.CursorX < cur.ScrollX+startX {
						cur.ScrollX--
					}
				}
			case termbox.KeyArrowRight:
				if lineLength >= ef.Cursor.CursorX+1 {
					cur.ChangeX(cur.CursorX + 1)
					if cur.CursorX >= cur.ScrollX+width {
						cur.ScrollX++
					}
				}
			case termbox.KeyPgup:
				cur.CursorY -= height
				cur.ScrollY -= height

				if cur.CursorY < 0 {
					cur.CursorY = 0
				}
				if cur.ScrollY < 0 {
					cur.ScrollY = 0
				}
			case termbox.KeyPgdn:
				cur.CursorY += height
				cur.ScrollY += height

				if cur.CursorY > len(ef.Content) {
					cur.CursorY = len(ef.Content) - 1
				}
				if cur.ScrollY < len(ef.Content)-height {
					cur.ScrollY = len(ef.Content) - height
				}
			case termbox.KeyEnter:
				tui.modifyWrapperFunc(tui.onEnter, height)
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if cur.CursorY == 0 && cur.GetCurXIndex() == 0 {
					continue inputLoop
				}
				tui.modifyWrapperFunc(tui.onBackspace, width)
			case termbox.KeyCtrlS:
				tui.handleSave(false)
			case termbox.KeySpace:
				tui.modifyWrapperFunc(tui.onSpace)
			default:
				if ev.Ch == 0 {
					continue inputLoop
				}
				tui.modifyWrapperFunc(tui.onCharInput, ev.Ch)
			}
		case termbox.EventError:
			fmt.Printf("Termbox error: %v\n", ev.Err)
			break inputLoop
		}
		tui.displayContent()
	}
	termbox.Close()
}

func (tui *TUI) onVerticalArrow(arrowType termbox.Key, height int) {
	ef := tui.ef
	cur := ef.Cursor
	startX := tui.startX
	totalLines := ef.GetTotalLines()
	if arrowType == termbox.KeyArrowDown {
		if cur.CursorY < totalLines-1 {
			cur.CursorY++
			if cur.CursorY >= cur.ScrollY+height {
				cur.ScrollY++
			}
		}
	} else if arrowType == termbox.KeyArrowUp {
		if cur.CursorY > 0 {
			cur.CursorY--
			if cur.CursorY < cur.ScrollY {
				cur.ScrollY--
			}
		}
	}

	lineLength := len(ef.Content[cur.CursorY])
	if totalLines <= cur.CursorY || lineLength <= 0 {
		cur.ChangeX(startX)
	} else if lineLength < cur.CursorX {
		cur.ChangeX(lineLength + startX - 1)
	}

}

func (tui *TUI) onSpace() {
	ef := tui.ef
	cur := ef.Cursor
	line, lineLength := ef.GetLine()
	xi := cur.GetCurXIndex()
	if xi == lineLength {
		ef.Content[cur.CursorY] = append(line[:xi], rune(' '))
	} else {
		ef.Content[cur.CursorY] = append(line[:xi+1], append([]rune{' '}, line[xi+1:]...)...)
	}
	cur.ChangeX(cur.CursorX + 1)
}

func (tui *TUI) onEnter(height int) {
	ef := tui.ef
	cur := ef.Cursor
	totalLines := ef.GetTotalLines()
	line, lineLength := ef.GetLine()

	if totalLines < cur.CursorY+1 {
		if cur.CursorY+1 > len(ef.Content) {
			ef.Content = append(ef.Content, []rune{'\n'})
		} else {
			ef.Content = append(ef.Content[:cur.CursorY+1], []rune{'\n'})
		}
	} else {
		xi := cur.GetCurXIndex()
		if lineLength > 0 {
			if xi > lineLength {
				breakLine := append(ef.Content[:cur.CursorY], line[:xi-1])
				ef.Content = append(breakLine, append([][]rune{{'\n'}}, ef.Content[cur.CursorY+1:]...)...)
			} else if xi == lineLength {
				breakLine := append(ef.Content[:cur.CursorY], line[:xi])
				ef.Content = append(breakLine, append([][]rune{{'\n'}}, ef.Content[cur.CursorY+1:]...)...)
			} else {
				breakLine := append(ef.Content[:cur.CursorY], line[:xi+1])
				after := line[xi+1:]
				ef.Content = append(breakLine, append([][]rune{after}, ef.Content[cur.CursorY+1:]...)...)
			}
		} else {
			lines := append(ef.Content[:cur.CursorY], []rune{'\n'})
			ef.Content = append(lines, ef.Content[cur.CursorY:]...)
		}
	}
	cur.CursorY++
	if cur.CursorY >= cur.ScrollY+height {
		cur.ScrollY++
	}
	cur.ChangeX(constants.StartX)
}

func (tui *TUI) onBackspace(width int) {
	ef := tui.ef
	cur := ef.Cursor
	xi := cur.GetCurXIndex()
	line, lineLength := ef.GetLine()
	if xi == 0 {
		var prevLines [][]rune
		if len(strings.TrimSpace(string(line))) == 0 {
			prevLines = append(ef.Content[:cur.CursorY-1], ef.Content[cur.CursorY-1])
		} else {
			prevLines = append(ef.Content[:cur.CursorY-1], append(ef.Content[cur.CursorY-1], line...))
		}
		ef.Content = append(prevLines, ef.Content[cur.CursorY+1:]...)
		cur.ChangeX(len(ef.Content[cur.CursorY-1]) + tui.startX)
		cur.CursorY--
		if cur.CursorX >= cur.ScrollX+width {
			cur.ScrollX += len(ef.Content[cur.CursorY-1])
		}
	} else {
		if lineLength > xi {
			ef.Content[cur.CursorY] = append(ef.Content[cur.CursorY][:xi-1], ef.Content[cur.CursorY][xi:]...)
		} else {
			ef.Content[cur.CursorY] = ef.Content[cur.CursorY][:xi-1]
		}
		cur.ChangeX(cur.CursorX - 1)
		if cur.CursorX >= cur.ScrollX+tui.startX {
			return
		}
		if width > cur.ScrollX {
			cur.ScrollX = 0
			return
		}
		cur.ScrollX -= width
	}
}

func (tui *TUI) onCharInput(char rune) {
	ef := tui.ef
	cur := ef.Cursor
	_, lineLength := ef.GetLine()
	xi := cur.GetCurXIndex()
	if cur.CursorY >= len(ef.Content) {
		ef.Content = append(ef.Content, []rune{char})
	} else if xi >= lineLength {
		ef.Content[cur.CursorY] = append(ef.Content[cur.CursorY], char)
	} else {
		prev := ef.Content[cur.CursorY][:xi+1]
		if strings.TrimSpace(string(prev)) == "" {
			ef.Content[cur.CursorY] = append([]rune{char}, ef.Content[cur.CursorY][xi+1:]...)
		} else {
			ef.Content[cur.CursorY] = append(prev, append([]rune{char}, ef.Content[cur.CursorY][xi+1:]...)...)
		}
	}
	cur.ChangeX(cur.CursorX + 1)
}
