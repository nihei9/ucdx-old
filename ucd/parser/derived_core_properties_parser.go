package parser

import (
	"io"

	"github.com/nihei9/ucdx/ucd/property"
)

// ParseDerivedCoreProperties parses the DerivedCoreProperties.txt.
func ParseDerivedCoreProperties(r io.Reader) (*property.DerivedCoreProperties, error) {
	props := map[property.PropertyName][]*property.CodePointRange{}
	p := newParser(r)
	for p.parse() {
		if len(p.fields) == 0 {
			continue
		}

		cp, err := p.fields[0].codePointRange()
		if err != nil {
			return nil, err
		}
		name, _ := p.fields[1].name()
		if name == property.PropNameAlphabetic || name == property.PropNameUppercase ||
			name == property.PropNameLowercase || name == property.PropNameIDStart ||
			name == property.PropNameIDContinue || name == property.PropNameXIDStart ||
			name == property.PropNameXIDContinue {
			props[name] = append(props[name], cp)
		}
	}
	if p.err != nil {
		return nil, p.err
	}

	return &property.DerivedCoreProperties{
		Entries: props,
	}, nil
}
