package filehandler

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
)

// TODO
// Request save
// Remember cursor position
// Implement better data structure - rope, tabulation, gap buffer, etc
// Undo / Redo
// Word wrap
// Refactor code
const (
	startX = 5
)

var (
	content          [][]rune
	cursorY          int
	cursorX          = startX
	fileName         string
	scrollY, scrollX int
	width, height    int
	isModified       bool
)

func InitHandler(filename string) {
	fileName = filename
	openAndReadFile()
	displayContent()
	handleInput()
}

func openAndReadFile() {
	contentByte, err := os.ReadFile(fileName)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		}
		return
	}

	for _, line := range strings.Split(string(contentByte), "\n") {
		content = append(content, []rune(line))
	}
}

// n is line number, i is y position
// line number (n) increases with line while y position (i) is static on the screen
func insertLineNum(n int, i int) {
	lineFormat := fmt.Sprintf("%*d", startX-1, n+1)
	for j, r := range lineFormat {
		termbox.SetCell(j, i, r, termbox.ColorWhite, termbox.ColorDarkGray)
	}
}

func displayContent() {
	if err := termbox.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing TUI: %v\n", err)
		return
	}
	width, height = termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	defer termbox.Flush()

	if len(content) == 0 {
		insertLineNum(0, 0)
		termbox.SetCell(0, 0, ' ', termbox.ColorWhite, termbox.ColorDarkGray)
	} else {
		width, height := termbox.Size()

		for y := 0; y < height && (scrollY+y) < len(content); y++ {
			line := content[scrollY+y]
			insertLineNum(scrollY+y, y)

			for x := 0; x+scrollX < len(line) && x+startX < width; x++ {
				termbox.SetCell(x+startX, y, rune(line[x+scrollX]), termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}

	termbox.SetCursor(cursorX-scrollX, cursorY-scrollY)
}

func displayMessage(message string) {
	for _, line := range splitLines(message) {
		for x, ch := range line {
			termbox.SetCell(x, height-1, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}

func splitLines(content string) []string {
	return strings.Split(content, "\n")
}

func promptSave() bool {
	displayMessage("You have unsaved changes. Save before exiting? (y/n)")
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Ch == 'y' || ev.Ch == 'Y' {
				return true
			} else if ev.Ch == 'n' || ev.Ch == 'N' {
				return false
			}
		case termbox.EventError:
			return false
		}
	}
}

func saveFile() {
	var bytes []byte
	for _, line := range content {
		bytes = append(bytes, []byte(strings.TrimSpace(string(line))+"\n")...)
	}

	err := os.WriteFile(fileName, bytes, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		return
	}
	isModified = false
	message := fmt.Sprintf("Written %d bytes to file. Press any key to continue", len(bytes))
	displayMessage(message)
	termbox.PollEvent()
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

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlQ:
				if isModified {
					if doSave := promptSave(); doSave {
						saveFile()
					}
				}
				break inputLoop
			case termbox.KeyArrowUp:
				if cursorY > 0 {
					cursorY--
					if cursorY < scrollY {
						scrollY--
					}
				}

				lineLength = len(content[cursorY])
				if totalLines <= cursorY || lineLength <= 0 {
					changeX(startX)
				} else if lineLength < xi {
					changeX(lineLength + startX - 1)
				}

			case termbox.KeyArrowDown:
				if cursorY < totalLines-1 {
					cursorY++
					if cursorY >= scrollY+height {
						scrollY++
					}
				}
				lineLength = len(content[cursorY])
				if totalLines <= cursorY || lineLength <= 0 {
					changeX(startX)
				} else if lineLength < xi {
					changeX(lineLength + startX - 1)
				}
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
				isModified = true
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if cursorY == 0 && xi == 0 {
					continue inputLoop
				}
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
					break
				} else if lineLength > xi {
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
				isModified = true
			case termbox.KeyCtrlS:
				saveFile()
			case termbox.KeySpace:
				if xi == lineLength {
					content[cursorY] = append(line[:xi], rune(' '))
				} else {
					content[cursorY] = append(line[:xi+1], append([]rune{' '}, line[xi+1:]...)...)
				}
				changeX(cursorX + 1)
				isModified = true
			default:
				if ev.Ch == 0 {
					continue inputLoop
				}

				if cursorY >= len(content) {
					content = append(content, []rune{ev.Ch})
				} else if xi >= lineLength {
					content[cursorY] = append(content[cursorY], ev.Ch)
				} else {
					prev := content[cursorY][:xi+1]
					if strings.TrimSpace(string(prev)) == "" {
						content[cursorY] = append([]rune{ev.Ch}, content[cursorY][xi+1:]...)
					} else {
						content[cursorY] = append(prev, append([]rune{ev.Ch}, content[cursorY][xi+1:]...)...)
					}
				}
				changeX(cursorX + 1)
				isModified = true
			}
		case termbox.EventError:
			fmt.Printf("Termbox error: %v\n", ev.Err)
			break inputLoop
		}
		displayContent()
	}
	termbox.Close()
}
