package ucd

import "fmt"

const UnicodeVersion = "13.0.0"

const (
	TxtUnicodeData          = "UnicodeData.txt"
	TxtNameAliases          = "NameAliases.txt"
	TxtPropertyValueAliases = "PropertyValueAliases.txt"
	TxtPropList             = "PropList.txt"
)

func MakeDataFileURL(dataFileName string) string {
	return fmt.Sprintf("https://www.unicode.org/Public/%v/ucd/%v", UnicodeVersion, dataFileName)
}
