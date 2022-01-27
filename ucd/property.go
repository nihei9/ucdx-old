package ucd

import (
	"fmt"

	"github.com/nihei9/ucdx/ucd/parser"
)

type PropertyName string

const (
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

func (u *UCD) lookupGeneralCategory(c rune) PropertyValueSymbol {
	for gc, cps := range u.UnicodeData.GeneralCategory {
		for _, cp := range cps {
			from, to := cp.Range()
			if c >= from && c <= to {
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
