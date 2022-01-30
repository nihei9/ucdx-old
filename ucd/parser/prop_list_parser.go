package parser

import (
	"io"

	"github.com/nihei9/ucdx/ucd/property"
)

// ParsePropList parses the PropList.txt.
func ParsePropList(r io.Reader) (*property.PropList, error) {
	var oa []*property.CodePointRange
	var ol []*property.CodePointRange
	var ou []*property.CodePointRange
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

		switch p.fields[1].symbol() {
		case "Other_Alphabetic":
			oa = append(oa, cp)
		case "Other_Lowercase":
			ol = append(ol, cp)
		case "Other_Uppercase":
			ou = append(ou, cp)
		case "White_Space":
			ws = append(ws, cp)
		}
	}
	if p.err != nil {
		return nil, p.err
	}

	return &property.PropList{
		OtherAlphabetic: oa,
		OtherLowercase:  ol,
		OtherUppercase:  ou,
		WhiteSpace:      ws,
	}, nil
}
