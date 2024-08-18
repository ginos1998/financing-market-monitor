package readers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func openCsvFile(fileNme string) ([][]string, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, errors.New("error getting working directory: " + err.Error())
	}
	filePath := filepath.Join(rootDir, fileNme)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("error opening file: " + err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("error closing file: ", err)
		}
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("error reading file: " + err.Error())
	}

	return records, nil
}

func checkCsvHeaders(headersFromCsv []string, requiredHeaders []string) bool {
	if len(headersFromCsv) != len(requiredHeaders) {
		return false
	}
	for i, header := range headersFromCsv {
		if header != requiredHeaders[i] {
			return false
		}
	}
	return true
}
