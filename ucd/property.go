package ucd

import (
	"fmt"
	"strings"

	"github.com/nihei9/ucdx/ucd/property"
)

type PropertySet struct {
	CP                    rune                                             `json:"code_point"`
	Properties            map[property.PropertyName]property.PropertyValue `json:"properties"`
	GeneralCategoryGroups []property.PropertyValueSymbol                   `json:"general_category_group"`
}

func (s *PropertySet) Lookup(propName property.PropertyName) *property.Property {
	v, ok := s.Properties[propName]
	if !ok {
		return nil
	}
	return property.NewProperty(propName, v)
}

type UCD struct {
	UnicodeData           *property.UnicodeData
	NameAliases           *property.NameAliases
	DerivedCoreProperties *property.DerivedCoreProperties
	PropertyAliases       *property.PropertyAliases
	PropertyValueAliases  *property.PropertyValueAliases
	PropList              *property.PropList
}

func (u *UCD) AnalizeCodePoint(c rune) *PropertySet {
	gc := u.lookupGeneralCategory(c)
	return &PropertySet{
		CP: c,
		Properties: map[property.PropertyName]property.PropertyValue{
			property.PropNameName:            u.lookupName(c),
			property.PropNameNameAlias:       u.lookupNameAlias(c),
			property.PropNameGeneralCategory: gc,
			property.PropNameAlphabetic:      u.isAlphabetic(c),
			property.PropNameUppercase:       u.isUppercase(c),
			property.PropNameLowercase:       u.isLowercase(c),
			property.PropNameIDStart:         u.isIDStart(c),
			property.PropNameIDContinue:      u.isIDContinue(c),
			property.PropNameXIDStart:        u.isXIDStart(c),
			property.PropNameXIDContinue:     u.isXIDContinue(c),
			property.PropNameWhiteSpace:      u.isWhiteSpace(c),
		},
		GeneralCategoryGroups: lookupGCGroups(gc),
	}
}

// 4.8 Table 4-8. Name Derivation Rule Prefix Strings in [Unicode].
var namePrefixes = map[string][]*property.CodePointRange{
	"HANGUL SYLLABLE ": {
		property.NewCodePointRange(0xAC00, 0xD7A3),
	},
	"CJK UNIFIED IDEOGRAPH-": {
		property.NewCodePointRange(0x3400, 0x4DBF),
		property.NewCodePointRange(0x4E00, 0x9FFC),
		property.NewCodePointRange(0x20000, 0x2A6DD),
		property.NewCodePointRange(0x2A700, 0x2B734),
		property.NewCodePointRange(0x2B740, 0x2B81D),
		property.NewCodePointRange(0x2B820, 0x2CEA1),
		property.NewCodePointRange(0x2CEB0, 0x2EBE0),
		property.NewCodePointRange(0x30000, 0x3134A),
	},
	"TANGUT IDEOGRAPH-": {
		property.NewCodePointRange(0x17000, 0x187F7),
		property.NewCodePointRange(0x18D00, 0x18D08),
	},
	"KHITAN SMALL SCRIPT CHARACTER-": {
		property.NewCodePointRange(0x18B00, 0x18CD5),
	},
	"NUSHU CHARACTER-": {
		property.NewCodePointRange(0x1B170, 0x1B2FB),
	},
	"CJK COMPATIBILITY IDEOGRAPH-": {
		property.NewCodePointRange(0xF900, 0xFA6D),
		property.NewCodePointRange(0xFA70, 0xFAD9),
		property.NewCodePointRange(0x2F800, 0x2FA1D),
	},
}

func (u *UCD) lookupName(c rune) property.PropertyName {
	for prefix, cps := range namePrefixes {
		for _, cp := range cps {
			if cp.Contain(c) {
				// TODO: Support the Name property for Hangul syllables following NR1.
				// See section 4.8 Name in [Unicode].
				if strings.HasPrefix(prefix, "HANGUL SYLLABLE") {
					return property.NewPropertyName("<Hangul Syllable>")
				}

				return property.NewPropertyName(fmt.Sprintf("%v%X", prefix, c))
			}
		}
	}
	for na, cp := range u.UnicodeData.Name {
		if cp.Contain(c) {
			return na
		}
	}
	return property.NewPropertyName("")
}

func (u *UCD) lookupNameAlias(c rune) property.PropertyNameList {
	for _, e := range u.NameAliases.Entries {
		if e.CP == c {
			as := make([]property.PropertyName, len(e.Aliases))
			for i, alias := range e.Aliases {
				as[i] = alias
			}
			return property.NewPropertyNameList(as)
		}
	}
	return nil
}

func (u *UCD) lookupGeneralCategory(c rune) property.PropertyValueSymbol {
	for gc, cps := range u.UnicodeData.GeneralCategory {
		for _, cp := range cps {
			if cp.Contain(c) {
				return gc
			}
		}
	}
	return u.PropertyValueAliases.DefaultValues[property.PropNameGeneralCategory].Value
}

func (u *UCD) isAlphabetic(c rune) property.PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[property.PropNameAlphabetic] {
		if cp.Contain(c) {
			return property.BinaryYes
		}
	}
	return property.BinaryNo
}

func (u *UCD) isLowercase(c rune) property.PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[property.PropNameLowercase] {
		if cp.Contain(c) {
			return property.BinaryYes
		}
	}
	return property.BinaryNo
}

func (u *UCD) isUppercase(c rune) property.PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[property.PropNameUppercase] {
		if cp.Contain(c) {
			return property.BinaryYes
		}
	}
	return property.BinaryNo
}

func (u *UCD) isIDStart(c rune) property.PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[property.PropNameIDStart] {
		if cp.Contain(c) {
			return property.BinaryYes
		}
	}
	return property.BinaryNo
}

func (u *UCD) isIDContinue(c rune) property.PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[property.PropNameIDContinue] {
		if cp.Contain(c) {
			return property.BinaryYes
		}
	}
	return property.BinaryNo
}

func (u *UCD) isXIDStart(c rune) property.PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[property.PropNameXIDStart] {
		if cp.Contain(c) {
			return property.BinaryYes
		}
	}
	return property.BinaryNo
}

func (u *UCD) isXIDContinue(c rune) property.PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[property.PropNameXIDContinue] {
		if cp.Contain(c) {
			return property.BinaryYes
		}
	}
	return property.BinaryNo
}

func (u *UCD) isWhiteSpace(c rune) property.PropertyValueBinary {
	for _, cp := range u.PropList.WhiteSpace {
		if cp.Contain(c) {
			return property.BinaryYes
		}
	}
	return property.BinaryNo
}
