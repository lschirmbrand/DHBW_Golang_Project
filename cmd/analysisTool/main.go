package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
)

type Operation string

const (
	Location Operation = "1"
	Person   Operation = "2"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Analyse Tool")
	fmt.Println("-------------")

	date := dateInputHandler(reader)
	operation := operationInputHandler(reader)

	var sb strings.Builder
	sb.WriteString("../../logs/log-")
	sb.WriteString(date)
	sb.WriteString(".txt")
	filePath := sb.String()
	sb.Reset()

	if _, err := os.Stat(filePath); err == nil {
		content := *readDataFromFile(filePath)
		data := *contentToArray(content)

		if operation == string(Person) {
			analysePersonsForLocation("", &data)
		} else if operation == string(Location) {
			analyseLocationsForPerson("", &data)
		}

	} else {
		fmt.Println("No logs exist for this day.")
	}
}

func dateInputHandler(reader *bufio.Reader) string {
	fmt.Println("Enter Date in format DD-MM-YYYY: ")

	for {
		text, _ := reader.ReadString('\n')
		ok, err := validateDateInput(text)
		check(err)
		if ok {
			return trimStringBasedOnOS(text)
		}
		fmt.Println("Format wrong or pointless. Retry.")
	}
}

func operationInputHandler(reader *bufio.Reader) string {
	fmt.Println("Would you like to analyse:")
	fmt.Println("Locations for a person \t[1]")
	fmt.Println("Visitor for a location \t[2]")

	for {
		text, _ := reader.ReadString('\n')
		ok, err := validateOperationInput(text)
		check(err)
		if ok {
			return trimStringBasedOnOS(text)
		}
		fmt.Println("Input was wrong. Retry.")
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func validateDateInput(date string) (bool, error) {
	return regexp.Match("^(0[1-9]|[12][0-9]|3[01])[-](0[1-9]|1[012])[-](19|20)", []byte(date))
}

func validateOperationInput(operation string) (bool, error) {
	return regexp.Match("\\b[1-2]\\b", []byte(operation))
}

func readDataFromFile(filePath string) *[]string {
	text, err := ioutil.ReadFile(filePath)
	check(err)
	out := strings.Split(string(text), "\n")
	return &out
}

func analyseLocationsForPerson(person string, data *[][]string) {

}

func analysePersonsForLocation(location string, data *[][]string) {

}

func trimStringBasedOnOS(text string) string {
	if runtime.GOOS == "windows" {
		text = strings.TrimSuffix(text, "\n")
		return strings.TrimSuffix(text, "\r")
	}
	return strings.TrimSuffix(text, "\n")
}

func contentToArray(content []string) *[][]string {
	data := make([][]string, len(content))
	for i := 0; i < len(content)-1; i++ {
		row := strings.Trim(content[i], ";")
		values := strings.Split(row, ",")
		data[i] = make([]string, len(values))
		for j, cell := range values {
			data[i][j] = cell
		}
	}
	return &data
}