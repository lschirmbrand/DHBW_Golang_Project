package main

import (
	"DHBW_Golang_Project/internal/config"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func buildFileLogPath(date string) string {
	return path.Join(*config.LogPath, "logs-"+date+".txt")
}

func buildFileCSVPath() string {
	return path.Join(*config.LogPath, "export-"+*config.Operation+"_"+*config.Query+".csv")
}

func readDataFromFile(filePath string) *[]string {
	text, err := ioutil.ReadFile(filePath)
	check(err)
	out := strings.Split(string(text), "\n")
	for i, _ := range out {
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
	f, e := os.Create(filePath)
	check(e)
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()
	e = writer.Write(*csvHeader)
	check(e)
	e = writer.Write(*results)
	check(e)
	fmt.Println("The query-result was exported to: " + filePath)
}

func writeContactsToCSV(contacts *[]contact, csvHeader *[]string, filePath string) {
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
			fmt.Println("Failed to write, aborting!")
			return
		}

	}
	fmt.Println("The query-result was exported to: " + filePath)
}

func createCSVHeader() *[]string {
	var infix string
	switch *config.Operation {
	case string(LOCATION):
		fallthrough
	case string(VISITOR):
		infix = "Results for the query: " + *config.Operation + " = "
	case string(CONTACT):
		infix = *config.Operation + " for the user: "
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
