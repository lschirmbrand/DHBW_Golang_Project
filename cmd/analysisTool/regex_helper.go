package main

import "regexp"

func validateDateInput(date string) (bool, error) {
	return regexp.Match("^([19|20].(0[1-9]|[1-9][1-9]))[-](0[1-9]|1[012])[-](0[1-9]|[12][0-9]|3[01])$", []byte(date))
}

func validateOperationInput(operation string) (bool, error) {
	return regexp.Match("\\b[1-2]\\b", []byte(operation))
}

func validateYesNoInput(operation string) (bool, error){
	return regexp.Match("(?i)\\b[y|n]\\b|yes\\b|no\\b", []byte(operation))
}