package parser

import (
	"io"

	"github.com/nihei9/ucdx/ucd/property"
)

// ParsePropertyAliases parses the PropertyAliases.txt.
func ParsePropertyAliases(r io.Reader) (*property.PropertyAliases, error) {
	aliases := []*property.PropertyAlias{}

	p := newParser(r)
	for p.parse() {
		if len(p.fields) == 0 {
			continue
		}

		abb, _ := p.fields[0].name()
		long, _ := p.fields[1].name()
		var others []property.PropertyName
		if len(p.fields) > 2 {
			for _, f := range p.fields[2:] {
				o, _ := f.name()
				others = append(others, o)
			}
		}
		aliases = append(aliases, &property.PropertyAlias{
			Abb:    abb,
			Long:   long,
			Others: others,
		})
	}
	if p.err != nil {
		return nil, p.err
	}

	return &property.PropertyAliases{
		Aliases: aliases,
	}, nil
}
