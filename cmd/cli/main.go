package main

import (
	"flag"

	"github.com/jasonblanchard/csv-exercise/pkg/orchestrator"
)

func main() {
	var inputDirectory, outputDirectory, errorDirectory string
	var clean bool
	flag.StringVar(&inputDirectory, "input", ".", "input directory")
	flag.StringVar(&outputDirectory, "output", ".", "output directory")
	flag.StringVar(&errorDirectory, "errors", ".", "error directory")
	flag.BoolVar(&clean, "clean", false, "remove *.csv files from output and error directories before running")

	flag.Parse()

	o := &orchestrator.Orchestrator{
		InputDirectory:  inputDirectory,
		OutputDirectory: outputDirectory,
		ErrorDirectory:  errorDirectory,
	}

	if clean == true {
		o.Clean()
	}

	o.Run()
}
