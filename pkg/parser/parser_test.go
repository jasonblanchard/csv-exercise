package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCsvToRecords(t *testing.T) {
	input := `INTERNAL_ID,FIRST_NAME,MIDDLE_NAME,LAST_NAME,PHONE_NUM
12345678,Bobby,Carmichael,Tables,555-555-5555
`
	output, err := CsvToRecords(input)
	if err != nil {
		panic(err)
	}
	expected := [][]string{
		{"INTERNAL_ID", "FIRST_NAME", "MIDDLE_NAME", "LAST_NAME", "PHONE_NUM"},
		{"12345678", "Bobby", "Carmichael", "Tables", "555-555-5555"},
	}
	assert.Equal(t, expected, output)
}

func TestRecordsToEntities(t *testing.T) {
	tests := []struct {
		input          [][]string
		expected       []*Entity
		expectedErrors []*RowError
	}{
		{
			input: [][]string{
				{"INTERNAL_ID", "FIRST_NAME", "MIDDLE_NAME", "LAST_NAME", "PHONE_NUM"},
				{"12345678", "Bobby", "Carmichael", "Tables", "555-555-5555"},
			},
			expected: []*Entity{
				{
					InternalID: 12345678,
					Name: &Name{
						First:  "Bobby",
						Middle: "Carmichael",
						Last:   "Tables",
					},
					Phone: "555-555-5555",
				},
			},
			expectedErrors: []*RowError{},
		},
		{
			input: [][]string{
				{"INTERNAL_ID", "FIRST_NAME", "MIDDLE_NAME", "LAST_NAME", "PHONE_NUM"},
				{"notavalidnumber", "Bobby", "Carmichael", "Tables", "555-555-5555"},
			},
			expected: []*Entity{
				{
					InternalID: 0,
					Name: &Name{
						First:  "Bobby",
						Middle: "Carmichael",
						Last:   "Tables",
					},
					Phone: "555-555-5555",
				},
			},
			expectedErrors: []*RowError{
				{
					LineNumber: 2,
					Message:    "Failed to parse ID",
				},
				{
					LineNumber: 2,
					Message:    "INTERNAL_ID must be an 8 digit positive integer",
				},
			},
		},
		{
			input: [][]string{
				{"INTERNAL_ID", "FIRST_NAME", "MIDDLE_NAME", "LAST_NAME", "PHONE_NUM"},
				{"12345678", "", "Carmichael", "Tables", "555-555-5555"},
			},
			expected: []*Entity{
				{
					InternalID: 12345678,
					Name: &Name{
						First:  "",
						Middle: "Carmichael",
						Last:   "Tables",
					},
					Phone: "555-555-5555",
				},
			},
			expectedErrors: []*RowError{
				{
					LineNumber: 2,
					Message:    "FIRST_NAME is required",
				},
			},
		},
	}

	for i, test := range tests {
		output, errors, err := RecordsToEntities(test.input)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, test.expected, output, fmt.Sprintf("Test %v", i), test)
		assert.Equal(t, test.expectedErrors, errors, fmt.Sprintf("Test %v", i), test)
	}
}

func TestValidateFirstName(t *testing.T) {
	tests := []struct {
		input    string
		expected error
	}{
		{input: "Bobby", expected: nil},
		{input: "", expected: fmt.Errorf("FIRST_NAME is required")},
		{input: "aaaaaaaaaaaaaaaa", expected: fmt.Errorf("FIRST_NAME must be fewer than 15 characters")},
	}

	for i, test := range tests {
		output := ValidateFirstName(test.input)
		assert.Equal(t, test.expected, output, fmt.Sprintf("Test %v", i), test)
	}
}

func TestValidateLastName(t *testing.T) {
	tests := []struct {
		input    string
		expected error
	}{
		{input: "Tables", expected: nil},
		{input: "", expected: fmt.Errorf("LAST_NAME is required")},
		{input: "aaaaaaaaaaaaaaaa", expected: fmt.Errorf("LAST_NAME must be fewer than 15 characters")},
	}

	for i, test := range tests {
		output := ValidateLastName(test.input)
		assert.Equal(t, test.expected, output, fmt.Sprintf("Test %v", i), test)
	}
}

func TestValidateMiddleName(t *testing.T) {
	tests := []struct {
		input    string
		expected error
	}{
		{input: "Carmichael", expected: nil},
		{input: "", expected: nil},
		{input: "aaaaaaaaaaaaaaaa", expected: fmt.Errorf("MIDDLE_NAME must be fewer than 15 characters")},
	}

	for i, test := range tests {
		output := ValidateMiddleName(test.input)
		assert.Equal(t, test.expected, output, fmt.Sprintf("Test %v", i), test)
	}
}

func TestValidateInternalID(t *testing.T) {
	tests := []struct {
		input    int
		expected error
	}{
		{input: 11111111, expected: nil},
		{input: -1234, expected: fmt.Errorf("INTERNAL_ID must be an 8 digit positive integer")},
		{input: 0, expected: fmt.Errorf("INTERNAL_ID must be an 8 digit positive integer")},
		{input: 123, expected: fmt.Errorf("INTERNAL_ID must be an 8 digit positive integer")},
	}

	for i, test := range tests {
		output := ValidateInternalID(test.input)
		assert.Equal(t, test.expected, output, fmt.Sprintf("Test %v", i), test)
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected error
	}{
		{input: "444-444-4444", expected: nil},
		{input: "", expected: fmt.Errorf("PHONE_NUMBER should match pattern ###-###-####")},
		{input: "4444", expected: fmt.Errorf("PHONE_NUMBER should match pattern ###-###-####")},
		{input: "444-444-44444", expected: fmt.Errorf("PHONE_NUMBER should match pattern ###-###-####")},
	}

	for i, test := range tests {
		output := ValidatePhoneNumber(test.input)
		assert.Equal(t, test.expected, output, fmt.Sprintf("Test %v", i), test)
	}
}
