package property

import (
	"fmt"
	"strings"
)

type CodePointRange [2]rune

func NewCodePointRange(from, to rune) *CodePointRange {
	cp := CodePointRange{}
	cp[0] = from
	cp[1] = to
	return &cp
}

func (r *CodePointRange) String() string {
	from, to := r.Range()
	return fmt.Sprintf("%X..%X", from, to)
}

func (r *CodePointRange) Range() (rune, rune) {
	return r[0], r[1]
}

func (r *CodePointRange) Contain(c rune) bool {
	from, to := r.Range()
	return c >= from && c <= to
}

type PropertyValue interface {
	fmt.Stringer
}

type PropertyName string

func NewPropertyName(s string) PropertyName {
	return PropertyName(s)
}

func (n PropertyName) String() string {
	return string(n)
}

const (
	PropNameName            PropertyName = "Name"
	PropNameNameAlias       PropertyName = "Name_Alias"
	PropNameGeneralCategory PropertyName = "General_Category"
	PropNameWhiteSpace      PropertyName = "White_Space"
	PropNameAlphabetic      PropertyName = "Alphabetic"
	PropNameLowercase       PropertyName = "Lowercase"
	PropNameUppercase       PropertyName = "Uppercase"
	PropNameIDStart         PropertyName = "ID_Start"
	PropNameIDContinue      PropertyName = "ID_Continue"
	PropNameXIDStart        PropertyName = "ID_XStart"
	PropNameXIDContinue     PropertyName = "ID_XContinue"
)

type PropertyNameList []PropertyName

func NewPropertyNameList(v []PropertyName) PropertyNameList {
	return PropertyNameList(v)
}

func (v PropertyNameList) String() string {
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

func NewSymbolPropertyValue(v string) PropertyValueSymbol {
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

func NewProperty(name PropertyName, value PropertyValue) *Property {
	return &Property{
		Name:  name,
		Value: value,
	}
}

type UnicodeData struct {
	Name            map[PropertyName]*CodePointRange          `json:"name"`
	GeneralCategory map[PropertyValueSymbol][]*CodePointRange `json:"general_category"`
}

func (u *UnicodeData) AddGC(gc PropertyValueSymbol, cp *CodePointRange) {
	// Section 4.2.11 Empty Fields in [UAX44]:
	// > The data file UnicodeData.txt defines many property values in each record. When a field in a data line
	// > for a code point is empty, that indicates that the property takes the default value for that code point.
	if gc == "" {
		return
	}

	cps, ok := u.GeneralCategory[gc]
	if ok {
		from1, to1 := cp.Range()
		i := len(cps) - 1
		from2, to2 := cps[i].Range()
		if from1-to2 == 1 {
			cps[i] = NewCodePointRange(from2, to1)
		} else {
			u.GeneralCategory[gc] = append(cps, cp)
		}
	} else {
		u.GeneralCategory[gc] = []*CodePointRange{cp}
	}
}

type NameAliasesEntry struct {
	CP      rune           `json:"cp"`
	Aliases []PropertyName `json:"aliases"`
}

type NameAliases struct {
	Entries []*NameAliasesEntry `json:"entries"`
}

type DerivedCoreProperties struct {
	Entries map[PropertyName][]*CodePointRange `json:"entries"`
}

type PropertyAlias struct {
	Abb    PropertyName   `json:"abb"`
	Long   PropertyName   `json:"long"`
	Others []PropertyName `json:"others"`
}

type PropertyAliases struct {
	Aliases []*PropertyAlias `json:"aliases"`
}

// PropertyValueAliase represents a set of aliases for a property value.
// `Abb` and `Long` are the preferred aliases.
type PropertyValueAliase struct {
	// Abb is an abbreviated symbolic name for a property value.
	Abb PropertyValueSymbol `json:"abb"`

	// Long is the long symbolic name for a property value.
	Long PropertyValueSymbol `json:"long"`

	// Others is a set of other aliases for a property value.
	Others []PropertyValueSymbol `json:"others,omitempty"`
}

type DefaultValue struct {
	Value PropertyValueSymbol `json:"value"`
	CP    *CodePointRange     `json:"cp"`
}

type PropertyValueAliases struct {
	Aliases       map[PropertyName][]*PropertyValueAliase `json:"aliases"`
	DefaultValues map[PropertyName]*DefaultValue          `json:"default_values"`
}

type PropList struct {
	WhiteSpace []*CodePointRange `json:"White_Space"`
}

type Unification struct {
	PropertyNames  map[PropertyName]PropertyName      `json:"property_names"`
	PropertyValues map[PropertyName]map[string]string `json:"property_values"`
}

func NewUnification(propAliases *PropertyAliases, propValAliases *PropertyValueAliases) *Unification {
	names := map[PropertyName]PropertyName{}
	for _, a := range propAliases.Aliases {
		names[a.Abb] = a.Long
		names[a.Long] = a.Long
		for _, o := range a.Others {
			names[o] = a.Long
		}
	}

	values := map[PropertyName]map[string]string{}
	for propName, aliases := range propValAliases.Aliases {
		if values[propName] == nil {
			values[propName] = map[string]string{}
		}
		for _, a := range aliases {
			values[propName][a.Abb.String()] = a.Long.String()
			values[propName][a.Long.String()] = a.Long.String()
			for _, o := range a.Others {
				values[propName][o.String()] = a.Long.String()
			}
		}
	}

	return &Unification{
		PropertyNames:  names,
		PropertyValues: values,
	}
}
