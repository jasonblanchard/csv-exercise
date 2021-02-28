package parser

import (
	"encoding/csv"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Name CSV name record
type Name struct {
	First  string `json:"first"`
	Middle string `json:"middle,omitempty"`
	Last   string `json:"last"`
}

// Entity individual record in the CSV
type Entity struct {
	InternalID int    `json:"id"`
	Name       *Name  `json:"name"`
	Phone      string `json:"phone"`
}

// RowError error that occurred when trying to parse a row
type RowError struct {
	LineNumber int
	Message    string
}

// CsvToRecords Convert CSV string into a slice or record slices
func CsvToRecords(s string) ([][]string, error) {
	r := csv.NewReader(strings.NewReader(s))
	return r.ReadAll()
}

// RecordsToEntities converts records from CSV reader to entity instances
// TODO: Output row-level errors here?
func RecordsToEntities(input [][]string) ([]*Entity, []*RowError, error) {
	rows := input[1:]
	fmt.Println(rows)

	output := []*Entity{}
	errors := []*RowError{}

	for i, record := range rows {
		lineNumber := i + 2 // Account for header + non-zero indexed CSV
		id, err := strconv.Atoi(record[0])

		if err != nil {
			errors = append(errors, &RowError{
				LineNumber: lineNumber,
				Message:    "Failed to parse ID",
			})
		}

		name := &Name{
			First:  record[1],
			Middle: record[2],
			Last:   record[3],
		}

		entity := &Entity{
			InternalID: id,
			Name:       name,
			Phone:      record[4],
		}

		output = append(output, entity)
		errors = append(errors, collectValidationErrors(entity, lineNumber)...)
	}

	return output, errors, nil
}

func collectValidationErrors(entity *Entity, lineNumber int) []*RowError {
	errors := []*RowError{}

	if err := ValidateFirstName(entity.Name.First); err != nil {
		errors = append(errors, &RowError{
			LineNumber: lineNumber,
			Message:    err.Error(),
		})
	}

	if err := ValidateLastName(entity.Name.Last); err != nil {
		errors = append(errors, &RowError{
			LineNumber: lineNumber,
			Message:    err.Error(),
		})
	}

	if err := ValidateMiddleName(entity.Name.Middle); err != nil {
		errors = append(errors, &RowError{
			LineNumber: lineNumber,
			Message:    err.Error(),
		})
	}

	if err := ValidateInternalID(entity.InternalID); err != nil {
		errors = append(errors, &RowError{
			LineNumber: lineNumber,
			Message:    err.Error(),
		})
	}

	if err := ValidatePhoneNumber(entity.Phone); err != nil {
		errors = append(errors, &RowError{
			LineNumber: lineNumber,
			Message:    err.Error(),
		})
	}

	return errors
}

// ValidateFirstName ensure first name field is valid
func ValidateFirstName(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("FIRST_NAME is required")
	}

	if len(value) > 15 {
		return fmt.Errorf("FIRST_NAME must be fewer than 15 characters")
	}

	return nil
}

// ValidateLastName ensure last name field is valid
func ValidateLastName(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("LAST_NAME is required")
	}

	if len(value) > 15 {
		return fmt.Errorf("LAST_NAME must be fewer than 15 characters")
	}

	return nil
}

// ValidateMiddleName ensure middle name field is valid
func ValidateMiddleName(value string) error {
	if len(value) > 15 {
		return fmt.Errorf("MIDDLE_NAME must be fewer than 15 characters")
	}

	return nil
}

// ValidateInternalID ensures internal id field is valid
func ValidateInternalID(value int) error {
	digits := len(fmt.Sprintf("%d", value))

	if (value <= 0) || (digits != 8) {
		return fmt.Errorf("INTERNAL_ID must be an 8 digit positive integer")
	}

	return nil
}

// ValidatePhoneNumber ensures phone number is valid
func ValidatePhoneNumber(value string) error {
	match, _ := regexp.MatchString("^\\d{3}-\\d{3}-\\d{4}$", value)
	if match == true {
		return nil
	}
	return fmt.Errorf("PHONE_NUMBER should match pattern ###-###-####")
}
