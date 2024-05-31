package file_handler

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
)

func promptSave() bool {
	displayMessage("You have unsaved changes. Save before exiting? (y/n)")
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		if ev.Ch == 'y' || ev.Ch == 'Y' {
			return true
		}
  }
  return false
}

func modifyWrapperFunc(handler interface{}, arg ...interface{}) {
	switch fn := handler.(type) {
	case func():
		fn()
	case func(rune):
		if len(arg) == 1 {
			if argValue, ok := arg[0].(rune); ok {
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
	isModified = true
}

func handleInput() {
	var (
		lineLength int
		line       []rune
	)

	changeX := func(val int) {
		if val < startX {
			cursorX = startX
		} else {
			cursorX = val
		}
	}

inputLoop:
	for {
		xi := cursorX - startX
		totalLines := len(content)
		if totalLines > cursorY {
			lineLength = len(content[cursorY])
			line = content[cursorY]
		}

		onVerticalArrow := func(arrowType termbox.Key) {
			if arrowType == termbox.KeyArrowDown {
				if cursorY < totalLines-1 {
					cursorY++
					if cursorY >= scrollY+height {
						scrollY++
					}
				}
			} else if arrowType == termbox.KeyArrowUp {
				if cursorY > 0 {
					cursorY--
					if cursorY < scrollY {
						scrollY--
					}
				}
			}

			lineLength = len(content[cursorY])
			if totalLines <= cursorY || lineLength <= 0 {
				changeX(startX)
			} else if lineLength < xi {
				changeX(lineLength + startX - 1)
			}
		}

		onEnter := func() {
			if totalLines < cursorY+1 {
				if cursorY+1 > len(content) {
					content = append(content, []rune{'\n'})
				} else {
					content = append(content[:cursorY+1], []rune{'\n'})
				}
			} else {
				if lineLength > 0 {
					if xi > lineLength {
						breakLine := append(content[:cursorY], line[:xi-1])
						content = append(breakLine, append([][]rune{{'\n'}}, content[cursorY+1:]...)...)
					} else if xi == lineLength {
						breakLine := append(content[:cursorY], line[:xi])
						content = append(breakLine, append([][]rune{{'\n'}}, content[cursorY+1:]...)...)
					} else {
						breakLine := append(content[:cursorY], line[:xi+1])
						after := line[xi+1:]
						content = append(breakLine, append([][]rune{after}, content[cursorY+1:]...)...)
					}
				} else {
					lines := append(content[:cursorY], []rune{'\n'})
					content = append(lines, content[cursorY:]...)
				}
			}
			cursorY++
			if cursorY >= scrollY+height {
				scrollY++
			}
			changeX(startX)
		}

		onBackspace := func() {
			if xi == 0 {
				var prevLines [][]rune
				if len(strings.TrimSpace(string(line))) == 0 {
					prevLines = append(content[:cursorY-1], content[cursorY-1])
				} else {
					prevLines = append(content[:cursorY-1], append(content[cursorY-1], line...))
				}
				content = append(prevLines, content[cursorY+1:]...)
				changeX(len(content[cursorY-1]) + startX)
				cursorY--
				if cursorX >= scrollX+width {
					scrollX += len(content[cursorY-1])
				}
			} else {
				if lineLength > xi {
					content[cursorY] = append(content[cursorY][:xi-1], content[cursorY][xi:]...)
				} else {
					content[cursorY] = content[cursorY][:xi-1]
				}
				changeX(cursorX - 1)
				if cursorX < scrollX+startX {
					if width > scrollX {
						scrollX = 0
					} else {
						scrollX -= width
					}
				}
			}
		}

		onCharInput := func(char rune) {
			if cursorY >= len(content) {
				content = append(content, []rune{char})
			} else if xi >= lineLength {
				content[cursorY] = append(content[cursorY], char)
			} else {
				prev := content[cursorY][:xi+1]
				if strings.TrimSpace(string(prev)) == "" {
					content[cursorY] = append([]rune{char}, content[cursorY][xi+1:]...)
				} else {
					content[cursorY] = append(prev, append([]rune{char}, content[cursorY][xi+1:]...)...)
				}
			}
			changeX(cursorX + 1)
		}

		onSpace := func() {
			if xi == lineLength {
				content[cursorY] = append(line[:xi], rune(' '))
			} else {
				content[cursorY] = append(line[:xi+1], append([]rune{' '}, line[xi+1:]...)...)
			}
			changeX(cursorX + 1)
		}

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlQ:
				closeFile()
				break inputLoop
			case termbox.KeyArrowUp:
				onVerticalArrow(termbox.KeyArrowUp)
			case termbox.KeyArrowDown:
				onVerticalArrow(termbox.KeyArrowDown)
			case termbox.KeyArrowLeft:
				if cursorX > startX {
					changeX(cursorX - 1)
					if cursorX < scrollX+startX {
						scrollX--
					}
				}
			case termbox.KeyArrowRight:
				if lineLength >= xi+1 {
					changeX(cursorX + 1)
					if cursorX >= scrollX+width {
						scrollX++
					}
				}
			case termbox.KeyPgup:
				cursorY -= height
				scrollY -= height

				if cursorY < 0 {
					cursorY = 0
				}
				if scrollY < 0 {
					scrollY = 0
				}
			case termbox.KeyPgdn:
				cursorY += height
				scrollY += height

				if cursorY > len(content) {
					cursorY = len(content) - 1
				}
				if scrollY < len(content)-height {
					scrollY = len(content) - height
				}
			case termbox.KeyEnter:
				modifyWrapperFunc(onEnter)
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if cursorY == 0 && xi == 0 {
					continue inputLoop
				}
				modifyWrapperFunc(onBackspace)
			case termbox.KeyCtrlS:
				saveFile()
			case termbox.KeySpace:
				modifyWrapperFunc(onSpace)
			default:
				if ev.Ch == 0 {
					continue inputLoop
				}
				modifyWrapperFunc(onCharInput, ev.Ch)
			}
		case termbox.EventError:
			fmt.Printf("Termbox error: %v\n", ev.Err)
			break inputLoop
		}
		displayContent()
	}
	termbox.Close()
}
