package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
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
	var wg sync.WaitGroup
	var mu sync.Mutex

	fmt.Println("Reading comics...")

	for {
		if errCount >= 2 || count >= 5000 {
			break
		}
		count++

		wg.Add(1)
		url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", count)
		go func(url string, count int) {
			fmt.Printf("Reading comic #%d\n", count)

			defer wg.Done()
			resp, err := http.Get(url)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Skipping #%d | Status: %d\n", count, resp.StatusCode)
				mu.Lock()
				errCount++
				mu.Unlock()
				return
			}

			defer resp.Body.Close()
			var response *Comic
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				fmt.Fprintf(os.Stderr, "Skipping #%d | Failed while decoding\n", count)
				mu.Lock()
				errCount++
				mu.Unlock()
				return
			}

			mu.Lock()
			comics = append(comics, response)
			mu.Unlock()
		}(url, count)
	}

	wg.Wait()
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
