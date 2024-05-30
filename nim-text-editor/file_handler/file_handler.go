package filehandler

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
)

const (
	startX = 5
)

var (
	content [][]rune
	cursorY int
	cursorX = startX
)

func InitHandler(filename string) {
	openAndReadFile(filename)
	displayContent()
	handleInput()
	defer termbox.Close()
}

func openAndReadFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	contentByte, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", filename)
		return
	}

	for _, line := range strings.Split(string(contentByte), "\n") {
		content = append(content, []rune(line))
	}
}

func displayContent() {
	if err := termbox.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing TUI: %v\n", err)
		return
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	defer termbox.Flush()

	for i, line := range content {
		lineNum := fmt.Sprintf("%*d", startX-1, i+1)
		for j, n := range lineNum {
			termbox.SetCell(j, i, n, termbox.ColorWhite, termbox.ColorDefault)
		}
		for j, ch := range string(line) {
			termbox.SetCell(j+5, i, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.SetCursor(cursorX, cursorY)
}

func handleInput() {
inputLoop:
	for {
		xi := getXIndex()
		totalLines := len(content)
		lineLength := len(content[cursorY])
		line := content[cursorY]

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlQ:
				break inputLoop
			case termbox.KeyArrowUp:
				if cursorY > 0 {
					cursorY--
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
				}
			case termbox.KeyArrowRight:
				if lineLength > xi+1 {
					changeX(cursorX + 1)
				}
			case termbox.KeyEnter:
				if totalLines < cursorY+1 {
					content = append(content[:cursorY+1], []rune{'\n'})
				} else {
					if lineLength > 1 {
						breakLine := append(content[:cursorY], line[:xi+1])
						if xi >= lineLength {
							content = append(breakLine, append([][]rune{{'\n'}}, content[cursorY+1:]...)...)
						} else {
							content = append(breakLine, append([][]rune{line[xi+1:]}, content[cursorY+1:]...)...)
						}
					} else {
						lines := append(content[:cursorY+1], []rune{'\n'})
						content = append(lines, content[cursorY+1:]...)
					}
				}
				cursorY++
				changeX(startX)
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if cursorY == 0 && xi == 0 {
					return
				}
				if xi == 0 {
					var prevLines [][]rune
					if len(strings.TrimSpace(string(line))) == 0 {
						prevLines = append(content[:cursorY-1], content[cursorY-1])
					} else {
						prevLines = append(content[:cursorY-1], append(content[cursorY-1], line...))
					}
					content = append(prevLines, content[cursorY+1:]...)
					changeX(len(content[cursorY-1]) + startX - 1)
					cursorY--
					break
				} else if lineLength > xi {
					content[cursorY] = append(content[cursorY][:xi-1], content[cursorY][xi:]...)
				} else {
					content[cursorY] = content[cursorY][:xi-1]
				}
				changeX(cursorX - 1)
			default:
				if cursorY >= len(content) {
					content = append(content, []rune{ev.Ch})
				} else if xi >= lineLength {
					content[cursorY] = append(content[cursorY], ev.Ch)
				} else {
					content[cursorY] = append(content[cursorY][:xi], append([]rune{ev.Ch}, content[cursorY][xi:]...)...)
				}
				changeX(cursorX + 1)
			}
		case termbox.EventError:
			fmt.Printf("Termbox error: %v\n", ev.Err)
			break inputLoop
		}
		displayContent()
	}
}

func changeX(val int) {
	if val < startX {
		cursorX = startX
	} else {
		cursorX = val
	}
}

func getXIndex() int {
	return cursorX - startX
}
