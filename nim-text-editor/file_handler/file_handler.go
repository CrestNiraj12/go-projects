package filehandler

import (
	"fmt"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
)

const (
	startX = 5
)

var (
	content  [][]rune
	cursorY  int
	cursorX  = startX
	fileName string
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
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
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

	lineNum := func(i int) {
		lineFormat := fmt.Sprintf("%*d", startX-1, i+1)
		for j, n := range lineFormat {
			termbox.SetCell(j, i, n, termbox.ColorWhite, termbox.ColorDarkGray)
		}
	}

	if len(content) == 0 {
		lineNum(0)
		termbox.SetCell(0, 0, ' ', termbox.ColorWhite, termbox.ColorDarkGray)
	} else {
		for i, line := range content {
			lineNum(i)
			for j, ch := range string(line) {
				termbox.SetCell(j+5, i, ch, termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}

	termbox.SetCursor(cursorX, cursorY)
}

func saveFile() {
	termbox.Close()

	var bytes []byte
	for _, line := range content {
		bytes = append(bytes, []byte(strings.TrimSpace(string(line))+"\n")...)
	}

	err := os.WriteFile(fileName, bytes, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		return
	}
	fmt.Printf("\nWritten %d bytes to file\n", len(bytes))
	os.Exit(0)
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
		if totalLines != 0 {
			lineLength = len(content[cursorY])
			line = content[cursorY]
		}

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
			case termbox.KeyCtrlS:
				saveFile()
			case termbox.KeySpace:
				content[cursorY] = append(line[:xi+1], append([]rune{' '}, line[xi+1:]...)...)
				changeX(cursorX + 1)
			default:
				if ev.Ch == 0 {
					continue inputLoop
				}
				if cursorY >= len(content) {
					content = append(content, []rune{ev.Ch})
				} else if xi >= lineLength {
					content[cursorY] = append(content[cursorY], ev.Ch)
				} else {
					content[cursorY] = append(content[cursorY][:xi+1], append([]rune{ev.Ch}, content[cursorY][xi+1:]...)...)
				}
				changeX(cursorX + 1)
			}
		case termbox.EventError:
			fmt.Printf("Termbox error: %v\n", ev.Err)
			break inputLoop
		}
		displayContent()
	}
	termbox.Close()
}
