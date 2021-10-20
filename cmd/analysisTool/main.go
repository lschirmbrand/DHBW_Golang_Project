package main

import (
	"DHBW_Golang_Project/pkg/journal"
	"bufio"
	"fmt"
	"io/ioutil"
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
	Location Operation = "1"
	Person   Operation = "2"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Analyse Tool")
	fmt.Println("-------------")

	date := dateInputHandler(reader)
	operation := operationInputHandler(reader)
	filePath := buildFilePath(date)

	if _, err := os.Stat(filePath); err == nil {
		content := *readDataFromFile(filePath)
		data := *contentToArray(&content)

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

func buildFilePath(date string) string {
	var sb strings.Builder
	sb.WriteString("../../logs/log-")
	sb.WriteString(date)
	sb.WriteString(".txt")
	return sb.String()
}

func analyseLocationsForPerson(person string, data *[]journal.Credentials) {

}

func analysePersonsForLocation(location string, data *[]journal.Credentials) {

}

func trimStringBasedOnOS(text string) string {
	if runtime.GOOS == "windows" {
		text = strings.TrimSuffix(text, "\n")
		return strings.TrimSuffix(text, "\r")
	}
	return strings.TrimSuffix(text, "\n")
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
		cred.Name = cells[0]
		cred.Address = cells[1]
		cred.Location = cells[2]
		var err error
		cred.TimeCome, err = time.Parse("02-01-2006 15:04:05", cells[3])
		check(err)
		cred.TimeGone, err = time.Parse("02-01-2006 15:04:05", cells[4])
		check(err)

	}
	return cred
}

func jobFactory(content []string) <-chan job {
	jobs := make(chan job)
	go func() {
		for i := 0; i < len(content)-1; i++ {
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
