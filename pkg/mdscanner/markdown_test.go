package mdscanner

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSingle(t *testing.T) {
	md := `This is a MUST example. MAY be two things.`

	in := strings.NewReader(md)
	found, err := Markdown(in)

	if err != nil {
		t.Error(err)
	}

	wanted := []Found{
		{
			Offset:   10,
			Lines:    []int{1},
			Section:  "0",
			Word:     "MUST",
			Sentence: "This is a MUST example.",
		},
		{
			Index:    1,
			Lines:    []int{1},
			Section:  "0",
			Word:     "MAY",
			Sentence: "MAY be two things.",
		},
	}

	if got := found; !cmp.Equal(got, wanted) {
		t.Error("Found (-want, +got):", cmp.Diff(wanted, got))
	}

	want := "This is a **MUST** example."
	if got := found[0].WhichWord(); !cmp.Equal(got, want) {
		t.Error("WhichWord (-want, +got):", cmp.Diff(want, got))
	}

	want = "**MAY** be two things."
	if got := found[1].WhichWord(); !cmp.Equal(got, want) {
		t.Error("WhichWord (-want, +got):", cmp.Diff(want, got))
	}
}

func TestReadme(t *testing.T) {
	md := `# Lorem markdownum optari illum

Temptatae atque usque MUST maerens moribundo Cererem. Cervix MAY et ut oculos iuveni
sublime dabit cera, **monstraverat animique**. Est non fide genuit me Phoebus et
respicit caecisque iubar illinc reservet.

## Agitasse ubi non profugus movent

- Neve atque MUST de heros
- Concedite SHOULD emisit
- MAY Tactae honorem multos

### Munychiosque ne

- Ab agat Caesar consiliis MUST crimine inquit
- Clipei SHOULD qui gemino dominus si habebat SHOULD NOT subiecta

## Qui Diti veniebat rursus

Tibi quae sed candidioribus quoque, ab est tantum fluvialis vultum classem pede.
Instruit exhausto exosus Amor **causas** amore ut MAY orbi potest rasa lunam
militiae *illum adhuc remisit* creatis.`

	in := strings.NewReader(md)
	found, err := Markdown(in)

	if err != nil {
		t.Error(err)
	}

	wanted := []Found{
		{
			Offset:   22,
			Lines:    []int{3},
			Section:  "1",
			Word:     "MUST",
			Sentence: "Temptatae atque usque MUST maerens moribundo Cererem.",
		},
		{
			Offset:   7,
			Index:    1,
			Lines:    []int{3, 4},
			Section:  "1",
			Word:     "MAY",
			Sentence: "Cervix MAY et ut oculos iuveni sublime dabit cera, **monstraverat animique**.",
		},
		{
			Offset:   13,
			Lines:    []int{9, 10},
			Section:  "1.1",
			Word:     "MUST",
			Sentence: "- Neve atque MUST de heros",
		},
		{
			Offset:   12,
			Index:    1,
			Lines:    []int{11},
			Section:  "1.1",
			Word:     "SHOULD",
			Sentence: "- Concedite SHOULD emisit",
		},
		{
			Offset:   2,
			Index:    2,
			Lines:    []int{12},
			Section:  "1.1",
			Word:     "MAY",
			Sentence: "- MAY Tactae honorem multos",
		},
		{
			Offset:   27,
			Lines:    []int{15, 16},
			Section:  "1.1.1",
			Word:     "MUST",
			Sentence: "- Ab agat Caesar consiliis MUST crimine inquit",
		},
		{
			Offset:   9,
			Index:    1,
			Lines:    []int{17},
			Section:  "1.1.1",
			Word:     "SHOULD",
			Sentence: "- Clipei SHOULD qui gemino dominus si habebat SHOULD NOT subiecta",
		},
		{
			Offset:   46,
			Index:    2,
			Lines:    []int{17},
			Section:  "1.1.1",
			Word:     "SHOULD NOT",
			Sentence: "- Clipei SHOULD qui gemino dominus si habebat SHOULD NOT subiecta",
		},
		{
			Offset:   11,
			Lines:    []int{21, 22},
			Section:  "1.2",
			Word:     "MAY",
			Sentence: "* amore ut MAY orbi potest rasa lunam militiae *illum adhuc remisit",
		},
	}

	if got := found; !cmp.Equal(got, wanted) {
		t.Error("Get (-want, +got):", cmp.Diff(wanted, got))
	}
}

func TestWithHeaders(t *testing.T) {
	md := `Before sections MUST work.
# Section 1
MUST work here too. RECOMMENDED if there are two.
## Subsection 1.1
This is a MUST example.
This is a MUST NOT example.
This is a REQUIRED example.
This is a SHOULD example.
This is a SHOULD NOT example.
This is a MAY example.
This is a MAY example.
This is a RECOMMENDED example.
This is a NOT RECOMMENDED example.
`

	in := strings.NewReader(md)
	found, err := Markdown(in)
	if err != nil {
		t.Error(err)
	}

	wanted := []Found{
		{
			Offset:   16,
			Lines:    []int{1},
			Section:  "0",
			Word:     "MUST",
			Sentence: "Before sections MUST work.",
		},
		{Lines: []int{3}, Section: "1", Word: "MUST", Sentence: "MUST work here too."},
		{
			Index:    1,
			Lines:    []int{3},
			Section:  "1",
			Word:     "RECOMMENDED",
			Sentence: "RECOMMENDED if there are two.",
		},
		{
			Offset:   10,
			Lines:    []int{5},
			Section:  "1.1",
			Word:     "MUST",
			Sentence: "This is a MUST example.",
		},
		{
			Offset:   10,
			Index:    1,
			Lines:    []int{6},
			Section:  "1.1",
			Word:     "MUST NOT",
			Sentence: "This is a MUST NOT example.",
		},
		{
			Offset:   10,
			Index:    2,
			Lines:    []int{7},
			Section:  "1.1",
			Word:     "REQUIRED",
			Sentence: "This is a REQUIRED example.",
		},
		{
			Offset:   10,
			Index:    3,
			Lines:    []int{8},
			Section:  "1.1",
			Word:     "SHOULD",
			Sentence: "This is a SHOULD example.",
		},
		{
			Offset:   10,
			Index:    4,
			Lines:    []int{9},
			Section:  "1.1",
			Word:     "SHOULD NOT",
			Sentence: "This is a SHOULD NOT example.",
		},
		{
			Offset:   10,
			Index:    5,
			Lines:    []int{10},
			Section:  "1.1",
			Word:     "MAY",
			Sentence: "This is a MAY example.",
		},
		{
			Offset:   10,
			Index:    6,
			Lines:    []int{11},
			Section:  "1.1",
			Word:     "MAY",
			Sentence: "This is a MAY example.",
		},
		{
			Offset:   10,
			Index:    7,
			Lines:    []int{12},
			Section:  "1.1",
			Word:     "RECOMMENDED",
			Sentence: "This is a RECOMMENDED example.",
		},
		{
			Offset:   10,
			Index:    8,
			Lines:    []int{13},
			Section:  "1.1",
			Word:     "NOT RECOMMENDED",
			Sentence: "This is a NOT RECOMMENDED example.",
		},
	}

	if got := found; !cmp.Equal(got, wanted) {
		t.Error("Get (-want, +got):", cmp.Diff(wanted, got))
	}
}
