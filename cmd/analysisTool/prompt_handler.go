package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startAnalyticalToolDialog() {
	reader := bufio.NewReader(os.Stdin)
	promptFormatter(1)
	fmt.Println("-         Analyse-Tool         -")
	fmt.Println("--------------------------------")

	var filePath string

	for {
		date := dateInputHandler(reader)
		filePath = buildFilePath(date)
		if _, err := os.Stat(filePath); err == nil {
			break
		} else {
			fmt.Println("No logs exist for the entered date.")
			fmt.Println("Retry with different date, or abort.")
		}
	}

	operation := operationInputHandler(reader)
	content := *readDataFromFile(filePath)
	data := *contentToArray(&content)

	if operation == string(VISITORID) {
		request, ok := searchRequestHandler(reader, string(LOCATION))
		if ok {
			analysePersonsForLocation(request, &data, reader)
		}
	} else if operation == string(LOCATIONID) {
		request, ok := searchRequestHandler(reader, string(VISITOR))
		if ok {
			analyseLocationsForPerson(request, &data, reader)
		}
	}
}

func dateInputHandler(reader *bufio.Reader) string {
	promptFormatter(1)
	fmt.Println("Enter date in format YYYY-MM-DD: ")

	for {
		text, _ := reader.ReadString('\n')
		text = trimStringBasedOnOS(text, true)
		ok, err := validateDateInput(text)
		check(err)
		if ok {
			return text
		}
		fmt.Println("Format wrong or pointless. Retry.")
	}
}

func operationInputHandler(reader *bufio.Reader) string {
	promptFormatter(1)
	fmt.Println("Would you like to analyse:")
	fmt.Println("Locations for a " + string(VISITOR) + "\t[1]")
	fmt.Println("Visitors for a " + string(LOCATION) + "\t[2]")

	for {
		text, _ := reader.ReadString('\n')
		ok, err := validateOperationInput(text)
		check(err)
		if ok {
			return trimStringBasedOnOS(text, true)
		}
		fmt.Println("Input was wrong. Retry.")
	}
}

func searchRequestHandler(reader *bufio.Reader, operation string) (string, bool) {
	promptFormatter(1)
	fmt.Println("You requested to search by: " + operation)
	fmt.Println("Please enter the keyword you are searching for:")
	input, e := reader.ReadString('\n')
	if check(e) {
		return trimStringBasedOnOS(input, true), true
	}
	return "", false
}

func exportHandler(reader *bufio.Reader, length int) bool {
	promptFormatter(1)
	fmt.Println("The requested query resulted in ", length, " elements.")
	fmt.Println("Do you want to export the query? [y/n]")
	for {
		input, e := reader.ReadString('\n')
		check(e)
		ok, e := validateYesNoInput(input)
		check(e)
		if ok {
			return strings.EqualFold(trimStringBasedOnOS(input, true), "y")
		}
		fmt.Println("Input was incorrect, retry.")
	}
}

func queryLengthHandler(slice []string) int {
	return len(slice)
}

func promptFormatter(newlines int){
	for i := 0; i < newlines; i++ {
		fmt.Println()
	}
}
