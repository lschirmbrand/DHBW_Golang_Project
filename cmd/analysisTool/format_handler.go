package main

import (
	"regexp"
	"strings"
)

func checkFlagFunctionality(datePtr *string, operationPtr *string, selectedOperation *Operation, queryPtr *string) (bool, *[]string) {
	flagsOk := true
	fails := make([]string, 0)
	ok, _ := validateDateInput(*datePtr)
	if !ok {
		flagsOk = false
		fails = append(fails, "Date input was incorrect.")
	}
	ok, _ = validateOperationInput(*operationPtr)
	if !ok {
		flagsOk = false
		fails = append(fails, "Operation input was incorrect.")
	} else {
		switch strings.ToLower(*operationPtr) {
		case strings.ToLower(string(VISITOR)):
			*selectedOperation = VISITOR
		case strings.ToLower(string(CONTACT)):
			*selectedOperation = CONTACT
		default:
			*selectedOperation = LOCATION
		}
	}
	ok = validateQueryInput(*queryPtr)
	if !ok {
		flagsOk = false
		fails = append(fails, "Query input was incorrect.")
	}
	return flagsOk, &fails
}

func validateDateInput(date string) (bool, error) {
	return regexp.Match("^([19|20].(0[1-9]|[1-9][1-9]))[-](0[1-9]|1[012])[-](0[1-9]|[12][0-9]|3[01])$", []byte(date))
}

func validateOperationInput(operation string) (bool, error) {
	return regexp.Match("(?i)\\bvisitor\\b|location\\b|contact\\b", []byte(operation))
}

func validateYesNoInput(operation string) (bool, error) {
	return regexp.Match("(?i)\\b[y|n]\\b|yes\\b|no\\b", []byte(operation))
}

func validateQueryInput(s string) bool {
	return len(s) > 0
}
