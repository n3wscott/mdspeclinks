package mdscanner

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

func Markdown(in io.Reader, out io.Writer) error {
	lt := new(levelTracker)

	var tags []string

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {

			parts := strings.Split(line, " ")
			lt.Next(len(parts[0]))
		}
		updated := strings.Builder{}
		l := line
		for {
			if word, found := hasSpecWord(l); found {
				tag, short := lt.Tag(word)

				i := strings.Index(l, word)

				updated.WriteString(l[:i])
				updated.WriteString(fmt.Sprintf(`<a name="%s"></a>%s<sup>[%s](#%s)</sup>`, tag, word, short, tag))
				tags = append(tags, tag)

				l = l[i+len(word):]
			} else {
				updated.WriteString(l)
				break
			}
		}

		if _, err := fmt.Fprintln(out, updated.String()); err != nil {
			return err
		}
	}

	return scanner.Err()
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

func (t *levelTracker) Tag(prefix string) ( /* tag */ string /* short */, string) {
	var state []string
	for _, s := range t.state {
		state = append(state, strconv.Itoa(s))
	}
	if len(state) == 0 {
		state = append(state, "0")
	}
	t.prefixes++
	prefix = strings.ReplaceAll(prefix, " ", "_")
	tag := fmt.Sprintf("%s-%s-%s", prefix, strings.Join(state, "."), strconv.Itoa(t.prefixes))
	short := fmt.Sprintf("%s-%s", strings.Join(state, "."), strconv.Itoa(t.prefixes))
	return strings.ToLower(tag), short
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
