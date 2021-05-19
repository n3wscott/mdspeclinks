package yaml

import (
	"crypto/md5"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/n3wscott/mdspeclinks/pkg/mdscanner"
	"gopkg.in/yaml.v3"
)

type SpecRequirements struct {
	Specification string        `json:"specification"`
	Processed     string        `json:"Processed"`
	Requirements  []Requirement `json:"requirements"`
}

type Requirement struct {
	ID      string `json:"id"`
	Word    string `json:"word"`
	Lines   string `json:"lines"`
	Section string `json:"section"`
	Offset  int    `json:"offset"`
	Text    string `json:"text"`
	Link    string `json:"link"`
	MD5     string `json:"md5"`
}

func Generate(file string, found []mdscanner.Found, out io.Writer) error {
	sr := new(SpecRequirements)
	sr.Specification = file
	sr.Processed = time.Now().Format(time.RFC3339)

	_, name := path.Split(file)
	name = strings.TrimSuffix(name, ".md")
	name = fmt.Sprintf("%s%s", strings.ToUpper(string(name[0])), name[1:])

	for _, f := range found {
		h := md5.Sum([]byte(f.Sentence))
		r := Requirement{
			ID:      fmt.Sprintf("%s%s%sC%d", name, strings.ReplaceAll(f.Word, " ", "_"), f.Section, f.Offset),
			Word:    f.Word,
			Lines:   f.Line(),
			Offset:  f.Offset,
			Section: f.Section,
			Text:    f.Sentence,
			Link:    f.BlameLink(toBlame(file)),
			MD5:     fmt.Sprintf("%x", h),
		}
		sr.Requirements = append(sr.Requirements, r)
	}

	encoder := yaml.NewEncoder(out)
	encoder.SetIndent(2)

	return encoder.Encode(sr)
}

func Unmarshal(in []byte) (*SpecRequirements, error) {
	sr := new(SpecRequirements)
	err := yaml.Unmarshal(in, sr)
	return sr, err
}

func toBlame(file string) string {
	return strings.Replace(file, "/blob/", "/blame/", 1)
}
