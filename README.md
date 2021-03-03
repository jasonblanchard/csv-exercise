## How-To
Build the program:
```
go build -o csv2json cmd/cli/main.go
```

Run it:
```
./csv2json -input example/input -output example/output -errors example/errors
```

Pass in the `-clean` flag to clean out the output and error directories before processing new files.

Add csv file to `example/input` and see the output in `example/output` and errors in `example/errors`

## Assumptions
1. It was unclear how to resolve these two specs: `files will be considered new if the file name has not been recorded as processed before` and `in the event of file name collision, the latest file should overwrite the earlier version`. For now, the application is keeping track of processed files and NOT processing them again until you restart the program, so I'm not sure you'd ever hit the name collision case.
1. Given ^, it wasn't obvious what the expected behavior was for the source file if it has already been processed. For now, I'm leaving it in the input directory and NOT removing it.
1. There's an annoying edge case when you try to process a file that gets a CSV parse error, like if a row has the wrong number of fields, it'll leave the file in the input directory until you remove it, correct the error and try again. I'm assuming this is an acceptable flow.
1. In all cases, I'm reading and writing entire files instead of reading/writing them in buffered chunks. This assumes that the CSV files aren't too large for that to cause a problem for taget systems. What "too large" means probably requires some benchmarking and a better understanding of the target customer.
1. This CLI is not going to get any more complicated than this. If not, I'd reach for something like [Cobra](https://github.com/spf13/cobra).
1. Using a library to do the directory watching is more stable than rolling my own. Initially, I was just continuously checking the source directory in a for loop, but that led to all sorts of problems trying to re-process the same files multiple times, etc.
1. The field validations could probably be refactored to be bit more DRY and reusable. I'm keeping them all separate for now and eating the duplication until we know that's the right abstraction.