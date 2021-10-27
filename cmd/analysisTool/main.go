package main

import (
	"DHBW_Golang_Project/pkg/journal"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type Operation string

type job struct {
	content string
}

type result struct {
	cred journal.Credentials
}

const (
	Location           Operation = "1"
	Person             Operation = "2"
	PATHTOLOGS                   = "logs/log-"
	DATEFORMATWITHTIME           = "02-01-2006 15:04:05"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Analyse Tool")
	fmt.Println("-------------")

	date := dateInputHandler(reader)
	filePath := buildFilePath(date)

	if _, err := os.Stat(filePath); err == nil {
		operation := operationInputHandler(reader)
		content := *readDataFromFile(filePath)
		data := *contentToArray(&content)

		if operation == string(Person) {
			request, ok := searchRequestHandler(reader, "person")
			if ok {
				analysePersonsForLocation(request, &data)
			}
		} else if operation == string(Location) {
			request, ok := searchRequestHandler(reader, "location")
			if ok {
				analyseLocationsForPerson(request, &data)
			}
		}

	} else {
		fmt.Println("No logs exist for this day.")
	}
}

func dateInputHandler(reader *bufio.Reader) string {
	fmt.Println("Enter Date in format YYYY-MM-DD: ")

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
	fmt.Println("Would you like to analyse:")
	fmt.Println("Locations for a person \t[1]")
	fmt.Println("Visitor for a location \t[2]")

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
	fmt.Println("You requested to search by: " + operation)
	fmt.Println("Please enter the keyword you are searching for:")
	input, e := reader.ReadString('\n')
	if check(e) {
		return trimStringBasedOnOS(input, true), true
	}
	return "", false
}

func check(e error) bool {
	if e != nil {
		log.Fatalln(e)
		return false
	}
	return true
}

func validateDateInput(date string) (bool, error) {
	return regexp.Match("^(([19|20].(0[1-9]|[1-9][1-9])))[-](0[1-9]|1[012])[-](0[1-9]|[12][0-9]|3[01])$", []byte(date))
}

func validateOperationInput(operation string) (bool, error) {
	return regexp.Match("\\b[1-2]\\b", []byte(operation))
}

func readDataFromFile(filePath string) *[]string {
	text, err := ioutil.ReadFile(filePath)
	check(err)
	out := strings.Split(string(text), "\x0d")
	if len(out) > 0 {
		out = out[:len(out)-1]
		fmt.Println(len(out))
	}
	return &out
}

func buildFilePath(date string) string {
	return PATHTOLOGS + date + ".txt"
}

func analyseLocationsForPerson(person string, data *[]journal.Credentials) {
	s := make([]journal.Credentials, 0)

	for _, entry := range *data {
		if strings.EqualFold(entry.Name, person) {
			fmt.Println(entry.Name, entry.Location)
			s = append(s, entry)
		}
	}

	for _, entry := range s {
		fmt.Println(entry.Name, entry.Location)
	}
}

func analysePersonsForLocation(location string, data *[]journal.Credentials) {

}

func trimStringBasedOnOS(text string, isSuffix bool) string {
	isWindows := runtime.GOOS == "windows"
	if isSuffix {
		if isWindows {
			text = strings.TrimSuffix(text, "\x0a\x0d")
			return strings.TrimSuffix(text, "\r\n")
		}
		text = strings.TrimSuffix(text, "\x0d")
		return strings.TrimSuffix(text, "\n")
	} else {
		text = strings.TrimPrefix(text, "\x0d")
		return strings.TrimPrefix(text, "\n")
	}
}

func contentToArray(content *[]string) *[]journal.Credentials {
	data := make([]journal.Credentials, len(*content))
	jobs := jobFactory(*content)
	results, imageDone := resultCollector(&data)

	workersDone := make(chan bool)
	workers := runtime.NumCPU()
	for i := 0; i < workers; i++ {
		go worker(jobs, results, workersDone)
	}

	for i := 0; i < workers; i++ {
		<-workersDone
	}

	close(results)
	<-imageDone
	return &data
}

func splitDataRowToCells(row string) journal.Credentials {
	var cred journal.Credentials
	row = strings.Trim(row, ";")
	cells := strings.Split(row, ",")
	if len(cells) > 0 {
		cred.Name = trimStringBasedOnOS(strings.ToLower(cells[0]), false)
		cred.Address = cells[1]
		cred.Location = strings.ToLower(cells[2])
		var err error
		cred.TimeCome, err = time.Parse(DATEFORMATWITHTIME, cells[3])
		check(err)
		cred.TimeGone, err = time.Parse(DATEFORMATWITHTIME, cells[4])
		check(err)
	}
	return cred
}

func jobFactory(content []string) <-chan job {
	jobs := make(chan job)
	go func() {
		for i := 0; i < len(content); i++ {
			jobs <- job{content[i]}
		}
		close(jobs)
	}()
	return jobs
}

func worker(jobs <-chan job, results chan<- result, done chan<- bool) {
	for job := range jobs {
		results <- result{splitDataRowToCells(job.content)}
	}
	done <- true
}

func resultCollector(data *[]journal.Credentials) (chan<- result, <-chan bool) {
	results := make(chan result)
	done := make(chan bool)
	go func() {
		i := 0
		for result := range results {
			(*data)[i] = result.cred
			i++
		}
		done <- true
	}()

	return results, done
}
