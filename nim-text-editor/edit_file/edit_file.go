package editFile

import (
	constants "nim-text-editor/constants"
)

type FileCursor struct {
	CursorX, CursorY, ScrollY, ScrollX int
}

type EditFile struct {
	Cursor     *FileCursor
	Content    [][]rune
	IsModified bool
	FileName   string
}

func (ef *EditFile) GetTotalLines() int {
	return len(ef.Content)
}

func (ef *EditFile) GetLine() (line []rune, lineLength int) {
	if ef.GetTotalLines() > ef.Cursor.CursorY {
		line = ef.Content[ef.Cursor.CursorY]
		lineLength = len(line)
	}
	return
}

func (cur *FileCursor) GetCurXIndex() int {
	return cur.CursorX - constants.StartX
}

func (cur *FileCursor) ChangeX(val int) {
	startX := constants.StartX
	if val < startX {
		cur.CursorX = startX
		return
	}
	cur.CursorX = val
}
