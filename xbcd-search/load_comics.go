package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Comic struct {
	Num        int    `json:"num"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
	Year       string `json:"year"`
	Month      string `json:"month"`
	Day        string `json:"day"`
	ImageURL   string `json:"img"`
}

func LoadComics() {
  var fileName string
	if len(os.Args) <= 2 {
		fileName = "comics.json"
	} else {
		fileName = os.Args[2]
	}

	var count, errCount int
	var comics []*Comic

	fmt.Println("Reading comics...")

	for {
		if errCount >= 2 {
			break
		}

		count++
		url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", count)
		resp, err := http.Get(url)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Skipping #%d | Status: %d\n", count, resp.StatusCode)
			errCount++
			continue
		}

		defer resp.Body.Close()

		var response *Comic
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Fprintf(os.Stderr, "Skipping #%d | Failed while decoding\n", count)
			errCount++
			continue
		}

		comics = append(comics, response)
	}

	bytes, err := json.Marshal(&comics)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed while unmarshaling")
		return
	}

	if err := os.WriteFile(fileName, bytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("Successfully loaded %d comics!\n", len(comics))
}
