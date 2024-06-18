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
	data := strings.Split(string(bytes), "\n")
	for _, d := range data {
		var item Comic

		if err := json.Unmarshal([]byte(d), &item); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to unmarshal data")
			continue
		}

		for _, k := range keywords {
			if !contains(item.Title, k) && !contains(item.Transcript, k) {
				continue
			}
      fmt.Printf("%d. URL: %s | Date: %4s/%2s/%2s | Title: %s\n", count+1, item.ImageURL, item.Year, item.Month, item.Day, item.Title)
      count++
			break
		}
	}
  fmt.Printf("Found %d comics!\n", count)
}

func contains(v string, k string) bool {
	return strings.Contains(v, k)
}
