package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func SearchComics() {
	if len(os.Args) <= 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run . search_comics <keyword1> <keyword2>...")
		os.Exit(-1)
	}

	keywords := os.Args[2:]
	bytes, err := os.ReadFile("comics.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read comics.json")
		os.Exit(-1)
	}

	var count int
	var item []*Comic

	if err := json.Unmarshal(bytes, &item); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to unmarshal data")
	}

	for _, c := range item {
		for _, k := range keywords {
			if !contains(c.Title, k) && !contains(c.Transcript, k) {
				continue
			}
			fmt.Printf("%d. Title: %s | Date: %4s/%02s/%02s | URL: %s\n", count+1, c.Title, c.Year, c.Month, c.Day, c.ImageURL)
			count++
			break
		}
	}
	fmt.Printf("Found %d comics!\n", count)
}

func contains(v string, k string) bool {
	return strings.Contains(v, k)
}
