package parser

import (
	"io"

	"github.com/nihei9/ucdx/ucd/property"
)

// ParseUnicodeData parses the UnicodeData.txt.
func ParseUnicodeData(r io.Reader) (*property.UnicodeData, error) {
	ud := &property.UnicodeData{
		Name:            map[property.PropertyName]*property.CodePointRange{},
		GeneralCategory: map[property.PropertyValueSymbol][]*property.CodePointRange{},
	}

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
		if ok {
			ud.Name[name] = cp
		}
		gc := p.fields[2].normalizedSymbol()
		ud.AddGC(gc, cp)
	}
	if p.err != nil {
		return nil, p.err
	}

	return ud, nil
}
