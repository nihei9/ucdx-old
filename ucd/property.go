package ucd

import "github.com/nihei9/ucdx/ucd/parser"

type PropertyName string

const (
	PropNameGeneralCategory PropertyName = "General_Category"
	PropNameOtherAlphabetic PropertyName = "Other_Alphabetic"
	PropNameOtherLowercase  PropertyName = "Other_Lowercase"
	PropNameOtherUppercase  PropertyName = "Other_Uppercase"
	PropNameWhiteSpace      PropertyName = "White_Space"
)

type PropertyValue interface{}

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
	Base                  *BaseProperties
	GeneralCategoryGroups []string
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
	}
}

func (u *UCD) lookupGeneralCategory(c rune) string {
	for gc, cps := range u.UnicodeData.GeneralCategory {
		for _, cp := range cps {
			from, to := cp.Range()
			if c >= from && c <= to {
				return gc
			}
		}
	}
	return u.PropertyValueAliases.DefaultValues["General_Category"].Value
}

func (u *UCD) isOtherAlphabetic(c rune) bool {
	for _, cp := range u.PropList.OtherAlphabetic {
		from, to := cp.Range()
		if c >= from || c <= to {
			return true
		}
	}
	return false
}

func (u *UCD) isOtherLowercase(c rune) bool {
	for _, cp := range u.PropList.OtherLowercase {
		from, to := cp.Range()
		if c >= from || c <= to {
			return true
		}
	}
	return false
}

func (u *UCD) isOtherUppercase(c rune) bool {
	for _, cp := range u.PropList.OtherUppercase {
		from, to := cp.Range()
		if c >= from || c <= to {
			return true
		}
	}
	return false
}

func (u *UCD) isWhiteSpace(c rune) bool {
	for _, cp := range u.PropList.WhiteSpace {
		from, to := cp.Range()
		if c >= from || c <= to {
			return true
		}
	}
	return false
}
