package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Count struct {
	fileName  string
	charCount int
	wordCount int
	lineCount int
}

type FormatArg struct {
	header string
	width  int
}

const (
	file       = "File"
	char       = "Characters"
	word       = "Words"
	line       = "Lines"
	maxHeaders = 4
)

func getFilesAndLabels(args []string) (files []string, labels []string) {
	labels = make([]string, 0, 4)
	labels = append(labels, file)
	if strings.HasPrefix(args[0], "-") {
		files = args[1:]
		var label string
		switch args[0] {
		case "-l":
			label = line
		case "-c":
			label = char
		case "-w":
			label = word
		}
		labels = append(labels, label)
	} else {
		files = args
		labels = append(labels, char, word, line)
	}
	return files, labels
}

func isValidHeader(labels []string, label string, currHeader string) bool {
	return slices.Contains(labels, label) && currHeader == label
}

func printDivider(width []int) {
	fmt.Print("|")
	for _, w := range width {
		fmt.Print(strings.Repeat("-", w))
		fmt.Print("|")
	}
	fmt.Println()
}

func printBorder(width []int) {
	fmt.Print(">")
	for i, w := range width {
		fmt.Print(strings.Repeat("-", w))
		fmt.Print("-")
		if i == len(width)-1 {
			fmt.Print("<")
		}
	}
	fmt.Println()
}

func throwErrorAndExit() {
	fmt.Fprintln(os.Stderr, "Error: No arguments found\nUsage: go run main.go [-l] file1 file2 ...")
	os.Exit(-1)
}

// TODO
// Make Arguments Dynamic -l, -c, -w
// Make Print dynamic according to args [Store the counts and calculate and print later]
// Adjust file width according to largest text file name
func main() {
	if len(os.Args) < 2 {
		throwErrorAndExit()
	}
	args := os.Args[1:]
	files, labels := getFilesAndLabels(args)
	if len(files) == 0 {
		throwErrorAndExit()
	}

	var counts []Count
	for _, fName := range files {
		var wc, cc, lc int
		file, err := os.Open(fName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error while reading file")
			os.Exit(-1)
		}

		scan := bufio.NewScanner(file)
		for scan.Scan() {
			line := scan.Text()
			cc += len(line)
			wc += len(strings.Fields(line))
			lc++
		}

		counts = append(counts, Count{fName, cc, wc, lc})
		file.Close()
	}

	width := make([]int, 0, len(labels))
	for _, header := range labels {
		width = append(width, len(header))
	}

	var tc, tw, tl int
	for _, count := range counts {
		for i, header := range labels {
			var currLen int
			if header == file {
				currLen = len(count.fileName)
			} else if isValidHeader(labels, char, header) {
				currLen = len(fmt.Sprintf("%d", count.charCount))
			} else if isValidHeader(labels, word, header) {
				currLen = len(fmt.Sprintf("%d", count.wordCount))
			} else if isValidHeader(labels, line, header) {
				currLen = len(fmt.Sprintf("%d", count.lineCount))
			}

			if currLen > width[i] {
				width[i] = currLen
			}
		}
		tc += count.charCount
		tw += count.wordCount
		tl += count.lineCount
	}

	stringFormat := "%-*s"
	digitFormat := "%*d"
	printBorder(width)

	fmt.Print("|")
	for i := range labels {
		fmt.Printf(stringFormat+"|", width[i], labels[i])
	}
	fmt.Println()
	printDivider(width)

	for _, count := range counts {
		fmt.Printf("|"+stringFormat+"|", width[0], count.fileName)
		for i, header := range labels[1:] {
			var value int
			if isValidHeader(labels, char, header) {
				value = count.charCount
			} else if isValidHeader(labels, word, header) {
				value = count.wordCount
			} else if isValidHeader(labels, line, header) {
				value = count.lineCount
			}
			fmt.Printf(digitFormat+"|", width[i+1], value)
		}
		fmt.Println()
	}

	if len(files) > 1 {
		printDivider(width)
		fmt.Printf("|"+stringFormat+"|", width[0], "Total")
		for i, header := range labels[1:] {
			var value int
			if isValidHeader(labels, char, header) {
				value = tc
			} else if isValidHeader(labels, word, header) {
				value = tw
			} else if isValidHeader(labels, line, header) {
				value = tl
			}
			fmt.Printf(digitFormat+"|", width[i+1], value)
		}
		fmt.Println()
	}
	printBorder(width)
}
