package mdscanner

import (
	"fmt"
	"os"
	"strings"
)

func Example_single() {
	md := `This is a MUST example.`

	in := strings.NewReader(md)
	if err := Markdown(in, os.Stdout); err != nil {
		fmt.Println(err)
	}

	// Output:
	// This is a <a name="must-0-1"></a>MUST<sup>[0-1](#must-0-1)</sup> example.
}

func Example_readme() {
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
	if err := Markdown(in, os.Stdout); err != nil {
		fmt.Println(err)
	}

	// Output:
	// # Lorem markdownum optari illum
	//
	// Temptatae atque usque <a name="must-1-1"></a>MUST<sup>[1-1](#must-1-1)</sup> maerens moribundo Cererem. Cervix <a name="may-1-2"></a>MAY<sup>[1-2](#may-1-2)</sup> et ut oculos iuveni
	// sublime dabit cera, **monstraverat animique**. Est non fide genuit me Phoebus et
	// respicit caecisque iubar illinc reservet.
	//
	// ## Agitasse ubi non profugus movent
	//
	// - Neve atque <a name="must-1.1-1"></a>MUST<sup>[1.1-1](#must-1.1-1)</sup> de heros
	// - Concedite <a name="should-1.1-2"></a>SHOULD<sup>[1.1-2](#should-1.1-2)</sup> emisit
	// - <a name="may-1.1-3"></a>MAY<sup>[1.1-3](#may-1.1-3)</sup> Tactae honorem multos
	//
	// ### Munychiosque ne
	//
	// - Ab agat Caesar consiliis <a name="must-1.1.1-1"></a>MUST<sup>[1.1.1-1](#must-1.1.1-1)</sup> crimine inquit
	// - Clipei <a name="should-1.1.1-2"></a>SHOULD<sup>[1.1.1-2](#should-1.1.1-2)</sup> qui gemino dominus si habebat <a name="should_not-1.1.1-3"></a>SHOULD NOT<sup>[1.1.1-3](#should_not-1.1.1-3)</sup> subiecta
	//
	// ## Qui Diti veniebat rursus
	//
	// Tibi quae sed candidioribus quoque, ab est tantum fluvialis vultum classem pede.
	// Instruit exhausto exosus Amor **causas** amore ut <a name="may-1.2-1"></a>MAY<sup>[1.2-1](#may-1.2-1)</sup> orbi potest rasa lunam
	// militiae *illum adhuc remisit* creatis.
}

func Example_withHeaders() {
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
	if err := Markdown(in, os.Stdout); err != nil {
		fmt.Println(err)
	}

	// Output:
	// Before sections <a name="must-0-1"></a>MUST<sup>[0-1](#must-0-1)</sup> work.
	// # Section 1
	// <a name="must-1-1"></a>MUST<sup>[1-1](#must-1-1)</sup> work here too. <a name="recommended-1-2"></a>RECOMMENDED<sup>[1-2](#recommended-1-2)</sup> if there are two.
	// # Section 1
	// ## Subsection 1.1
	// This is a <a name="must-1.1-1"></a>MUST<sup>[1.1-1](#must-1.1-1)</sup> example.
	// This is a <a name="must_not-1.1-2"></a>MUST NOT<sup>[1.1-2](#must_not-1.1-2)</sup> example.
	// This is a <a name="required-1.1-3"></a>REQUIRED<sup>[1.1-3](#required-1.1-3)</sup> example.
	// This is a <a name="should-1.1-4"></a>SHOULD<sup>[1.1-4](#should-1.1-4)</sup> example.
	// This is a <a name="should_not-1.1-5"></a>SHOULD NOT<sup>[1.1-5](#should_not-1.1-5)</sup> example.
	// This is a <a name="may-1.1-6"></a>MAY<sup>[1.1-6](#may-1.1-6)</sup> example.
	// This is a <a name="may-1.1-7"></a>MAY<sup>[1.1-7](#may-1.1-7)</sup> example.
	// This is a <a name="recommended-1.1-8"></a>RECOMMENDED<sup>[1.1-8](#recommended-1.1-8)</sup> example.
	// This is a <a name="not_recommended-1.1-9"></a>NOT RECOMMENDED<sup>[1.1-9](#not_recommended-1.1-9)</sup> example.
}
