package parser

import "io"

type DerivedCoreProperties struct {
	Entries map[string][]*CodePointRange `json:"entries"`
}

// ParseDerivedCoreProperties parses the DerivedCoreProperties.txt.
func ParseDerivedCoreProperties(r io.Reader) (*DerivedCoreProperties, error) {
	props := map[string][]*CodePointRange{}
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
		if sym == "Alphabetic" || sym == "Uppercase" || sym == "Lowercase" {
			props[sym] = append(props[sym], cp)
		}
	}
	if p.err != nil {
		return nil, p.err
	}

	return &DerivedCoreProperties{
		Entries: props,
	}, nil
}
