package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func readDataFromFile(filePath string) *[]string {
	text, err := ioutil.ReadFile(filePath)
	check(err)
	out := strings.Split(string(text), "\n")
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	return &out
}

func writeSessionsToCSV(results *[]string, selector string, operation Operation, filePath string) {
	f, e := os.Create(filePath)
	check(e)
	defer f.Close()

	csvHeader := make([]string, 1)
	csvHeader[0] = "Results for: " + selector

	w := csv.NewWriter(f)
	e = w.Write(csvHeader)
	check(e)
	e = w.Write(*(results))
	w.Flush()
	check(e)

	fmt.Println("The query-result was exported to: " + filePath)
}

func writeContactsToCSV(contacts *[]contact, filePath string) {
	f, e := os.Create(filePath)
	check(e)
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	for _, contact := range *contacts {
		row := contact.toSlice()
		if err := writer.Write(*row); err != nil {
			fmt.Println("Failed to write, aborting!")
			return
		}

	}
	fmt.Println("The query-result was exported to: " + filePath)
}

func buildFileLogPath(date string) string {
	return PATHTOLOGS + date + ".txt"
}

func buildFileCSVPath(operation Operation, selector string) string {
	return PATHTOCSV + string(operation) + "_" + selector + ".csv"
}

func (contact *contact) toSlice() *[]string {
	field := make([]string, 3)
	field[0] = contact.session.Name
	field[1] = string(contact.session.Location)
	field[2] = contact.duration.String()
	return &field
}
