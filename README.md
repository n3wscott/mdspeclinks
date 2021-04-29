# mdspeclinks

A tool to inject spec language links into markdown. `mdspeclinks` will look for
"MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT",
"RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" (see
[RFC 2119](https://tools.ietf.org/html/rfc2119) for interpretation) and attach a
link to the spec word to be linked from external sources.

## Installation

`mdspeclinks` can be installed and upgraded by running:

```shell
go get github.com/n3wscott/mdspeclinks
```

## Usage

```bash
mdspeclinks path/to/file
```

This will read `path/to/file` and output the updated document to standard out.
Chain with append on the shell:

```bash
mdspeclinks path/to/file > path/to/new-file
```

### Example

Given the following markdown:

```
# Lorem markdownum optari illum

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
militiae *illum adhuc remisit* creatis.
```

The tool would output:

# Lorem markdownum optari illum

Temptatae atque usque <a name="must-1-1"></a>MUST<sup>[1-1](#must-1-1)</sup>
maerens moribundo Cererem. Cervix
<a name="may-1-2"></a>MAY<sup>[1-2](#may-1-2)</sup> et ut oculos iuveni sublime
dabit cera, **monstraverat animique**. Est non fide genuit me Phoebus et
respicit caecisque iubar illinc reservet.

## Agitasse ubi non profugus movent

- Neve atque <a name="must-1.1-1"></a>MUST<sup>[1.1-1](#must-1.1-1)</sup> de
  heros
- Concedite <a name="should-1.1-2"></a>SHOULD<sup>[1.1-2](#should-1.1-2)</sup>
  emisit
- <a name="may-1.1-3"></a>MAY<sup>[1.1-3](#may-1.1-3)</sup> Tactae honorem
  multos

### Munychiosque ne

- Ab agat Caesar consiliis
  <a name="must-1.1.1-1"></a>MUST<sup>[1.1.1-1](#must-1.1.1-1)</sup> crimine
  inquit
- Clipei
  <a name="should-1.1.1-2"></a>SHOULD<sup>[1.1.1-2](#should-1.1.1-2)</sup> qui
  gemino dominus si habebat <a name="should_not-1.1.1-3"></a>SHOULD
  NOT<sup>[1.1.1-3](#should_not-1.1.1-3)</sup> subiecta

## Qui Diti veniebat rursus

Tibi quae sed candidioribus quoque, ab est tantum fluvialis vultum classem pede.
Instruit exhausto exosus Amor **causas** amore ut
<a name="may-1.2-1"></a>MAY<sup>[1.2-1](#may-1.2-1)</sup> orbi potest rasa lunam
militiae _illum adhuc remisit_ creatis.
