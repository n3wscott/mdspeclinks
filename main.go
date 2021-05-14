package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/n3wscott/mdspeclinks/pkg/mdscanner"
)

func main() {
	var filename string
	if len(os.Args) == 2 && strings.HasPrefix(os.Args[1], "https://github.com") {
		filename = os.Args[1]
	} else {
		fmt.Fprintf(os.Stderr, "Please provide a github blob url to the spec language markdown file.\nUsage:\n\tmdspeclinks https://github.com/org/repo/blob/ref/path/file.md\n")
		os.Exit(1)
	}

	//file, err := os.Open(filename)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Failed to open file %q: %v\n", filename, err)
	//	os.Exit(1)
	//}
	//defer file.Close()

	resp, err := http.Get(toRaw(filename))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch file %q: %v\n", filename, err)
		os.Exit(1)
	}

	if found, err := mdscanner.Markdown(resp.Body); err != nil {
		fmt.Fprintf(os.Stderr, "Failed processing markdown: %v", err)
	} else {
		for _, f := range found {
			fmt.Fprintf(os.Stdout, "%s - %q\n", f.BlameLink(toBlame(filename)), f.WhichWord())
		}
	}
}

func toRaw(file string) string {
	hasBlob := strings.Replace(file, "https://github.com", "https://raw.githubusercontent.com", 1)
	return strings.Replace(hasBlob, "/blob/", "/", 1)
}

func toBlame(file string) string {
	return strings.Replace(file, "/blob/", "/blame/", 1)
}

//                https://github.com/n3wscott/mdspeclinks/blob/main/example.md
// https://github.com/n3wscott/mdspeclinks/blame/main/example.md#L17
// https://raw.githubusercontent.com/n3wscott/mdspeclinks/main/example.md
