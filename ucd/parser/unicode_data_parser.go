package parser

import "io"

type UnicodeData struct {
	GeneralCategory map[string][]*CodePointRange `json:"general_category"`
}

func (u *UnicodeData) addGC(gc string, cp *CodePointRange) {
	// Section 4.2.11 Empty Fields in [UAX44]:
	// > The data file UnicodeData.txt defines many property values in each record. When a field in a data line
	// > for a code point is empty, that indicates that the property takes the default value for that code point.
	if gc == "" {
		return
	}

	cps, ok := u.GeneralCategory[gc]
	if ok {
		cpFrom, cpTo := cp.Range()
		c := cps[len(cps)-1]
		cFrom, cTo := c.Range()
		if cpFrom-cTo == 1 {
			c.Rewrite(cFrom, cpTo)
		} else {
			u.GeneralCategory[gc] = append(cps, cp)
		}
	} else {
		u.GeneralCategory[gc] = []*CodePointRange{cp}
	}
}

// ParseUnicodeData parses the UnicodeData.txt.
func ParseUnicodeData(r io.Reader) (*UnicodeData, error) {
	ud := &UnicodeData{
		GeneralCategory: map[string][]*CodePointRange{},
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
		gc := p.fields[2].normalizedSymbol()
		ud.addGC(gc, cp)
	}
	if p.err != nil {
		return nil, p.err
	}

	return ud, nil
}
