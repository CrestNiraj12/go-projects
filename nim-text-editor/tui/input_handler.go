package tui

import (
	"fmt"
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

func (ef *EditFile) handleSave(prompt bool) {
	if !ef.isModified {
		return
	}
	var doSave = true
	if prompt {
		doSave = promptSave()
	}
	if !doSave {
		return
	}
	if bytes, ok := ef.Save(); ok {
		ef.isModified = false
		message := fmt.Sprintf("Written %d bytes to file. Press any key to continue", len(bytes))
		displayMessage(message)
		termbox.PollEvent()
	}
}

func (ef *EditFile) handleInput() {
	var (
		lineLength int
		line       []rune
	)

	width, height := termbox.Size()

	changeX := func(val int) {
		if val < startX {
			ef.cursorX = startX
			return
		}
		ef.cursorX = val
	}

inputLoop:
	for {
		xi := ef.cursorX - startX
		totalLines := len(ef.content)
		if totalLines > ef.cursorY {
			lineLength = len(ef.content[ef.cursorY])
			line = ef.content[ef.cursorY]
		}

		onVerticalArrow := func(arrowType termbox.Key) {
			if arrowType == termbox.KeyArrowDown {
				if ef.cursorY < totalLines-1 {
					ef.cursorY++
					if ef.cursorY >= ef.scrollY+height {
						ef.scrollY++
					}
				}
			} else if arrowType == termbox.KeyArrowUp {
				if ef.cursorY > 0 {
					ef.cursorY--
					if ef.cursorY < ef.scrollY {
						ef.scrollY--
					}
				}
			}

			lineLength = len(ef.content[ef.cursorY])
			if totalLines <= ef.cursorY || lineLength <= 0 {
				changeX(startX)
			} else if lineLength < xi {
				changeX(lineLength + startX - 1)
			}
		}

		onEnter := func() {
			if totalLines < ef.cursorY+1 {
				if ef.cursorY+1 > len(ef.content) {
					ef.content = append(ef.content, []rune{'\n'})
				} else {
					ef.content = append(ef.content[:ef.cursorY+1], []rune{'\n'})
				}
			} else {
				if lineLength > 0 {
					if xi > lineLength {
						breakLine := append(ef.content[:ef.cursorY], line[:xi-1])
						ef.content = append(breakLine, append([][]rune{{'\n'}}, ef.content[ef.cursorY+1:]...)...)
					} else if xi == lineLength {
						breakLine := append(ef.content[:ef.cursorY], line[:xi])
						ef.content = append(breakLine, append([][]rune{{'\n'}}, ef.content[ef.cursorY+1:]...)...)
					} else {
						breakLine := append(ef.content[:ef.cursorY], line[:xi+1])
						after := line[xi+1:]
						ef.content = append(breakLine, append([][]rune{after}, ef.content[ef.cursorY+1:]...)...)
					}
				} else {
					lines := append(ef.content[:ef.cursorY], []rune{'\n'})
					ef.content = append(lines, ef.content[ef.cursorY:]...)
				}
			}
			ef.cursorY++
			if ef.cursorY >= ef.scrollY+height {
				ef.scrollY++
			}
			changeX(startX)
		}

		onBackspace := func() {
			if xi == 0 {
				var prevLines [][]rune
				if len(strings.TrimSpace(string(line))) == 0 {
					prevLines = append(ef.content[:ef.cursorY-1], ef.content[ef.cursorY-1])
				} else {
					prevLines = append(ef.content[:ef.cursorY-1], append(ef.content[ef.cursorY-1], line...))
				}
				ef.content = append(prevLines, ef.content[ef.cursorY+1:]...)
				changeX(len(ef.content[ef.cursorY-1]) + startX)
				ef.cursorY--
				if ef.cursorX >= ef.scrollX+width {
					ef.scrollX += len(ef.content[ef.cursorY-1])
				}
			} else {
				if lineLength > xi {
					ef.content[ef.cursorY] = append(ef.content[ef.cursorY][:xi-1], ef.content[ef.cursorY][xi:]...)
				} else {
					ef.content[ef.cursorY] = ef.content[ef.cursorY][:xi-1]
				}
				changeX(ef.cursorX - 1)
				if ef.cursorX < ef.scrollX+startX {
					if width > ef.scrollX {
						ef.scrollX = 0
					} else {
						ef.scrollX -= width
					}
				}
			}
		}

		onCharInput := func(char rune) {
			if ef.cursorY >= len(ef.content) {
				ef.content = append(ef.content, []rune{char})
			} else if xi >= lineLength {
				ef.content[ef.cursorY] = append(ef.content[ef.cursorY], char)
			} else {
				prev := ef.content[ef.cursorY][:xi+1]
				if strings.TrimSpace(string(prev)) == "" {
					ef.content[ef.cursorY] = append([]rune{char}, ef.content[ef.cursorY][xi+1:]...)
				} else {
					ef.content[ef.cursorY] = append(prev, append([]rune{char}, ef.content[ef.cursorY][xi+1:]...)...)
				}
			}
			changeX(ef.cursorX + 1)
		}

		onSpace := func() {
			if xi == lineLength {
				ef.content[ef.cursorY] = append(line[:xi], rune(' '))
			} else {
				ef.content[ef.cursorY] = append(line[:xi+1], append([]rune{' '}, line[xi+1:]...)...)
			}
			changeX(ef.cursorX + 1)
		}

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlQ:
				ef.handleSave(true)
				break inputLoop
			case termbox.KeyArrowUp:
				onVerticalArrow(termbox.KeyArrowUp)
			case termbox.KeyArrowDown:
				onVerticalArrow(termbox.KeyArrowDown)
			case termbox.KeyArrowLeft:
				if ef.cursorX > startX {
					changeX(ef.cursorX - 1)
					if ef.cursorX < ef.scrollX+startX {
						ef.scrollX--
					}
				}
			case termbox.KeyArrowRight:
				if lineLength >= xi+1 {
					changeX(ef.cursorX + 1)
					if ef.cursorX >= ef.scrollX+width {
						ef.scrollX++
					}
				}
			case termbox.KeyPgup:
				ef.cursorY -= height
				ef.scrollY -= height

				if ef.cursorY < 0 {
					ef.cursorY = 0
				}
				if ef.scrollY < 0 {
					ef.scrollY = 0
				}
			case termbox.KeyPgdn:
				ef.cursorY += height
				ef.scrollY += height

				if ef.cursorY > len(ef.content) {
					ef.cursorY = len(ef.content) - 1
				}
				if ef.scrollY < len(ef.content)-height {
					ef.scrollY = len(ef.content) - height
				}
			case termbox.KeyEnter:
				ef.modifyWrapperFunc(onEnter)
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if ef.cursorY == 0 && xi == 0 {
					continue inputLoop
				}
				ef.modifyWrapperFunc(onBackspace)
			case termbox.KeyCtrlS:
				ef.handleSave(false)
			case termbox.KeySpace:
				ef.modifyWrapperFunc(onSpace)
			default:
				if ev.Ch == 0 {
					continue inputLoop
				}
				ef.modifyWrapperFunc(onCharInput, ev.Ch)
			}
		case termbox.EventError:
			fmt.Printf("Termbox error: %v\n", ev.Err)
			break inputLoop
		}
		ef.displayContent()
	}
	termbox.Close()
}
