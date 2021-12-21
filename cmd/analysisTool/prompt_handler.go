package main

import (
	"DHBW_Golang_Project/pkg/config"
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func exportContacts(contacts *[]contact){
	if exportHandler(len(*contacts)) {
		filePath := buildFileCSVPath()
		csvHeader := createCSVHeader()
		writeContactsToCSV(contacts, csvHeader, filePath)
	}
}

func exportLocations(qryResults *[]string){
	if assertQueryExport(qryResults) {
		filePath := buildFileCSVPath()
		csvHeader := createCSVHeader()
		writeSessionsToCSV(qryResults, filePath, csvHeader)
	}
}

func exportVisitors(qryResults *[]string){
	fmt.Println("\nResults of query for: " + *config.Operation + " = " + *config.Query + ":\n")
	for _, out := range *qryResults {
		fmt.Println(out)
	}
	fmt.Print("\n")

	if assertQueryExport(qryResults) {
		filePath := buildFileCSVPath()
		csvHeader := createCSVHeader()
		writeSessionsToCSV(qryResults, filePath, csvHeader)
	}
}

func assertQueryExport(s *[]string) bool {
	qLen := queryLengthHandler(*s)
	if qLen > 0 {
		if exportHandler(qLen) {
			return true
		} else {
			fmt.Println("Results of query wont get exported. \nAborting.")
			return false
		}
	} else {
		fmt.Println("No results were found for the queried selector.")
		return false
	}
}

func exportHandler(length int) bool {
	fmt.Println("The requested query resulted in ", length, " elements.")
	fmt.Println("Do you want to additionally export the query to csv? [y/n]")
	reader := bufio.NewReader(os.Stdin)
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

func trimStringBasedOnOS(text string, isSuffix bool) string {
	isWindows := runtime.GOOS == "windows"
	if isSuffix {
		if isWindows {
			text = strings.TrimSuffix(text, "\x0d")
			text = strings.TrimSuffix(text, "\x0a")
			text = strings.TrimSuffix(text, "\n")
			return strings.TrimSuffix(text, "\r")
		}
		text = strings.TrimSuffix(text, "\x0d")
		return strings.TrimSuffix(text, "\n")
	} else {
		text = strings.TrimPrefix(text, "\x0d")
		return strings.TrimPrefix(text, "\n")
	}
}
