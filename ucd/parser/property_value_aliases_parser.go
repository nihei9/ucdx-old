package parser

import (
	"io"

	"github.com/nihei9/ucdx/ucd/property"
)

// ParsePropertyValueAliases parses the PropertyValueAliases.txt.
func ParsePropertyValueAliases(r io.Reader) (*property.PropertyValueAliases, error) {
	aliases := map[string]*property.PropertyValueAliase{}
	defaultValues := map[string]*property.DefaultValue{}

	p := newParser(r)
	for p.parse() {
		// The format of the data file is explained in section 5.8.2 Property Value Aliases in [UAX44].
		if len(p.fields) > 0 {
			var others []string
			if len(p.fields) >= 4 {
				rest := p.fields[3:]
				others = make([]string, len(rest))
				for i, f := range rest {
					others[i] = f.normalizedSymbol()
				}
			}
			aliases[p.fields[0].normalizedSymbol()] = &property.PropertyValueAliase{
				Abb:    p.fields[1].normalizedSymbol(),
				Long:   p.fields[2].normalizedSymbol(),
				Others: others,
			}
		}

		// The format of the default values is explained in section 4.2.10 @missing Conventions in [UAX44].
		if len(p.defaultFields) > 0 && p.defaultFields[1].symbol() == "General_Category" {
			cp, err := p.defaultFields[0].codePointRange()
			if err != nil {
				return nil, err
			}
			prop := p.defaultFields[1].symbol()
			val := p.defaultFields[2].normalizedSymbol()

			defaultValues[prop] = &property.DefaultValue{
				Value: val,
				CP:    cp,
			}
		}
	}
	if p.err != nil {
		return nil, p.err
	}
	return &property.PropertyValueAliases{
		Aliases:       aliases,
		DefaultValues: defaultValues,
	}, nil
}
