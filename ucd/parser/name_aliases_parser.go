package parser

import (
	"io"
	"sort"
)

type NameAliasesEntry struct {
	CP      rune     `json:"cp"`
	Aliases []string `json:"aliases"`
}

type NameAliases struct {
	Entries []*NameAliasesEntry `json:"entries"`
}

// ParseNameAliases parses the NameAliases.txt.
func ParseNameAliases(r io.Reader) (*NameAliases, error) {
	aliases := map[rune][]string{}
	p := newParser(r)
	for p.parse() {
		if len(p.fields) == 0 {
			continue
		}

		cp, err := p.fields[0].codePointRange()
		if err != nil {
			return nil, err
		}
		na, ok := p.fields[1].name()
		if !ok {
			continue
		}
		c, _ := cp.Range()
		aliases[c] = append(aliases[c], na)
	}
	if p.err != nil {
		return nil, p.err
	}

	entries := make([]*NameAliasesEntry, 0, len(aliases))
	for c, as := range aliases {
		entries = append(entries, &NameAliasesEntry{
			CP:      c,
			Aliases: as,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CP < entries[j].CP
	})

	return &NameAliases{
		Entries: entries,
	}, nil
}
