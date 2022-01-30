package parser

import (
	"io"

	"github.com/nihei9/ucdx/ucd/property"
)

// ParseDerivedCoreProperties parses the DerivedCoreProperties.txt.
func ParseDerivedCoreProperties(r io.Reader) (*property.DerivedCoreProperties, error) {
	props := map[string][]*property.CodePointRange{}
	p := newParser(r)
	for p.parse() {
		if len(p.fields) == 0 {
			continue
		}

		cp, err := p.fields[0].codePointRange()
		if err != nil {
			return nil, err
		}
		sym := p.fields[1].symbol()
		if sym == "Alphabetic" || sym == "Uppercase" || sym == "Lowercase" ||
			sym == "ID_Start" || sym == "ID_Continue" || sym == "XID_Start" || sym == "XID_Continue" {
			props[sym] = append(props[sym], cp)
		}
	}
	if p.err != nil {
		return nil, p.err
	}

	return &property.DerivedCoreProperties{
		Entries: props,
	}, nil
}
