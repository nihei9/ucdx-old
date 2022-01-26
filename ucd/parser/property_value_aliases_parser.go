package parser

import "io"

type PropertyValueAliases struct {
	Aliases       map[string]*PropertyValueAliase `json:"aliases"`
	DefaultValues map[string]*DefaultValue        `json:"default_values"`
}

// PropertyValueAliase represents a set of aliases for a property value.
// `Abb` and `Long` are the preferred aliases.
type PropertyValueAliase struct {
	// Abb is an abbreviated symbolic name for a property value.
	Abb string `json:"abb"`

	// Long is the long symbolic name for a property value.
	Long string `json:"long"`

	// Others is a set of other aliases for a property value.
	Others []string `json:"others,omitempty"`
}

type DefaultValue struct {
	Value string          `json:"value"`
	CP    *CodePointRange `json:"cp"`
}

// ParsePropertyValueAliases parses the PropertyValueAliases.txt.
func ParsePropertyValueAliases(r io.Reader) (*PropertyValueAliases, error) {
	aliases := map[string]*PropertyValueAliase{}
	defaultValues := map[string]*DefaultValue{}

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
			aliases[p.fields[0].normalizedSymbol()] = &PropertyValueAliase{
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

			defaultValues[prop] = &DefaultValue{
				Value: val,
				CP:    cp,
			}
		}
	}
	if p.err != nil {
		return nil, p.err
	}
	return &PropertyValueAliases{
		Aliases:       aliases,
		DefaultValues: defaultValues,
	}, nil
}
