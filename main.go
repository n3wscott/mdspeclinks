package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("./example.md")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lt := new(levelTracker)

	var tags []string

	scanner := bufio.NewScanner(file)
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
				tag := lt.Tag(word)

				i := strings.Index(l, word)

				updated.WriteString(l[:i])
				updated.WriteString(fmt.Sprintf(`<a name="%s"></a>%s`, tag, word))
				tags = append(tags, tag)

				l = l[i+len(word):]
			} else {
				updated.WriteString(l)
				break
			}
		}

		fmt.Println(updated.String())
	}

	fmt.Println("<!---")
	for _, tag := range tags {
		fmt.Printf("- [%s](#%s)\n", tag, tag)
	}
	fmt.Println("---!>")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// TODO: there is an edge case where the MUST/SHOULD NOT is split on two lines.
func hasSpecWord(line string) (string, bool) {
	for _, word := range []string{"MUST NOT", "MUST", "SHOULD NOT", "SHOULD", "MAY"} {
		if strings.Contains(line, word) {
			// TODO: this will not work with lines with more than one spec word... will have to make tokens and scan them.
			return word, true
		}
	}
	return "", false
}

type levelTracker struct {
	state    []int
	prefixes int
}

func (t *levelTracker) Tag(prefix string) string {
	var state []string
	for _, s := range t.state {
		state = append(state, strconv.Itoa(s))
	}
	t.prefixes++
	prefix = strings.ReplaceAll(prefix, " ", "_")
	tag := fmt.Sprintf("%s-%s-%s", prefix, strings.Join(state, "."), strconv.Itoa(t.prefixes))
	return strings.ToLower(tag)
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
