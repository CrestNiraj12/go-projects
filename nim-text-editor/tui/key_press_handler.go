package tui

import (
	constants "nim-text-editor/constants"
	"strings"

	"github.com/nsf/termbox-go"
)

func (tui *TUI) onVerticalArrow(arrowType termbox.Key) {
	ef := tui.ef
	cur := ef.Cursor
	startX := tui.startX
	totalLines := ef.GetTotalLines()
	if arrowType == termbox.KeyArrowDown {
		if cur.CursorY < totalLines-1 {
			cur.CursorY++
			if cur.CursorY >= cur.ScrollY+tui.height {
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

	_, lineLength := ef.GetLine()
	if totalLines <= cur.CursorY || lineLength <= 0 {
		cur.ChangeX(startX)
	} else if lineLength < cur.GetCurXIndex() {
		cur.ChangeX(lineLength + startX - 1)
	}

	tui.scrollX()
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

func (tui *TUI) onEnter() {
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
	if cur.CursorY >= cur.ScrollY+tui.height {
		cur.ScrollY++
	}
	cur.ChangeX(constants.StartX)
}

func (tui *TUI) onBackspace() {
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
		if cur.CursorX >= cur.ScrollX+tui.width {
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
		if tui.width > cur.ScrollX {
			cur.ScrollX = 0
			return
		}
		cur.ScrollX -= tui.width
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

func (tui *TUI) scrollX() {
	cur := tui.ef.Cursor
	startX := tui.startX
	if cur.CursorX < cur.ScrollX+startX {
		cur.ScrollX--
	} else if cur.CursorX >= cur.ScrollX+tui.width {
		cur.ScrollX++
	}
}
