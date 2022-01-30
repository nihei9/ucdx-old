package ucd

import (
	"fmt"
	"strings"

	"github.com/nihei9/ucdx/ucd/property"
)

type PropertySet struct {
	Properties            map[property.PropertyName]*property.Property
	GeneralCategoryGroups []property.PropertyValueSymbol
}

func (s *PropertySet) Lookup(propName property.PropertyName) *property.Property {
	return s.Properties[propName]
}

type UCD struct {
	UnicodeData           *property.UnicodeData
	NameAliases           *property.NameAliases
	DerivedCoreProperties *property.DerivedCoreProperties
	PropertyValueAliases  *property.PropertyValueAliases
	PropList              *property.PropList
}

func (u *UCD) AnalizeCodePoint(c rune) *PropertySet {
	gc := u.lookupGeneralCategory(c)
	return &PropertySet{
		Properties: map[property.PropertyName]*property.Property{
			property.PropNameName:            property.NewProperty(property.PropNameName, u.lookupName(c)),
			property.PropNameNameAlias:       property.NewProperty(property.PropNameNameAlias, u.lookupNameAlias(c)),
			property.PropNameGeneralCategory: property.NewProperty(property.PropNameGeneralCategory, gc),
			property.PropNameAlphabetic:      property.NewProperty(property.PropNameAlphabetic, u.isAlphabetic(c)),
			property.PropNameUppercase:       property.NewProperty(property.PropNameUppercase, u.isUppercase(c)),
			property.PropNameLowercase:       property.NewProperty(property.PropNameLowercase, u.isLowercase(c)),
			property.PropNameIDStart:         property.NewProperty(property.PropNameIDStart, u.isIDStart(c)),
			property.PropNameIDContinue:      property.NewProperty(property.PropNameIDContinue, u.isIDContinue(c)),
			property.PropNameXIDStart:        property.NewProperty(property.PropNameXIDStart, u.isXIDStart(c)),
			property.PropNameXIDContinue:     property.NewProperty(property.PropNameXIDContinue, u.isXIDContinue(c)),
			property.PropNameWhiteSpace:      property.NewProperty(property.PropNameWhiteSpace, u.isWhiteSpace(c)),
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

func (u *UCD) lookupName(c rune) property.PropertyValueName {
	for prefix, cps := range namePrefixes {
		for _, cp := range cps {
			if cp.Contain(c) {
				// TODO: Support the Name property for Hangul syllables following NR1.
				// See section 4.8 Name in [Unicode].
				if strings.HasPrefix(prefix, "HANGUL SYLLABLE") {
					return property.NewNamePropertyValue("<Hangul Syllable>")
				}

				return property.NewNamePropertyValue(fmt.Sprintf("%v%X", prefix, c))
			}
		}
	}
	for na, cp := range u.UnicodeData.Name {
		if cp.Contain(c) {
			return property.NewNamePropertyValue(na)
		}
	}
	return property.NewNamePropertyValue("")
}

func (u *UCD) lookupNameAlias(c rune) property.PropertyValueNameList {
	for _, e := range u.NameAliases.Entries {
		if e.CP == c {
			as := make([]property.PropertyValueName, len(e.Aliases))
			for i, alias := range e.Aliases {
				as[i] = property.NewNamePropertyValue(alias)
			}
			return property.NewNameListPropertyValue(as)
		}
	}
	return nil
}

func (u *UCD) lookupGeneralCategory(c rune) property.PropertyValueSymbol {
	for gc, cps := range u.UnicodeData.GeneralCategory {
		for _, cp := range cps {
			if cp.Contain(c) {
				return property.NewSymbolPropertyValue(gc)
			}
		}
	}
	return property.NewSymbolPropertyValue(u.PropertyValueAliases.DefaultValues["General_Category"].Value)
}

func (u *UCD) isAlphabetic(c rune) property.PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[string(property.PropNameAlphabetic)] {
		if cp.Contain(c) {
			return property.BinaryYes
		}
	}
	return property.BinaryNo
}

func (u *UCD) isLowercase(c rune) property.PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[string(property.PropNameLowercase)] {
		if cp.Contain(c) {
			return property.BinaryYes
		}
	}
	return property.BinaryNo
}

func (u *UCD) isUppercase(c rune) property.PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[string(property.PropNameUppercase)] {
		if cp.Contain(c) {
			return property.BinaryYes
		}
	}
	return property.BinaryNo
}

func (u *UCD) isIDStart(c rune) property.PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[string(property.PropNameIDStart)] {
		if cp.Contain(c) {
			return property.BinaryYes
		}
	}
	return property.BinaryNo
}

func (u *UCD) isIDContinue(c rune) property.PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[string(property.PropNameIDContinue)] {
		if cp.Contain(c) {
			return property.BinaryYes
		}
	}
	return property.BinaryNo
}

func (u *UCD) isXIDStart(c rune) property.PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[string(property.PropNameXIDStart)] {
		if cp.Contain(c) {
			return property.BinaryYes
		}
	}
	return property.BinaryNo
}

func (u *UCD) isXIDContinue(c rune) property.PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[string(property.PropNameXIDContinue)] {
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
