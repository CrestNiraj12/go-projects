package main

import (
	"fmt"
	"os"
)

const (
	load   = "load_comics"
	search = "search_comics"
)

func main() {
  if len(os.Args) <= 1 {
    throwInvalidExecution()
	}

	cmd := os.Args[1]
	switch cmd {
	case load:
		LoadComics()
	case search:
		SearchComics()
	default:
		throwInvalidExecution()
	}
}

func throwInvalidExecution() {
  fmt.Fprintln(os.Stderr, ">> Usage: go run . [command] [...arguments]\nCommands:\nload_comics <json filename>\nsearch_comics <keyword1> <keyword2>...")
	os.Exit(0)
}
