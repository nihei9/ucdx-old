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
)

type DBConfig struct {
	AppDirPath string
}

func MakeDB(config *DBConfig) error {
	dataFileNames := []string{
		ucd.TxtUnicodeData,
		ucd.TxtNameAliases,
		ucd.TxtDerivedCoreProperties,
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
	var ud *parser.UnicodeData
	{
		d, err := os.ReadFile(makeParsedDataFilePath(appDirPath, ucd.TxtUnicodeData))
		if err != nil {
			return nil, err
		}
		ud = &parser.UnicodeData{}
		err = json.Unmarshal(d, ud)
		if err != nil {
			return nil, err
		}
	}

	var nameAliases *parser.NameAliases
	{
		d, err := os.ReadFile(makeParsedDataFilePath(appDirPath, ucd.TxtNameAliases))
		if err != nil {
			return nil, err
		}
		nameAliases = &parser.NameAliases{}
		err = json.Unmarshal(d, nameAliases)
		if err != nil {
			return nil, err
		}
	}

	var derivedCoreProps *parser.DerivedCoreProperties
	{
		d, err := os.ReadFile(makeParsedDataFilePath(appDirPath, ucd.TxtDerivedCoreProperties))
		if err != nil {
			return nil, err
		}
		derivedCoreProps = &parser.DerivedCoreProperties{}
		err = json.Unmarshal(d, derivedCoreProps)
		if err != nil {
			return nil, err
		}
	}

	var propValAliases *parser.PropertyValueAliases
	{
		d, err := os.ReadFile(makeParsedDataFilePath(appDirPath, ucd.TxtPropertyValueAliases))
		if err != nil {
			return nil, err
		}
		propValAliases = &parser.PropertyValueAliases{}
		err = json.Unmarshal(d, propValAliases)
		if err != nil {
			return nil, err
		}
	}

	var propList *parser.PropList
	{
		d, err := os.ReadFile(makeParsedDataFilePath(appDirPath, ucd.TxtPropList))
		if err != nil {
			return nil, err
		}
		propList = &parser.PropList{}
		err = json.Unmarshal(d, propList)
		if err != nil {
			return nil, err
		}
	}

	return &ucd.UCD{
		UnicodeData:           ud,
		NameAliases:           nameAliases,
		DerivedCoreProperties: derivedCoreProps,
		PropertyValueAliases:  propValAliases,
		PropList:              propList,
	}, nil
}

func makeParsedDataFilePath(appDirPath string, srcDataFileName string) string {
	return filepath.Join(appDirPath, "db", fmt.Sprintf("%v.json", strings.TrimSuffix(srcDataFileName, ".txt")))
}
