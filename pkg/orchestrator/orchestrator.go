package orchestrator

import (
	"encoding/json"
	"fmt"
	"io/fs"
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
type Orchestrator struct {
	InputDirectory  string
	OutputDirectory string
	ErrorDirectory  string
	Processed       []string
}

// Run Watch directory and handle new files
func (o *Orchestrator) Run() {
	fmt.Println(fmt.Sprintf("Watching for CSV files in %s", o.InputDirectory))

	// for {
	// 	files, err := ioutil.ReadDir(o.InputDirectory)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	for _, file := range files {
	// 		if match, _ := regexp.Match("^.+\\.csv$", []byte(file.Name())); match == true {
	// 			fmt.Println(fmt.Sprintf("Processing %s/%s", o.InputDirectory, file.Name()))
	// 			err := o.HandleFile(file.Name())
	// 			if err != nil {
	// 				fmt.Println(err)
	// 			}
	// 		}
	// 	}

	// 	// Wait a bit between cycles to let the calling code catch up and prevent unnecessary work.
	// 	time.Sleep(100 * time.Millisecond)
	// }

	// Process files already in the input directory
	files, err := ioutil.ReadDir(o.InputDirectory)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if match, _ := regexp.Match("^.+\\.csv$", []byte(file.Name())); match == true {
			fmt.Println(fmt.Sprintf("Processing %s/%s", o.InputDirectory, file.Name()))
			err := o.HandleFile(file.Name())
			if err != nil {
				fmt.Println(err)
			}
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
				err := o.HandleFile(event.FileInfo.Name())
				if err != nil {
					// TODO: What should happen in this case? Remove the file? Try to re-process it again?
					fmt.Println(err)
				}
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
	// Noop if we've already processed this file
	if contains(o.Processed, filename) {
		return nil
	}

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
	}

	err = os.Remove(inputFile)
	if err != nil {
		return err
	}

	// NOTE: I think this could create a race condition between orchestrator passes because we're changing shared state inside the goroutines.
	// TODO: Clarify this part of the spec: "in the event of file name collision, the latest file should overwrite the earlier version."
	o.Processed = append(o.Processed, filename)

	return nil
}

// Clean remove *.csv files from directories
func (o *Orchestrator) Clean() {
	outputFiles, err := ioutil.ReadDir(o.OutputDirectory)
	if err != nil {
		panic(err)
	}
	cleanDir(o.OutputDirectory, outputFiles)

	errorFiles, err := ioutil.ReadDir(o.ErrorDirectory)
	if err != nil {
		panic(err)
	}
	cleanDir(o.ErrorDirectory, errorFiles)
}

func cleanDir(dir string, files []fs.FileInfo) error {
	for _, file := range files {
		if match, _ := regexp.Match("^.+\\.json|.+\\.csv$", []byte(file.Name())); match == true {
			path := fmt.Sprintf("%s/%s", dir, file.Name())
			fmt.Println(fmt.Sprintf("Removing %s", path))
			err := os.Remove(path)
			if err != nil {
				fmt.Println(err)
			}
		}
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
