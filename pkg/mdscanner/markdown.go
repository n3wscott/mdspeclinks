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
	Offset   int
	Index    int
	Lines    []int
	Section  string
	Word     string
	Sentence string
}

func (f Found) BlameLink(link string) string {
	return fmt.Sprintf("%s?w=%s&c=%d#%s", link, strings.ReplaceAll(f.Word, " ", "_"), f.Offset, f.Line())
}

func (f Found) Line() string {
	if len(f.Lines) == 0 {
		return ""
	}
	if len(f.Lines) == 1 {
		return fmt.Sprintf("L%d", f.Lines[0])
	} else {
		return fmt.Sprintf("L%d-L%d", f.Lines[0], f.Lines[len(f.Lines)-1])
	}
}

func (f Found) WhichWord() string {
	return fmt.Sprintf("%s**%s**%s", f.Sentence[:f.Offset], f.Word, f.Sentence[f.Offset+len(f.Word):])
}

func Markdown(in io.Reader) ([]Found, error) {
	lt := new(levelTracker)

	var found []Found
	st := new(sentenceTracker)
	lineCount := 0
	var lines []int
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		lineCount++
		lines = append(lines, lineCount)
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			parts := strings.Split(line, " ")
			lt.Next(len(parts[0]))
		}

		line = strings.TrimSpace(line)
		i := 0
		for {
			sentence, offset, term := st.IngestUntil(i, line)
			if sentence != "" {
				l := sentence
				o := 0
				for {
					word, has := hasSpecWord(l)
					if !has {
						break
					}
					at := strings.Index(l, word)
					o += at

					l = l[at+len(word):]

					section, index := lt.Section()
					found = append(found, Found{
						Offset:   o,
						Index:    index,
						Word:     word,
						Section:  section,
						Sentence: sentence,
						Lines:    lines,
					})

					o += len(word)
				}
			}
			if term {
				if offset == 0 || i+offset >= len(line) {
					lines = []int(nil)
				} else {
					lines = []int{lineCount}
				}
			}

			i += offset
			if i >= len(line) {
				break
			}
			line = line[offset:]
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

// IngestUntil reads line until a sentence terminates, and returns back the number of chars consumed from line.
func (t *sentenceTracker) IngestUntil(offset int, line string) (string /* sentence */, int /* consumed */, bool /* terminated */) {
	if t.text == nil {
		t.text = new(strings.Builder)
	}

	// Special case a blank line.
	if offset == 0 && len(line) == 0 {
		sentence := strings.TrimSpace(t.text.String())
		t.text = new(strings.Builder)
		return sentence, len(line), true
	} else if len(line) == 0 {
		return "", len(line), true
	}

	for i := 0; i < len(line); i++ {
		if kind, ending := isEnding(offset, i, line); ending {
			sentence := strings.TrimSpace(t.text.String())
			switch kind {
			case "ignore":
				t.text = new(strings.Builder)
				return sentence, len(line), true
			case "pre":
				t.text = new(strings.Builder)
				if sentence != "" {
					return sentence, i, true
				} else {
					t.text.WriteByte(line[i])
				}
			case "post":
				t.text.WriteByte(line[i])
				sentence = strings.TrimSpace(t.text.String())
				t.text = new(strings.Builder)
				return sentence, i + 1, true
			}
		} else {
			t.text.WriteByte(line[i])
		}
	}
	// Last, add space to let removing lines not stick words together.
	t.text.WriteByte(' ')
	return "", len(line), false
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
	index int
}

func (t *levelTracker) Section() (string, int) {
	index := t.index
	t.index++
	var state []string
	for _, s := range t.state {
		state = append(state, strconv.Itoa(s))
	}
	if len(state) == 0 {
		state = append(state, "0")
	}
	return strings.Join(state, "."), index
}

func (t *levelTracker) Next(level int) {
	t.index = 0
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
