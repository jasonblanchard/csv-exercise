package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jasonblanchard/csv-exercise/pkg/orchestrator"
)

func main() {
	var inputDirectory, outputDirectory, errorDirectory string
	var clean bool
	flag.StringVar(&inputDirectory, "input", "", "input directory")
	flag.StringVar(&outputDirectory, "output", ".", "output directory")
	flag.StringVar(&errorDirectory, "errors", ".", "error directory")
	flag.BoolVar(&clean, "clean", false, "remove *.csv files from output and error directories before running")

	flag.Parse()

	// Require explicit input directory because it's doing a destructive action (deleting *.csv)
	if inputDirectory == "" {
		fmt.Println("Input directory is required")
		os.Exit(1)
	}

	o := &orchestrator.Orchestrator{
		InputDirectory:     inputDirectory,
		OutputDirectory:    outputDirectory,
		ErrorDirectory:     errorDirectory,
		FinishedProcessing: make(chan string),
	}

	if clean == true {
		o.Clean()
	}

	o.Run()
}
