package parser

import (
	"fmt"
	"io"

	"github.com/nihei9/ucdx/ucd/property"
)

// ParseUnicodeData parses the UnicodeData.txt.
func ParseUnicodeData(r io.Reader) (*property.UnicodeData, error) {
	ud := &property.UnicodeData{
		Name:            map[property.PropertyName]*property.CodePointRange{},
		GeneralCategory: map[property.PropertyValueSymbol][]*property.CodePointRange{},
	}

	inRange := false
	var firstCP rune
	var firstGC property.PropertyValueSymbol
	p := newParser(r)
	for p.parse() {
		if len(p.fields) == 0 {
			continue
		}

		cp, err := p.fields[0].codePointRange()
		if err != nil {
			return nil, err
		}

		if inRange {
			if !p.fields[1].rangeLast() {
				// TODO: ERROR
				return nil, fmt.Errorf("")
			}
			lastCP, _ := cp.Range()
			ud.AddGC(firstGC, property.NewCodePointRange(firstCP, lastCP))
			inRange = false

			continue
		}
		name, ok := p.fields[1].name()
		if ok {
			ud.Name[name] = cp
		} else {
			if p.fields[1].rangeStart() {
				inRange = true
				firstCP, _ = cp.Range()
			}
		}

		gc := p.fields[2].normalizedSymbol()
		if inRange {
			firstGC = gc
		} else {
			ud.AddGC(gc, cp)
		}
	}
	if p.err != nil {
		return nil, p.err
	}

	return ud, nil
}
