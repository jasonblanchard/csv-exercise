package orchestrator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

	nameWithoutExtension := strings.Split(filename, ".")[0]

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

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", outputDirectory, nameWithoutExtension), json, 0644)
	if err != nil {
		return err
	}

	if len(rowErrors) > 0 {
		err := ioutil.WriteFile(fmt.Sprintf("%s/%s.csv", errorDirectory, nameWithoutExtension), []byte(rowErrorsCsv), 0644)
		if err != nil {
			return err
		}
	} else {
		// TODO: Note sure it this is desireable, putting here for convenience
		err = os.Remove(fmt.Sprintf("%s/%s.csv", errorDirectory, nameWithoutExtension))
		if err != nil {
			return err
		}
	}

	err = os.Remove(fmt.Sprintf("%s/%s", inputDirectory, filename))
	if err != nil {
		return err
	}

	return nil
}
