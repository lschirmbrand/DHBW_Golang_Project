package main

import (
	"DHBW_Golang_Project/internal/config"
	"regexp"
)

func checkFlagFunctionality() (bool, *[]string) {
	/*
		The function is used at the start of the execution
		to check, whether als passed flags are set and
		formatted correctly.
		In case of the passed date, additionally the function
		checks, whether a log file exists for this date.
	 */
	flagsOk := true
	fails := make([]string, 0)
	ok, _ := validateDateInput(*config.Date)
	if !ok {
		flagsOk = false
		fails = append(fails, "Date input was incorrect.")
	}
	if !checkFileExistence(buildFileLogPath(*config.Date)){
		flagsOk = false
		fails = append(fails, "No log file exists for the day: " + *config.Date)
	}
	ok, _ = validateOperationInput(*config.Operation)
	if !ok {
		flagsOk = false
		fails = append(fails, "Operation input was incorrect.")
	}
	ok = validateQueryInput(*config.Query)
	if !ok {
		flagsOk = false
		fails = append(fails, "Query input was incorrect.")
	}
	// All flag errors will be returned
	return flagsOk, &fails
}

func validateDateInput(date string) (bool, error) {
	// Function uses regex to validate the date input format
	return regexp.Match("^([19|20].(0[1-9]|[1-9][1-9]))[-](0[1-9]|1[012])[-](0[1-9]|[12][0-9]|3[01])$", []byte(date))
}

func validateOperationInput(operation string) (bool, error) {
	// Function uses regex to validate the operation input format
	return regexp.Match("(?i)\\bvisitor\\b|location\\b|contact\\b", []byte(operation))
}

func validateYesNoInput(operation string) (bool, error) {
	// Function uses regex to validate, whether a result should get exported
	return regexp.Match("(?i)\\b[y|n]\\b|yes\\b|no\\b", []byte(operation))
}

func validateQueryInput(s string) bool {
	// Function checks for the length of the query input
	return len(s) > 0
}
