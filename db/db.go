package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nihei9/ucdx/ucd"
	"github.com/nihei9/ucdx/ucd/parser"
	"github.com/nihei9/ucdx/ucd/property"
)

type DBConfig struct {
	AppDirPath string
}

func MakeDB(config *DBConfig) error {
	dataFileNames := []string{
		ucd.TxtUnicodeData,
		ucd.TxtNameAliases,
		ucd.TxtDerivedCoreProperties,
		ucd.TxtPropertyAliases,
		ucd.TxtPropertyValueAliases,
		ucd.TxtPropList,
	}

	tempDirPath, err := os.MkdirTemp(config.AppDirPath, "db-*")
	if err != nil {
		return err
	}

	for _, dataFileName := range dataFileNames {
		err := fetchDataFile(dataFileName, tempDirPath)
		if err != nil {
			return err
		}
	}

	for _, dataFileName := range dataFileNames {
		err := parseDataFile(tempDirPath, dataFileName)
		if err != nil {
			return err
		}
	}

	{
		var propAliases *property.PropertyAliases
		{
			b, err := ioutil.ReadFile(filepath.Join(tempDirPath, makeParsedDataFileName(ucd.TxtPropertyAliases)))
			if err != nil {
				return err
			}
			propAliases = &property.PropertyAliases{}
			err = json.Unmarshal(b, propAliases)
			if err != nil {
				return err
			}
		}

		var propValAliases *property.PropertyValueAliases
		{
			b, err := ioutil.ReadFile(filepath.Join(tempDirPath, makeParsedDataFileName(ucd.TxtPropertyValueAliases)))
			if err != nil {
				return err
			}
			propValAliases = &property.PropertyValueAliases{}
			err = json.Unmarshal(b, propValAliases)
			if err != nil {
				return err
			}
		}

		uni := property.NewUnification(propAliases, propValAliases)
		b, err := json.Marshal(uni)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(tempDirPath, "unification.json"), b, 0644)
		if err != nil {
			return err
		}
	}

	dbDirPath := filepath.Join(config.AppDirPath, "db")
	err = os.RemoveAll(dbDirPath)
	if err != nil {
		return err
	}
	err = os.Rename(tempDirPath, dbDirPath)
	if err != nil {
		return err
	}

	return nil
}

func fetchDataFile(dataFileName string, dirPath string) error {
	res, err := http.Get(ucd.MakeDataFileURL(dataFileName))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dirPath, dataFileName), d, 0644)
}

func parseDataFile(dirPath string, dataFileName string) error {
	f, err := os.Open(filepath.Join(dirPath, dataFileName))
	if err != nil {
		return err
	}
	defer f.Close()

	var data interface{}
	switch dataFileName {
	case ucd.TxtUnicodeData:
		data, err = parser.ParseUnicodeData(f)
	case ucd.TxtNameAliases:
		data, err = parser.ParseNameAliases(f)
	case ucd.TxtDerivedCoreProperties:
		data, err = parser.ParseDerivedCoreProperties(f)
	case ucd.TxtPropertyAliases:
		data, err = parser.ParsePropertyAliases(f)
	case ucd.TxtPropertyValueAliases:
		data, err = parser.ParsePropertyValueAliases(f)
	case ucd.TxtPropList:
		data, err = parser.ParsePropList(f)
	default:
		return fmt.Errorf("unknown data file name: %v", dataFileName)
	}
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	jsonFileName := fmt.Sprintf("%v.json", strings.TrimSuffix(dataFileName, ".txt"))
	return ioutil.WriteFile(filepath.Join(dirPath, jsonFileName), jsonData, 0644)
}

func OpenDB(appDirPath string) (*ucd.UCD, error) {
	var ud *property.UnicodeData
	{
		d, err := os.ReadFile(makeParsedDataFilePath(appDirPath, ucd.TxtUnicodeData))
		if err != nil {
			return nil, err
		}
		ud = &property.UnicodeData{}
		err = json.Unmarshal(d, ud)
		if err != nil {
			return nil, err
		}
	}

	var nameAliases *property.NameAliases
	{
		d, err := os.ReadFile(makeParsedDataFilePath(appDirPath, ucd.TxtNameAliases))
		if err != nil {
			return nil, err
		}
		nameAliases = &property.NameAliases{}
		err = json.Unmarshal(d, nameAliases)
		if err != nil {
			return nil, err
		}
	}

	var derivedCoreProps *property.DerivedCoreProperties
	{
		d, err := os.ReadFile(makeParsedDataFilePath(appDirPath, ucd.TxtDerivedCoreProperties))
		if err != nil {
			return nil, err
		}
		derivedCoreProps = &property.DerivedCoreProperties{}
		err = json.Unmarshal(d, derivedCoreProps)
		if err != nil {
			return nil, err
		}
	}

	var propAliases *property.PropertyAliases
	{
		d, err := os.ReadFile(makeParsedDataFilePath(appDirPath, ucd.TxtPropertyAliases))
		if err != nil {
			return nil, err
		}
		propAliases = &property.PropertyAliases{}
		err = json.Unmarshal(d, propAliases)
		if err != nil {
			return nil, err
		}
	}

	var propValAliases *property.PropertyValueAliases
	{
		d, err := os.ReadFile(makeParsedDataFilePath(appDirPath, ucd.TxtPropertyValueAliases))
		if err != nil {
			return nil, err
		}
		propValAliases = &property.PropertyValueAliases{}
		err = json.Unmarshal(d, propValAliases)
		if err != nil {
			return nil, err
		}
	}

	var propList *property.PropList
	{
		d, err := os.ReadFile(makeParsedDataFilePath(appDirPath, ucd.TxtPropList))
		if err != nil {
			return nil, err
		}
		propList = &property.PropList{}
		err = json.Unmarshal(d, propList)
		if err != nil {
			return nil, err
		}
	}

	var unification *property.Unification
	{
		d, err := ioutil.ReadFile(filepath.Join(appDirPath, "db", "unification.json"))
		if err != nil {
			return nil, err
		}
		unification = &property.Unification{}
		err = json.Unmarshal(d, unification)
		if err != nil {
			return nil, err
		}
	}

	return &ucd.UCD{
		UnicodeData:           ud,
		NameAliases:           nameAliases,
		DerivedCoreProperties: derivedCoreProps,
		PropertyAliases:       propAliases,
		PropertyValueAliases:  propValAliases,
		PropList:              propList,
		Unification:           unification,
	}, nil
}

func makeParsedDataFilePath(appDirPath string, srcDataFileName string) string {
	return filepath.Join(appDirPath, "db", makeParsedDataFileName(srcDataFileName))
}

func makeParsedDataFileName(srcDataFileName string) string {
	return fmt.Sprintf("%v.json", strings.TrimSuffix(srcDataFileName, ".txt"))
}
