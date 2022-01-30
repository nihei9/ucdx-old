package parser

import (
	"io"

	"github.com/nihei9/ucdx/ucd/property"
)

// ParsePropList parses the PropList.txt.
func ParsePropList(r io.Reader) (*property.PropList, error) {
	var ws []*property.CodePointRange
	p := newParser(r)
	for p.parse() {
		if len(p.fields) == 0 {
			continue
		}

		cp, err := p.fields[0].codePointRange()
		if err != nil {
			return nil, err
		}

		if propName, _ := p.fields[1].name(); propName == property.PropNameWhiteSpace {
			ws = append(ws, cp)
		}
	}
	if p.err != nil {
		return nil, p.err
	}

	return &property.PropList{
		WhiteSpace: ws,
	}, nil
}
