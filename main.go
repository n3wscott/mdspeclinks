package main

import (
	"fmt"
	"os"

	"github.com/n3wscott/mdspeclinks/pkg/mdscanner"
)

func main() {
	var filename string
	if len(os.Args) == 2 {
		filename = os.Args[1]
	} else {
		fmt.Fprintf(os.Stderr, "Please provide a file path to spec language markdown.\nUsage:\n\tmdspeclinks <path>\n")
		os.Exit(1)
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file %q: %v\n", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	if err := mdscanner.Markdown(file, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Failed processing markdown: %v", err)
	}
}
