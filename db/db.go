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
		ucd.TxtPropertyValueAliases,
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
	case ucd.TxtPropertyValueAliases:
		data, err = parser.ParsePropertyValueAliases(f)
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
