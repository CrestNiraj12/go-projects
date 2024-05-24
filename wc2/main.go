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

type TotalCount struct {
	totalChar  int
	totalWords int
	totalLines int
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

func throwErrorAndExit() {
	fmt.Fprintln(os.Stderr, "\nError: No arguments found\nUsage: go run main.go [-l] file1 file2 ...")
	os.Exit(-1)
}

func openFileAndStoreCount(files []string, counts *[]Count) {
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

		*counts = append(*counts, Count{fName, cc, wc, lc})
		file.Close()
	}
}

func main() {
	if len(os.Args) < 2 {
		throwErrorAndExit()
	}
	args := os.Args[1:]
	files, labels := getFilesAndLabels(args)
	if len(files) == 0 {
		throwErrorAndExit()
	}

	counts := make([]Count, 0, len(files))
	openFileAndStoreCount(files, &counts)
   
	width := make([]int, 0, len(labels))
	for _, header := range labels {
		width = append(width, len(header))
	}

	isValidLabel := func(label string) bool {
		return slices.Contains(labels, label)
	}

	getCountByLabel := func(count interface{}, label string) (val int) {
		if !isValidLabel(label) {
			return val
		}

		switch v := count.(type) {
		case *Count:
			switch label {
			case char:
				val = (*v).charCount
			case word:
				val = (*v).wordCount
			case line:
				val = (*v).lineCount
			}
		case *TotalCount:
			switch label {
			case char:
				val = (*v).totalChar
			case word:
				val = (*v).totalWords
			case line:
				val = (*v).totalLines
			}
		}

		return val
	}

	var totalCount TotalCount
	for _, count := range counts {
		for i, header := range labels {
			var currLen int
			if header == file {
				currLen = len(count.fileName)
			} else {
				currLen = len(fmt.Sprintf("%d", getCountByLabel(&count, header)))
			}

			if currLen > width[i] {
				width[i] = currLen
			}
		}
		totalCount.totalChar += count.charCount
		totalCount.totalWords += count.wordCount
		totalCount.totalLines += count.lineCount
	}

	stringFormat := "%-*s"
	digitFormat := "%*d"
	printDivider := func() {
		fmt.Print("|")
		for _, w := range width {
			fmt.Print(strings.Repeat("-", w))
			fmt.Print("|")
		}
		fmt.Println()
	}

	printBorder := func() {
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
	printBorder()

	fmt.Print("|")
	for i := range labels {
		fmt.Printf(stringFormat+"|", width[i], labels[i])
	}
	fmt.Println()
	printDivider()

	for _, count := range counts {
		fmt.Printf("|"+stringFormat+"|", width[0], count.fileName)
		for i, header := range labels[1:] {
			fmt.Printf(digitFormat+"|", width[i+1], getCountByLabel(&count, header))
		}
		fmt.Println()
	}

	if len(files) > 1 {
		printDivider()
		fmt.Printf("|"+stringFormat+"|", width[0], "Total")
		for i, header := range labels[1:] {
			fmt.Printf(digitFormat+"|", width[i+1], getCountByLabel(&totalCount, header))
		}
		fmt.Println()
	}
	printBorder()
}
