package main

import (
	"flag"

	"github.com/jasonblanchard/csv-exercise/pkg/orchestrator"
)

func main() {
	var inputDirectory, outputDirectory, errorDirectory string
	flag.StringVar(&inputDirectory, "input", ".", "input directory")
	flag.StringVar(&outputDirectory, "output", ".", "output directory")
	flag.StringVar(&errorDirectory, "errors", ".", "error directory")

	flag.Parse()

	o := &orchestrator.Orchestrator{
		InputDirectory:  inputDirectory,
		OutputDirectory: outputDirectory,
		ErrorDirectory:  errorDirectory,
	}

	// err := o.HandleFile("example.csv")
	// if err != nil {
	// 	panic(err)
	// }
	o.Run()
}
