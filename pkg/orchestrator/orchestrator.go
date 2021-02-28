package orchestrator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/jasonblanchard/csv-exercise/pkg/parser"
)

// Read file from input-directory/[file].csv
// Run parser.RecordsToEntities
// Write entities to output-directory [file]].json
// Write errors to error-directory [file].csv
// Remove input-directory/[file].csv

// HandleFile parses file and writes results to output and error directory
func HandleFile(inputDirectory string, filename string, outputDirectory string, errorDirectory string) error {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", inputDirectory, filename))

	if err != nil {
		return err
	}

	records, err := parser.CsvToRecords(string(data))
	if err != nil {
		return err
	}

	entities, rowErrors, err := parser.RecordsToEntities(records)
	if err != nil {
		return err
	}

	json, err := json.Marshal(entities)
	if err != nil {
		return err
	}

	rowErrorsCsv, err := parser.ErrorsToCSV(rowErrors)
	if err != nil {
		return err
	}

	fmt.Println(string(json))
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", outputDirectory, "TODO.json"), json, 0644)
	if err != nil {
		return err
	}

	if len(rowErrors) > 0 {
		err := ioutil.WriteFile(fmt.Sprintf("%s/%s", errorDirectory, "TODO.csv"), []byte(rowErrorsCsv), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
