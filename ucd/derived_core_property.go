package ucd

// See section 5.3 Property Definitions in [UAX44] for the definition of DerivedCoreProperties.
var derivedCoreProperties = map[PropertyName][]*Property{
	PropNameAlphabetic: {
		newProperty(PropNameLowercase, BinaryYes),
		newProperty(PropNameUppercase, BinaryYes),
		newProperty(PropNameGeneralCategory, newSymbolPropertyValue("lt")),
		newProperty(PropNameGeneralCategory, newSymbolPropertyValue("lm")),
		newProperty(PropNameGeneralCategory, newSymbolPropertyValue("lo")),
		newProperty(PropNameGeneralCategory, newSymbolPropertyValue("nl")),
		newProperty(PropNameOtherAlphabetic, BinaryYes),
	},
	PropNameLowercase: {
		newProperty(PropNameGeneralCategory, newSymbolPropertyValue("ll")),
		newProperty(PropNameOtherLowercase, BinaryYes),
	},
	PropNameUppercase: {
		newProperty(PropNameGeneralCategory, newSymbolPropertyValue("lu")),
		newProperty(PropNameOtherUppercase, BinaryYes),
	},
}

type DerivedCoreProperties struct {
	Properies map[PropertyName]*Property
}

func calcDerivedCoreProperties(base *BaseProperties) *DerivedCoreProperties {
	return &DerivedCoreProperties{
		Properies: map[PropertyName]*Property{
			PropNameAlphabetic: newProperty(PropNameAlphabetic, isDerivable(base, PropNameAlphabetic)),
			PropNameLowercase:  newProperty(PropNameLowercase, isDerivable(base, PropNameLowercase)),
			PropNameUppercase:  newProperty(PropNameUppercase, isDerivable(base, PropNameUppercase)),
		},
	}
}

func isDerivable(base *BaseProperties, derivedPropName PropertyName) PropertyValueBinary {
	elems := derivedCoreProperties[derivedPropName]
	for _, elem := range elems {
		baseProp, ok := base.Properties[elem.Name]
		if ok {
			if baseProp.equal(elem) {
				return BinaryYes
			}
			continue
		}
		if isDerivable(base, elem.Name) {
			return BinaryYes
		}
	}
	return BinaryNo
}
