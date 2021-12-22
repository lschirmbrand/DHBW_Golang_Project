package main

import (
	"DHBW_Golang_Project/internal/config"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
)

func buildFileLogPath(date string) string {
	// Helper-Function for the creation of the logpath
	return path.Join(*config.LogPath, "logs-"+date+".txt")
}

func buildFileCSVPath() string {
	// Helper-Function for the creation of the exportpath
	return path.Join(*config.LogPath, "export-"+*config.Operation+"_"+*config.Query+".csv")
}

func readDataFromFile(filePath string) *[]string {
	// Function reads the content of a file
	text, err := ioutil.ReadFile(filePath)
	check(err)
	out := strings.Split(string(text), "\n")
	for i := range out {
		out[i] = trimStringBasedOnOS(out[i], true)
		if i == (len(out) - 1) {
			out[i] = strings.TrimSuffix(out[i], ";")
		}
	}
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	return &out
}

func writeSessionsToCSV(results *[]string, csvHeader *[]string, filePath string) {
	if _, err := os.Stat(*config.LogPath); os.IsNotExist(err) {
		os.MkdirAll(*config.LogPath, 0755)
	}
	f, e := os.Create(filePath)
	check(e)
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()
	e = writer.Write(*csvHeader)
	check(e)
	e = writer.Write(*results)
	check(e)
	if !(*config.Testcase) {
		fmt.Println("The query-result was exported to: " + filePath)
	}
}

func writeContactsToCSV(contacts *[]contact, csvHeader *[]string, filePath string) {
	if _, err := os.Stat(*config.LogPath); os.IsNotExist(err) {
		os.MkdirAll(*config.LogPath, 0755)
	}
	f, e := os.Create(filePath)
	check(e)
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	e = writer.Write(*csvHeader)
	check(e)
	for _, contact := range *contacts {
		row := contact.toSlice()
		if err := writer.Write(*row); err != nil {
			if !*config.Testcase {
				fmt.Println("Failed to write, aborting!")
			}
			return
		}

	}
	if !*config.Testcase {
		fmt.Println("The query-result was exported to: " + filePath)
	}
}

func createCSVHeader() *[]string {
	/*
		The function is used to create the caption, which
		gets printed above the csv-export
	 */
	var infix string
	switch *config.Operation {
	case string(LOCATION):
		fallthrough
	case string(VISITOR):
		infix = "Results for the query: " + *config.Operation + " = "
	case string(CONTACT):
		infix = *config.Operation + "s for the user: "
	}
	csvHeader := make([]string, 1)
	csvHeader[0] = infix + *config.Query
	return &csvHeader
}

func (contact *contact) toSlice() *[]string {
	field := make([]string, 3)
	field[0] = contact.session.Name
	field[1] = string(contact.session.Location)
	field[2] = contact.duration.String()
	return &field
}

func checkFileExistence(path string) bool {
	// Helper function to check whether a file exists
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func trimStringBasedOnOS(text string, isSuffix bool) string {
	// Function is used to trim linebreaks based on the OS
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
