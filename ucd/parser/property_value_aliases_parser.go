package parser

import (
	"io"

	"github.com/nihei9/ucdx/ucd/property"
)

// ParsePropertyValueAliases parses the PropertyValueAliases.txt.
func ParsePropertyValueAliases(r io.Reader) (*property.PropertyValueAliases, error) {
	aliases := map[property.PropertyName][]*property.PropertyValueAliase{}
	defaultValues := map[property.PropertyName]*property.DefaultValue{}

	p := newParser(r)
	for p.parse() {
		// The format of the data file is explained in section 5.8.2 Property Value Aliases in [UAX44].
		if len(p.fields) > 0 {
			var others []property.PropertyValueSymbol
			if len(p.fields) >= 4 {
				rest := p.fields[3:]
				others = make([]property.PropertyValueSymbol, len(rest))
				for i, f := range rest {
					others[i] = f.normalizedSymbol()
				}
			}
			propName, _ := p.fields[0].name()
			aliases[propName] = append(aliases[propName], &property.PropertyValueAliase{
				Abb:    p.fields[1].normalizedSymbol(),
				Long:   p.fields[2].normalizedSymbol(),
				Others: others,
			})
		}

		// The format of the default values is explained in section 4.2.10 @missing Conventions in [UAX44].
		if len(p.defaultFields) > 0 {
			if propName, _ := p.defaultFields[1].name(); propName == property.PropNameGeneralCategory {
				cp, err := p.defaultFields[0].codePointRange()
				if err != nil {
					return nil, err
				}
				defaultValues[propName] = &property.DefaultValue{
					Value: p.defaultFields[2].normalizedSymbol(),
					CP:    cp,
				}
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
