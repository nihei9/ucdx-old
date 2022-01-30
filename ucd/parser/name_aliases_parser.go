package parser

import (
	"io"
	"sort"

	"github.com/nihei9/ucdx/ucd/property"
)

// ParseNameAliases parses the NameAliases.txt.
func ParseNameAliases(r io.Reader) (*property.NameAliases, error) {
	aliases := map[rune][]property.PropertyName{}
	p := newParser(r)
	for p.parse() {
		if len(p.fields) == 0 {
			continue
		}

		cp, err := p.fields[0].codePointRange()
		if err != nil {
			return nil, err
		}
		name, ok := p.fields[1].name()
		if !ok {
			continue
		}
		c, _ := cp.Range()
		aliases[c] = append(aliases[c], name)
	}
	if p.err != nil {
		return nil, p.err
	}

	entries := make([]*property.NameAliasesEntry, 0, len(aliases))
	for c, as := range aliases {
		entries = append(entries, &property.NameAliasesEntry{
			CP:      c,
			Aliases: as,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CP < entries[j].CP
	})

	return &property.NameAliases{
		Entries: entries,
	}, nil
}
