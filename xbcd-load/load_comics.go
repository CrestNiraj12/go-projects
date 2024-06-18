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
	if len(os.Args) <= 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run . load_comics <json filename>")
		os.Exit(-1)
	}

	fileName := os.Args[2]
	var count, errCount int
	var bytes []byte

	for {
		count++
		url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", count)
		resp, err := http.Get(url)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Skipping #%d | Status: %d\n", count, resp.StatusCode)
			errCount++
			if errCount < 5 {
				continue
			}
			break
		}

		defer resp.Body.Close()

		var response Comic
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Fprintf(os.Stderr, "Skipping #%d | Failed while decoding\n", count)
			errCount++
			if errCount < 5 {
				continue
			}
		  break	
		}

		mResp, err := json.Marshal(response)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Skipping #%d | Failed while marshaling\n", count)
			errCount++
			if errCount < 5 {
				continue
			}
		  break	
		}

		bytes = append(bytes, append(mResp, byte('\n'))...)
		fmt.Printf("Read %d comics!\n", count)
	}

	if err := os.WriteFile(fileName, bytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("Successfully loaded %d comics!\n", count)
}
