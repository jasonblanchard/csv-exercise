package orchestrator

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleFile(t *testing.T) {
	inputDir, err := ioutil.TempDir("", "csv-exercise-input")
	outputDir, err := ioutil.TempDir("", "csv-exercise-output")
	errorDir, err := ioutil.TempDir("", "csv-exercise-error")
	defer os.RemoveAll(inputDir)
	defer os.RemoveAll(outputDir)
	defer os.RemoveAll(errorDir)

	if err != nil {
		panic(err)
	}

	o := &Orchestrator{
		InputDirectory:  inputDir,
		OutputDirectory: outputDir,
		ErrorDirectory:  errorDir,
	}

	csvInput := `INTERNAL_ID,FIRST_NAME,MIDDLE_NAME,LAST_NAME,PHONE_NUM
12345678,Bobby,,Tables,555-555-5555
55555555,Schmobby,Schmooble,Schmables,222-222-2222
6666666,,,Schneebly,222-2dee22-2222
`
	inputFilePath := fmt.Sprintf("%s/%s", inputDir, "input.csv")
	err = ioutil.WriteFile(inputFilePath, []byte(csvInput), 0644)
	if err != nil {
		panic(err)
	}

	err = o.HandleFile("input.csv")
	if err != nil {
		panic(err)
	}

	ouputFilePath := fmt.Sprintf("%s/%s", outputDir, "input.json")
	output, err := ioutil.ReadFile(ouputFilePath)
	if err != nil {
		panic(err)
	}

	errorFilePath := fmt.Sprintf("%s/%s", errorDir, "input.csv")
	errors, err := ioutil.ReadFile(errorFilePath)
	if err != nil {
		panic(err)
	}
	_, err = os.Stat(inputFilePath)

	assert.Equal(t, "[{\"id\":12345678,\"name\":{\"first\":\"Bobby\",\"last\":\"Tables\"},\"phone\":\"555-555-5555\"},{\"id\":55555555,\"name\":{\"first\":\"Schmobby\",\"middle\":\"Schmooble\",\"last\":\"Schmables\"},\"phone\":\"222-222-2222\"},{\"id\":6666666,\"name\":{\"first\":\"\",\"last\":\"Schneebly\"},\"phone\":\"222-2dee22-2222\"}]", string(output))
	assert.Equal(t, "LINE_NUM,ERROR_MSG\n4,FIRST_NAME is required\n4,INTERNAL_ID must be an 8 digit positive integer\n4,PHONE_NUMBER should match pattern ###-###-####\n", string(errors))
	assert.Equal(t, true, os.IsNotExist(err))

	// Simulate receiving the processed message
	o.Processed = append(o.Processed, "input.csv")

	// Ensure that it doesn't re-process a file with the same name
	csvInput = `INTERNAL_ID,FIRST_NAME,MIDDLE_NAME,LAST_NAME,PHONE_NUM
12345678,Bobby,,Tables II,555-555-5555
`
	err = ioutil.WriteFile(inputFilePath, []byte(csvInput), 0644)
	if err != nil {
		panic(err)
	}

	err = o.HandleFile("input.csv")
	if err != nil {
		panic(err)
	}

	ouputFilePath = fmt.Sprintf("%s/%s", outputDir, "input.json")
	output, err = ioutil.ReadFile(ouputFilePath)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "[{\"id\":12345678,\"name\":{\"first\":\"Bobby\",\"last\":\"Tables\"},\"phone\":\"555-555-5555\"},{\"id\":55555555,\"name\":{\"first\":\"Schmobby\",\"middle\":\"Schmooble\",\"last\":\"Schmables\"},\"phone\":\"222-222-2222\"},{\"id\":6666666,\"name\":{\"first\":\"\",\"last\":\"Schneebly\"},\"phone\":\"222-2dee22-2222\"}]", string(output))
}
