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
		{Line: 1, Column: 10, Word: "MUST", Context: "This is a MUST example. MAY be two things."},
		{Line: 1, Column: 24, Word: "MAY", Context: "This is a MUST example. MAY be two things."},
	}

	if got := found; !cmp.Equal(got, wanted) {
		t.Error("Found (-want, +got):", cmp.Diff(wanted, got))
	}

	want := "This is a **MUST** example. MAY be two things."
	if got := found[0].WhichWord(); !cmp.Equal(got, want) {
		t.Error("WhichWord (-want, +got):", cmp.Diff(want, got))
	}

	want = "This is a MUST example. **MAY** be two things."
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

	wanted := []Found{{
		Line: 3, Column: 22,
		Word:    "MUST",
		Context: "Temptatae atque usque MUST maerens moribundo Cererem. Cervix MAY et ut oculos iuveni",
	}, {
		Line: 3, Column: 61,
		Word:    "MAY",
		Context: "Temptatae atque usque MUST maerens moribundo Cererem. Cervix MAY et ut oculos iuveni",
	}, {
		Line: 9, Column: 13,
		Word:    "MUST",
		Context: "- Neve atque MUST de heros",
	}, {
		Line: 10, Column: 12, Word: "SHOULD", Context: "- Concedite SHOULD emisit",
	}, {
		Line: 11, Column: 2,
		Word:    "MAY",
		Context: "- MAY Tactae honorem multos",
	}, {
		Line: 15, Column: 27,
		Word:    "MUST",
		Context: "- Ab agat Caesar consiliis MUST crimine inquit",
	}, {
		Line: 16, Column: 9,
		Word:    "SHOULD",
		Context: "- Clipei SHOULD qui gemino dominus si habebat SHOULD NOT subiecta",
	}, {
		Line: 16, Column: 46,
		Word:    "SHOULD NOT",
		Context: "- Clipei SHOULD qui gemino dominus si habebat SHOULD NOT subiecta",
	}, {
		Line: 21, Column: 50,
		Word:    "MAY",
		Context: "Instruit exhausto exosus Amor **causas** amore ut MAY orbi potest rasa lunam",
	}}

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

	wanted := []Found{{
		Line: 1, Column: 16,
		Word: "MUST", Context: "Before sections MUST work.",
	}, {
		Line: 3, Column: 0,
		Word:    "MUST",
		Context: "MUST work here too. RECOMMENDED if there are two.",
	}, {
		Line: 3, Column: 20,
		Word:    "RECOMMENDED",
		Context: "MUST work here too. RECOMMENDED if there are two.",
	}, {
		Line: 5, Column: 10,
		Word:    "MUST",
		Context: "This is a MUST example.",
	}, {
		Line: 6, Column: 10,
		Word:    "MUST NOT",
		Context: "This is a MUST NOT example.",
	}, {
		Line: 7, Column: 10,
		Word:    "REQUIRED",
		Context: "This is a REQUIRED example.",
	}, {
		Line: 8, Column: 10,
		Word: "SHOULD", Context: "This is a SHOULD example.",
	}, {
		Line: 9, Column: 10,
		Word:    "SHOULD NOT",
		Context: "This is a SHOULD NOT example.",
	}, {
		Line: 10, Column: 10,
		Word: "MAY", Context: "This is a MAY example.",
	}, {
		Line: 11, Column: 10,
		Word: "MAY", Context: "This is a MAY example.",
	}, {
		Line: 12, Column: 10,
		Word:    "RECOMMENDED",
		Context: "This is a RECOMMENDED example.",
	}, {
		Line: 13, Column: 10,
		Word:    "NOT RECOMMENDED",
		Context: "This is a NOT RECOMMENDED example.",
	}}

	if got := found; !cmp.Equal(got, wanted) {
		t.Error("Get (-want, +got):", cmp.Diff(wanted, got))
	}
}
