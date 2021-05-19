package mdscanner

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// https://github.com/cloudevents/spec/blame/v1.0.1/spec.md#L13-L17

type Found struct {
	Line     int
	Column   int
	Section  string
	Word     string
	Context  string
	Sentence string
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
	var pendingFound []Found

	st := new(sentenceTracker)

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

				if sentences := st.Ingest(c, l[:i+len(word)]); len(sentences) > 0 {
					if pendingFound != nil {
						for i := range pendingFound {
							pendingFound[i].Sentence = sentences[len(sentences)-1]
						}
						found = append(found, pendingFound...)
						pendingFound = nil
					}
				}

				c += i
				l = l[i+len(word):]

				pendingFound = append(pendingFound, Found{
					Line:    at,
					Column:  c,
					Word:    word,
					Context: line,
					Section: lt.Section(),
				})

				c += len(word)
			} else {
				if sentences := st.Ingest(c, l); len(sentences) > 0 {
					if pendingFound != nil {
						for i := range pendingFound {
							pendingFound[i].Sentence = sentences[len(sentences)-1]
						}
						found = append(found, pendingFound...)
						pendingFound = nil
					}
				}
				break
			}
		}
	}

	if sentences := st.Ingest(0, "# EOF"); len(sentences) > 0 {
		if pendingFound != nil {
			for i := range pendingFound {
				pendingFound[i].Sentence = sentences[len(sentences)-1]
			}
			found = append(found, pendingFound...)
			pendingFound = nil
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

type sentenceTracker struct {
	text *strings.Builder
}

func (t *sentenceTracker) Ingest(offset int, line string) []string {
	line = strings.TrimSpace(line)
	if t.text == nil {
		t.text = new(strings.Builder)
	}

	var lines []string

	// Special case a blank line.
	if offset == 0 && len(line) == 0 {
		lines = append(lines, strings.TrimSpace(t.text.String()))
		t.text = new(strings.Builder)
		return lines
	} else if len(line) == 0 {
		return nil
	}

	for i := 0; i < len(line); i++ {
		if kind, ending := isEnding(offset, i, line); ending {
			switch kind {
			case "ignore":
				return lines
			case "pre":
				lines = append(lines, strings.TrimSpace(t.text.String()))
				t.text = new(strings.Builder)
				t.text.WriteByte(line[i])
			case "post":
				t.text.WriteByte(line[i])
				lines = append(lines, strings.TrimSpace(t.text.String()))
				t.text = new(strings.Builder)
			}
		} else {
			t.text.WriteByte(line[i])
		}
	}
	// Last, add space to let removing lines not stick words together.
	t.text.WriteByte(' ')
	return lines
}

// pre ending means this char terms the previous sentence, and a new one should start.
func isEnding(offset, at int, line string) (string, bool) {
	switch line[at] {
	case '#':
		if offset == 0 && at == 0 {
			return "ignore", true // only on the first line.
		}
	case '-', '*':
		return "pre", offset == 0 && len(line) > at+1 && line[at+1] == ' ' // only on the first line.
	case '.', '!', '?', ';':
		if len(line)-1 == at {
			return "post", true // line ends
		}
		return "post", line[at+1] == ' ' // next char is a space.
	}
	return "", false
}

type levelTracker struct {
	state []int
}

func (t *levelTracker) Section() string {
	var state []string
	for _, s := range t.state {
		state = append(state, strconv.Itoa(s))
	}
	if len(state) == 0 {
		state = append(state, "0")
	}
	return strings.Join(state, ".")
}

func (t *levelTracker) Next(level int) {
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
