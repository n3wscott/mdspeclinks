package mdscanner

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strings"
)

// https://github.com/cloudevents/spec/blame/v1.0.1/spec.md#L13-L17

type Found struct {
	Line    int
	Column  int
	Word    string
	Context string
}

func (f Found) BlameLink(link string) string {
	return fmt.Sprintf("%s?w=%s&c=%d#L%d", link, f.Word, f.Column, f.Line)
}

func (f Found) WhichWord() string {
	return fmt.Sprintf("%s**%s**%s", f.Context[:f.Column], f.Word, f.Context[f.Column+len(f.Word):])
}

func Markdown(in io.Reader) ([]Found, error) {
	lt := new(levelTracker)

	var found []Found

	at := 0
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		at++
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {

			parts := strings.Split(line, " ")
			lt.Next(len(parts[0]))
		}
		l := line
		c := 0
		for {
			if word, has := hasSpecWord(l); has {
				i := strings.Index(l, word)
				c += i
				l = l[i+len(word):]

				found = append(found, Found{
					Line:    at,
					Column:  c,
					Word:    word,
					Context: line,
				})

				c += len(word)
			} else {
				break
			}
		}

	}

	return found, scanner.Err()
}

// TODO: there is an edge case where the MUST/SHOULD NOT is split on two lines.
func hasSpecWord(line string) (string, bool) {
	found := ""
	foundAt := math.MaxInt32 // real big number.
	for _, word := range []string{"MUST NOT", "MUST", "REQUIRED", "SHOULD NOT", "SHOULD", "SHALL NOT", "SHALL", "NOT RECOMMENDED", "RECOMMENDED", "MAY"} {
		if strings.Contains(line, word) {
			if at := strings.Index(line, word); at < foundAt {
				found = word
				foundAt = at
				continue
			}
			// TODO: this will not work with lines with more than one spec word... will have to make tokens and scan them.
		}
	}
	return found, foundAt < math.MaxInt32
}

type levelTracker struct {
	state    []int
	prefixes int
}

func (t *levelTracker) Next(level int) {
	t.prefixes = 0
	if len(t.state) < level {
		for len(t.state) != level {
			t.state = append(t.state, 0)
		}
	} else if len(t.state) > level {
		for len(t.state) != level {
			t.state = t.state[:len(t.state)-1]
		}
	}
	t.state[len(t.state)-1]++
}
