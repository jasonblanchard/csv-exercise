package orchestrator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jasonblanchard/csv-exercise/pkg/parser"
	"github.com/radovskyb/watcher"
)

// Orchestrator orchestrates the application lifecycle
// TODO: Refactor to reader/writer interfaces? That would be more testable
type Orchestrator struct {
	InputDirectory  string
	OutputDirectory string
	ErrorDirectory  string
	Processing      []string
}

// Run Watch directory and handle new files
func (o *Orchestrator) Run() {
	// Process all the files in the directory already
	files, err := ioutil.ReadDir(o.InputDirectory)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if match, _ := regexp.Match("^.+\\.csv$", []byte(file.Name())); match == true {
			o.HandleFile(file.Name())
		}
	}

	// Watch for new files
	w := watcher.New()
	w.FilterOps(watcher.Create)
	r := regexp.MustCompile("^.+\\.csv$")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(fmt.Sprintf("Processing %s", event.FileInfo.Name()))
				o.HandleFile(event.FileInfo.Name())
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.Add(o.InputDirectory); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

// HandleFile parses file and writes results to output and error directory
func (o *Orchestrator) HandleFile(filename string) error {
	// TODO: Only continue if file has not been processed yet

	nameWithoutExtension := strings.Split(filename, ".")[0]
	inputFile := fmt.Sprintf("%s/%s", o.InputDirectory, filename)
	outputFile := fmt.Sprintf("%s/%s.json", o.OutputDirectory, nameWithoutExtension)
	errorFile := fmt.Sprintf("%s/%s.csv", o.ErrorDirectory, nameWithoutExtension)

	data, err := ioutil.ReadFile(inputFile)
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

	err = ioutil.WriteFile(outputFile, json, 0644)
	if err != nil {
		return err
	}

	if len(rowErrors) > 0 {
		err := ioutil.WriteFile(errorFile, []byte(rowErrorsCsv), 0644)
		if err != nil {
			return err
		}
	} else {
		// TODO: Note sure it this is desireable, putting here for convenience
		_, err := os.Stat(errorFile)
		if !os.IsNotExist(err) {
			err = os.Remove(errorFile)
			if err != nil {
				return err
			}
		}
	}

	err = os.Remove(inputFile)
	if err != nil {
		return err
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func remove(s []string, e string) []string {
	output := []string{}
	for _, a := range s {
		if a != e {
			output = append(output, a)
		}
	}
	return output
}
