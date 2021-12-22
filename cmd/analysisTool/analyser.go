package main

import (
	"DHBW_Golang_Project/internal/config"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func contactHandler(sessions *[]session) *[]contact {
	// Function initiates the contact verification
	contacts := make([]contact, 0)
	for _, entry := range *sessions {
		if strings.EqualFold(entry.Name, *config.Query) {
			newContacts := getOverlaps(&entry, sessions)
			contacts = append(contacts, *newContacts...)
		}
	}
	return &contacts
}

func visitorHandler(sessions *[]session) *[]string {
	// Function initiates the location analysing
	qryResults := analyseLocationsByVisitor(*config.Query, sessions)
	return qryResults
}

func locationHandler(sessions *[]session) *[]string {
	// Function initiates the visitor analysing
	qryResults := analyseVisitorsByLocation(*config.Query, sessions)
	return qryResults
}

func analyseLocationsByVisitor(visitor string, data *[]session) *[]string {
	// Function appends all locations visited by a person
	s := make([]string, 0)
	for _, entry := range *data {
		if strings.EqualFold(entry.Name, visitor) {
			s = append(s, string(entry.Location))
		}
	}
	return &s
}

func analyseVisitorsByLocation(location string, data *[]session) *[]string {
	// Function appends all locations visited by a person
	s := make([]string, 0)
	for _, entry := range *data {
		if strings.EqualFold(string(entry.Location), location) {
			s = append(s, entry.Name)
		}
	}
	return &s
}

func exportContacts(contacts *[]contact) {
	if exportHandler(len(*contacts)) {
		filePath := buildFileCSVPath()
		csvHeader := createCSVHeader()
		writeContactsToCSV(contacts, csvHeader, filePath)
	}
}

func exportLocationsForVisitor(qryResults *[]string) {
	if assertQueryExport(qryResults) {
		filePath := buildFileCSVPath()
		csvHeader := createCSVHeader()
		writeSessionsToCSV(qryResults, csvHeader, filePath)
	}
}

func exportVisitorsForLocation(qryResults *[]string) {
	if !*config.Testcase {
		fmt.Println("\nResults of query for: " + *config.Operation + " = " + *config.Query + ":\n")
		for _, out := range *qryResults {
			fmt.Println(out)
		}
		fmt.Print("\n")
	}

	if assertQueryExport(qryResults) {
		filePath := buildFileCSVPath()
		csvHeader := createCSVHeader()
		writeSessionsToCSV(qryResults, csvHeader, filePath)
	}
}

func assertQueryExport(s *[]string) bool {
	qLen := queryLengthHandler(*s)
	if qLen > 0 {
		if exportHandler(qLen) {
			return true
		} else {
			if !*config.Testcase {
				fmt.Println("Results of query wont get exported. \nAborting.")
			}
			return false
		}
	} else {
		if !*config.Testcase {
			fmt.Println("No results were found for the queried selector.")
		}
		return false
	}
}

func exportHandler(length int) bool {
	if !*config.Testcase {
		fmt.Println("The requested query resulted in ", length, " elements.")
		fmt.Println("Do you want to additionally export the query to csv? [y/n]")
	}
	reader := bufio.NewReader(os.Stdin)
	if *config.Testcase {
		return true
	}
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
