package ucd

import (
	"fmt"
	"strings"

	"github.com/nihei9/ucdx/ucd/parser"
)

type PropertyName string

const (
	PropNameName            PropertyName = "Name"
	PropNameGeneralCategory PropertyName = "General_Category"
	PropNameOtherAlphabetic PropertyName = "Other_Alphabetic"
	PropNameOtherLowercase  PropertyName = "Other_Lowercase"
	PropNameOtherUppercase  PropertyName = "Other_Uppercase"
	PropNameWhiteSpace      PropertyName = "White_Space"
	PropNameAlphabetic      PropertyName = "Alphabetic"
	PropNameLowercase       PropertyName = "Lowercase"
	PropNameUppercase       PropertyName = "Uppercase"
)

type PropertyValue interface {
	fmt.Stringer
	equal(o PropertyValue) bool
}

type PropertyValueName string

func newNamePropertyValue(v string) PropertyValueName {
	return PropertyValueName(v)
}

func (v PropertyValueName) String() string {
	return string(v)
}

func (v PropertyValueName) equal(o PropertyValue) bool {
	s, ok := o.(PropertyValueName)
	if !ok {
		return false
	}
	return v == s
}

type PropertyValueSymbol string

func newSymbolPropertyValue(v string) PropertyValueSymbol {
	return PropertyValueSymbol(v)
}

func (v PropertyValueSymbol) String() string {
	return string(v)
}

func (v PropertyValueSymbol) equal(o PropertyValue) bool {
	s, ok := o.(PropertyValueSymbol)
	if !ok {
		return false
	}
	return v == s
}

type PropertyValueBinary bool

const (
	BinaryYes PropertyValueBinary = true
	BinaryNo  PropertyValueBinary = false
)

func (v PropertyValueBinary) String() string {
	if v {
		return "Yes"
	}
	return "No"
}

func (v PropertyValueBinary) equal(o PropertyValue) bool {
	b, ok := o.(PropertyValueBinary)
	if !ok {
		return false
	}
	return v == b
}

type Property struct {
	Name  PropertyName
	Value PropertyValue
}

func newProperty(name PropertyName, value PropertyValue) *Property {
	return &Property{
		Name:  name,
		Value: value,
	}
}

func (p *Property) equal(o *Property) bool {
	return p.Name == o.Name && p.Value.equal(o.Value)
}

type BaseProperties struct {
	Properties map[PropertyName]*Property
}

type PropertySet struct {
	Base                  *BaseProperties
	GeneralCategoryGroups []PropertyValueSymbol
	DerivedCore           *DerivedCoreProperties
}

func (s *PropertySet) Lookup(propName PropertyName) *Property {
	if v, ok := s.Base.Properties[propName]; ok {
		return v
	}
	return s.DerivedCore.Properies[propName]
}

type UCD struct {
	UnicodeData          *parser.UnicodeData
	PropertyValueAliases *parser.PropertyValueAliases
	PropList             *parser.PropList
}

func (u *UCD) AnalizeCodePoint(c rune) *PropertySet {
	gc := u.lookupGeneralCategory(c)
	base := &BaseProperties{
		Properties: map[PropertyName]*Property{
			PropNameName:            newProperty(PropNameName, u.lookupName(c)),
			PropNameGeneralCategory: newProperty(PropNameGeneralCategory, gc),
			PropNameOtherAlphabetic: newProperty(PropNameOtherAlphabetic, u.isOtherAlphabetic(c)),
			PropNameOtherLowercase:  newProperty(PropNameOtherLowercase, u.isOtherLowercase(c)),
			PropNameOtherUppercase:  newProperty(PropNameOtherUppercase, u.isOtherUppercase(c)),
			PropNameWhiteSpace:      newProperty(PropNameWhiteSpace, u.isWhiteSpace(c)),
		},
	}
	return &PropertySet{
		Base:                  base,
		GeneralCategoryGroups: lookupGCGroups(gc),
		DerivedCore:           calcDerivedCoreProperties(base),
	}
}

// 4.8 Table 4-8. Name Derivation Rule Prefix Strings in [Unicode].
var namePrefixes = map[string][]*parser.CodePointRange{
	"HANGUL SYLLABLE ": {
		parser.NewCodePointRange(0xAC00, 0xD7A3),
	},
	"CJK UNIFIED IDEOGRAPH-": {
		parser.NewCodePointRange(0x3400, 0x4DBF),
		parser.NewCodePointRange(0x4E00, 0x9FFC),
		parser.NewCodePointRange(0x20000, 0x2A6DD),
		parser.NewCodePointRange(0x2A700, 0x2B734),
		parser.NewCodePointRange(0x2B740, 0x2B81D),
		parser.NewCodePointRange(0x2B820, 0x2CEA1),
		parser.NewCodePointRange(0x2CEB0, 0x2EBE0),
		parser.NewCodePointRange(0x30000, 0x3134A),
	},
	"TANGUT IDEOGRAPH-": {
		parser.NewCodePointRange(0x17000, 0x187F7),
		parser.NewCodePointRange(0x18D00, 0x18D08),
	},
	"KHITAN SMALL SCRIPT CHARACTER-": {
		parser.NewCodePointRange(0x18B00, 0x18CD5),
	},
	"NUSHU CHARACTER-": {
		parser.NewCodePointRange(0x1B170, 0x1B2FB),
	},
	"CJK COMPATIBILITY IDEOGRAPH-": {
		parser.NewCodePointRange(0xF900, 0xFA6D),
		parser.NewCodePointRange(0xFA70, 0xFAD9),
		parser.NewCodePointRange(0x2F800, 0x2FA1D),
	},
}

func (u *UCD) lookupName(c rune) PropertyValueName {
	for prefix, cps := range namePrefixes {
		for _, cp := range cps {
			if cp.Contain(c) {
				// TODO: Support the Name property for Hangul syllables following NR1.
				// See section 4.8 Name in [Unicode].
				if strings.HasPrefix(prefix, "HANGUL SYLLABLE") {
					return newNamePropertyValue("<Hangul Syllable>")
				}

				return newNamePropertyValue(fmt.Sprintf("%v%X", prefix, c))
			}
		}
	}
	for na, cp := range u.UnicodeData.Name {
		if cp.Contain(c) {
			return newNamePropertyValue(na)
		}
	}
	return newNamePropertyValue("")
}

func (u *UCD) lookupGeneralCategory(c rune) PropertyValueSymbol {
	for gc, cps := range u.UnicodeData.GeneralCategory {
		for _, cp := range cps {
			if cp.Contain(c) {
				return newSymbolPropertyValue(gc)
			}
		}
	}
	return newSymbolPropertyValue(u.PropertyValueAliases.DefaultValues["General_Category"].Value)
}

func (u *UCD) isOtherAlphabetic(c rune) PropertyValueBinary {
	for _, cp := range u.PropList.OtherAlphabetic {
		if cp.Contain(c) {
			return BinaryYes
		}
	}
	return BinaryNo
}

func (u *UCD) isOtherLowercase(c rune) PropertyValueBinary {
	for _, cp := range u.PropList.OtherLowercase {
		if cp.Contain(c) {
			return BinaryYes
		}
	}
	return BinaryNo
}

func (u *UCD) isOtherUppercase(c rune) PropertyValueBinary {
	for _, cp := range u.PropList.OtherUppercase {
		if cp.Contain(c) {
			return BinaryYes
		}
	}
	return BinaryNo
}

func (u *UCD) isWhiteSpace(c rune) PropertyValueBinary {
	for _, cp := range u.PropList.WhiteSpace {
		if cp.Contain(c) {
			return BinaryYes
		}
	}
	return BinaryNo
}
