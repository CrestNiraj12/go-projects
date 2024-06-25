package tui

import (
	constants "nim-text-editor/constants"
	editFile "nim-text-editor/edit_file"

	"github.com/nsf/termbox-go"
)

func (tui *TUI) onVerticalArrow(arrowType termbox.Key) {
	startX := tui.startX
	cur := tui.ef.Cursor
	switch arrowType {
	case termbox.KeyArrowDown:
		tui.moveDown()
		tui.moveVertically()
	case termbox.KeyArrowUp:
		tui.moveUp()
		tui.moveVertically()
	case termbox.KeyArrowLeft:
		tui.moveUp()
		_, lineLength := tui.ef.GetLine(0)
		cur.ChangeX(lineLength + startX - 1)
	case termbox.KeyArrowRight:
		tui.moveDown()
		cur.ChangeX(startX)
	}
	_, lineLength := tui.ef.GetLine(0)
	tui.scrollX(startX + lineLength)
}

func (tui *TUI) moveVertically() {
	ef := tui.ef
	cur := ef.Cursor
	startX := tui.startX
	_, totalLines := ef.GetContent()
	_, lineLength := ef.GetLine(0)
	if totalLines <= cur.CursorY || lineLength <= 0 {
		ef.SetXMemo()
		cur.ChangeX(startX)
	} else if lineLength < cur.GetCurXIndex() {
		ef.SetXMemo()
		cur.ChangeX(lineLength + startX - 1)
	} else {
		if ef.XMemoCur == 0 {
			return
		}
		cur.ChangeX(ef.XMemoCur)
	}
}

func (tui *TUI) moveDown() {
	cur := tui.ef.Cursor
	_, totalLines := tui.ef.GetContent()
	if cur.CursorY < totalLines-1 {
		cur.CursorY++
		if cur.CursorY >= cur.ScrollY+tui.height {
			cur.ScrollY++
		}
	}
}

func (tui *TUI) moveUp() {
	cur := tui.ef.Cursor
	if cur.CursorY > 0 {
		cur.CursorY--
		if cur.CursorY < cur.ScrollY {
			cur.ScrollY--
		}
	}
}

func (tui *TUI) onSpace() {
	ef := tui.ef
	cur := ef.Cursor
	tui.AddPieces([]rune{' '})
	cur.ChangeX(cur.CursorX + 1)
}

func (tui *TUI) onEnter() {
	ef := tui.ef
	cur := ef.Cursor
	_, lineLength := ef.GetLine(0)
	tui.AddPieces([]rune{'\n'})
	cur.CursorY++
	if cur.CursorY >= cur.ScrollY+tui.height {
		cur.ScrollY++
	}
	cur.ChangeX(constants.StartX)
	tui.scrollX(tui.startX + lineLength)
}

func (tui *TUI) onBackspace() {
	tui.RemovePiece()
	ef := tui.ef
	cur := ef.Cursor
	xi := cur.GetCurXIndex()
	if xi == 0 {
		_, prevLineLen := ef.GetLine(-1)
		cur.ChangeX(prevLineLen + tui.startX)
		cur.CursorY--
		if cur.CursorX >= cur.ScrollX+tui.width {
			cur.ScrollX += prevLineLen
		}
	} else {
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
	tui.AddPieces([]rune{char})
	cur.ChangeX(cur.CursorX + 1)
	tui.scrollX(1)
}

func (tui *TUI) AddPieces(input []rune) {
	ef := tui.ef
	*ef.Content.Add = append(*ef.Content.Add, input...)

	pieces := ef.Content.Pieces
	curX := ef.Cursor.GetCurXIndex()
	inputLen := len(string(input))
	if curX >= ef.GetFileLength() {
		if len(pieces) <= 1 {
			pieces = append(pieces, &editFile.PieceTable{
				Start:  curX,
				Length: inputLen,
				Source: editFile.ADD,
			})
			return
		} else {

			lastPiece := pieces[len(pieces)-1]
			if lastPiece.Source == editFile.ADD {
				lastPiece.Length += inputLen
			}
		}
		return
	}

	addPiecesArr := make([]*editFile.PieceTable, 0, 3+len(pieces)-1)

	for i, _ := range pieces {
		piece := pieces[i]
		if curX >= piece.Start && curX <= piece.Start+piece.Length {
			var source *[]rune
			if piece.Source == editFile.ADD {
				source = ef.Content.Add
			} else {
				source = ef.Content.Original
			}
			splitSizeLeading := len((*source)[:curX])
			if curX != piece.Start {
				addPiecesArr = append(addPiecesArr, &editFile.PieceTable{
					Start:  piece.Start,
					Length: splitSizeLeading,
					Source: piece.Source,
				})
			}

			addPiecesArr = append(addPiecesArr, &editFile.PieceTable{
				Start:  tui.ef.GetAddLength() - inputLen,
				Length: inputLen,
				Source: editFile.ADD,
			})

			if curX != piece.Start+piece.Length {
				addPiecesArr = append(addPiecesArr, &editFile.PieceTable{
					Start:  splitSizeLeading,
					Length: piece.Length - splitSizeLeading,
					Source: piece.Source,
				})
				pieces = append(pieces[:i], append(addPiecesArr, pieces[i+1:]...)...)
			}
			break
		}
	}
}

func (tui *TUI) RemovePiece() {
	ef := tui.ef
	pieces := ef.Content.Pieces
	curX := ef.Cursor.GetCurXIndex()
	if curX >= ef.GetFileLength() {
		lastPiece := len(pieces) - 1
		pieces[lastPiece].Length--
		if pieces[lastPiece].Length < 1 {
			pieces = pieces[:lastPiece]
		}
		return
	}

	removePiecesArr := make([]*editFile.PieceTable, 0, 2+len(pieces)-1)

	for i := range pieces {
		piece := pieces[i]
		if curX > piece.Start && curX <= piece.Start+piece.Length {
			var source *[]rune
			if piece.Source == editFile.ADD {
				source = ef.Content.Add
			} else {
				source = ef.Content.Original
			}
			if curX == piece.Start+1 {
				piece.Start++
				piece.Length--
			} else if curX == piece.Start+piece.Length {
				piece.Length--
			} else {
				splitSizeLeading := len((*source)[:curX])
				removePiecesArr = append(removePiecesArr, &editFile.PieceTable{
					Start:  piece.Start,
					Length: splitSizeLeading,
					Source: piece.Source,
				})
				removePiecesArr = append(removePiecesArr, &editFile.PieceTable{
					Start:  piece.Start + splitSizeLeading,
					Length: piece.Length - splitSizeLeading - 1,
					Source: piece.Source,
				})
				pieces = append(pieces[:i], append(removePiecesArr, pieces[i+1:]...)...)
				break
			}
			if piece.Length < 1 {
				pieces = append(pieces[:i], pieces[i+1:]...)
			}
			break
		}
	}
}

func (tui *TUI) scrollX(val int) {
	cur := tui.ef.Cursor
	if cur.CursorX >= tui.width+cur.ScrollX {
		if val == 1 {
			cur.ScrollX++
		} else {
			cur.ScrollX = val - tui.width
		}
	} else if cur.GetCurXIndex() < cur.ScrollX {
		if val == 1 {
			cur.ScrollX--
		} else {
			cur.ScrollX = 0
		}
	}
}
