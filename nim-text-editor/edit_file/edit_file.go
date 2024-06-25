package editFile

import (
	constants "nim-text-editor/constants"
	"strings"
)

const ORIGINAL = "original"
const ADD = "add"

type FileCursor struct {
	CursorX, CursorY, ScrollY, ScrollX int
}

type EditFile struct {
	Cursor     *FileCursor
	Content    *ContentTable
	IsModified bool
	FileName   string
	XMemoCur   int
}

type ContentTable struct {
	Original *[]rune
	Add      *[]rune
	Pieces   []*PieceTable
}

type PieceTable struct {
	Start, Length int
	Source        string
}

func (ef *EditFile) GetFileLength() (length int) {
	for _, piece := range ef.Content.Pieces {
		length += piece.Length
	}
	return
}

func (ef *EditFile) GetAddLength() int {
	return len(string(*ef.Content.Add))
}

func (ef *EditFile) GetOriginalLength() int {
	return len(string(*ef.Content.Original))
}

func (ef *EditFile) SplitLines(content string) []string {
	return strings.Split(content, "\n")
}

func (ef *EditFile) GetLine(offset int) (line []rune, lineLength int) {
	content, _ := ef.GetContent()
	lineString := ef.SplitLines(string(content))[ef.Cursor.CursorY+offset]
	line = []rune(lineString)
	lineLength = len(lineString)
	return
}

func (ef *EditFile) GetContent() (content []rune, totalLines int) {
	var buffer []rune

	countLines := func(text []rune) int {
		return len(ef.SplitLines(string(text)))
	}

	for _, piece := range ef.Content.Pieces {
		if piece.Source == ORIGINAL {
			buffer = *ef.Content.Original
		} else {
			buffer = *ef.Content.Add
		}
		textSegment := buffer[piece.Start : piece.Start+piece.Length]
		totalLines += countLines(textSegment) - 1
		content = append(content, textSegment...)
	}
	return
}

func (cur *FileCursor) GetCurXIndex() int {
	return cur.CursorX - constants.StartX
}

func (ef *EditFile) SetXMemo() {
	curX := ef.Cursor.CursorX
	if curX <= ef.XMemoCur {
		return
	}
	ef.XMemoCur = curX
}

func (cur *FileCursor) ChangeX(val int) {
	startX := constants.StartX
	if val < startX {
		cur.CursorX = startX
		return
	}
	cur.CursorX = val
}
