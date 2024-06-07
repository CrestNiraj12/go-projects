package tui

import (
	"fmt"
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
	if bytes, ok := tui.HandleSave(); ok {
		ef.IsModified = false
		message := fmt.Sprintf("Written %d bytes to file. Press any key to continue", len(bytes))
		displayMessage(message)
		termbox.PollEvent()
	}
}

func (tui *TUI) handleInput() {
	ef := tui.ef
	startX := tui.startX

inputLoop:
	for {
		totalLines := ef.GetTotalLines()
		width, height := termbox.Size()
		if width != tui.width || height != tui.height {
			tui.width, tui.height = width, height
		}
		cur := ef.Cursor
		_, lineLength := ef.GetLine()

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlQ:
				tui.handleSave(true)
				break inputLoop
			case termbox.KeyArrowUp:
				tui.onVerticalArrow(termbox.KeyArrowUp)
			case termbox.KeyArrowDown:
				tui.onVerticalArrow(termbox.KeyArrowDown)
			case termbox.KeyArrowLeft:
				if cur.CursorX > startX {
					cur.ChangeX(cur.CursorX - 1)
					tui.scrollX(1)
				} else {
          tui.onVerticalArrow(termbox.KeyArrowLeft)
        }
			case termbox.KeyArrowRight:
				if lineLength >= cur.GetCurXIndex()+1 {
					cur.ChangeX(cur.CursorX + 1)
					tui.scrollX(1)
				} else {
          tui.onVerticalArrow(termbox.KeyArrowRight)
				}
			case termbox.KeyPgup:
				cur.CursorY -= tui.height
				cur.ScrollY -= tui.height

				if cur.CursorY < 0 {
					cur.CursorY = 0
				}
				if cur.ScrollY < 0 {
					cur.ScrollY = 0
				}
			case termbox.KeyPgdn:
				cur.CursorY += tui.height
				cur.ScrollY += tui.height

				if cur.CursorY > totalLines {
					cur.CursorY = totalLines - 1
				}
				if cur.ScrollY < totalLines-tui.height {
					cur.ScrollY = totalLines - tui.height
				}
			case termbox.KeyEnter:
				tui.modifyWrapperFunc(tui.onEnter)
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if cur.CursorY == 0 && cur.GetCurXIndex() == 0 {
					continue inputLoop
				}
				tui.modifyWrapperFunc(tui.onBackspace)
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
