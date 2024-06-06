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
				if lineLength >= cur.CursorX+1 {
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
