package ucd

import (
	"fmt"
	"strings"

	"github.com/nihei9/ucdx/ucd/parser"
)

type PropertyName string

const (
	PropNameName            PropertyName = "Name"
	PropNameNameAlias       PropertyName = "Name_Alias"
	PropNameGeneralCategory PropertyName = "General_Category"
	PropNameWhiteSpace      PropertyName = "White_Space"
	PropNameAlphabetic      PropertyName = "Alphabetic"
	PropNameLowercase       PropertyName = "Lowercase"
	PropNameUppercase       PropertyName = "Uppercase"
)

type PropertyValue interface {
	fmt.Stringer
}

type PropertyValueName string

func newNamePropertyValue(v string) PropertyValueName {
	return PropertyValueName(v)
}

func (v PropertyValueName) String() string {
	return string(v)
}

type PropertyValueNameList []PropertyValueName

func newNameListPropertyValue(v []PropertyValueName) PropertyValueNameList {
	return PropertyValueNameList(v)
}

func (v PropertyValueNameList) String() string {
	if len(v) == 0 {
		return ""
	}

	var b strings.Builder
	fmt.Fprint(&b, v[0].String())
	for _, name := range v[1:] {
		fmt.Fprintf(&b, ", %v", name)
	}
	return b.String()
}

type PropertyValueSymbol string

func newSymbolPropertyValue(v string) PropertyValueSymbol {
	return PropertyValueSymbol(v)
}

func (v PropertyValueSymbol) String() string {
	return string(v)
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

type BaseProperties struct {
	Properties map[PropertyName]*Property
}

type PropertySet struct {
	Properties            map[PropertyName]*Property
	GeneralCategoryGroups []PropertyValueSymbol
}

func (s *PropertySet) Lookup(propName PropertyName) *Property {
	return s.Properties[propName]
}

type UCD struct {
	UnicodeData           *parser.UnicodeData
	NameAliases           *parser.NameAliases
	DerivedCoreProperties *parser.DerivedCoreProperties
	PropertyValueAliases  *parser.PropertyValueAliases
	PropList              *parser.PropList
}

func (u *UCD) AnalizeCodePoint(c rune) *PropertySet {
	gc := u.lookupGeneralCategory(c)
	return &PropertySet{
		Properties: map[PropertyName]*Property{
			PropNameName:            newProperty(PropNameName, u.lookupName(c)),
			PropNameNameAlias:       newProperty(PropNameNameAlias, u.lookupNameAlias(c)),
			PropNameGeneralCategory: newProperty(PropNameGeneralCategory, gc),
			PropNameAlphabetic:      newProperty(PropNameAlphabetic, u.isAlphabetic(c)),
			PropNameUppercase:       newProperty(PropNameUppercase, u.isUppercase(c)),
			PropNameLowercase:       newProperty(PropNameLowercase, u.isLowercase(c)),
			PropNameWhiteSpace:      newProperty(PropNameWhiteSpace, u.isWhiteSpace(c)),
		},
		GeneralCategoryGroups: lookupGCGroups(gc),
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

func (u *UCD) lookupNameAlias(c rune) PropertyValueNameList {
	for _, e := range u.NameAliases.Entries {
		if e.CP == c {
			as := make([]PropertyValueName, len(e.Aliases))
			for i, alias := range e.Aliases {
				as[i] = newNamePropertyValue(alias)
			}
			return newNameListPropertyValue(as)
		}
	}
	return nil
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

func (u *UCD) isAlphabetic(c rune) PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[string(PropNameAlphabetic)] {
		if cp.Contain(c) {
			return BinaryYes
		}
	}
	return BinaryNo
}

func (u *UCD) isLowercase(c rune) PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[string(PropNameLowercase)] {
		if cp.Contain(c) {
			return BinaryYes
		}
	}
	return BinaryNo
}

func (u *UCD) isUppercase(c rune) PropertyValueBinary {
	for _, cp := range u.DerivedCoreProperties.Entries[string(PropNameUppercase)] {
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
