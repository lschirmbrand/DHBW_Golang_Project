package main

import (
	"DHBW_Golang_Project/pkg/journal"
	"bufio"
	"fmt"
	"log"
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
	LOCATION           Operation = "Location"
	VISITOR            Operation = "Visitor"
	LOCATIONID         Operation = "1"
	VISITORID          Operation = "2"
	PATHTOLOGS                   = "logs/log-"
	PATHTOCSV                    = "logs/export-"
	DATEFORMATWITHTIME           = "02-01-2006 15:04:05"
)

func main() {
	startAnalyticalToolDialog()
}

func check(e error) bool {
	if e != nil {
		log.Fatalln(e)
		return false
	}
	return true
}

func buildFilePath(date string) string {
	return PATHTOLOGS + date + ".txt"
}

func analyseLocationsForPerson(visitor string, data *[]journal.Credentials, reader *bufio.Reader) {
	s := make([]string, 0)
	for _, entry := range *data {
		if strings.EqualFold(entry.Name, visitor) {
			s = append(s, entry.Location)
		}
	}
	assertQueryExport(s, reader, VISITOR, visitor)
}

func analysePersonsForLocation(location string, data *[]journal.Credentials, reader *bufio.Reader) {
	s := make([]string, 0)
	for _, entry := range *data {
		if strings.EqualFold(entry.Location, location) {
			s = append(s, entry.Name)
		}
	}
	assertQueryExport(s, reader, LOCATION, location)
}

func assertQueryExport(s []string, reader *bufio.Reader, operation Operation, selector string) {
	qLen := queryLengthHandler(s)
	if qLen > 0 {
		if exportHandler(reader, qLen) {
			logToCSVFile(s, selector, string(operation))
		} else {
			promptFormatter(1)
			fmt.Println("Results of query wont get exported. \nAborting.")
		}
	} else {
		fmt.Println("No results were found for the queried selector.")
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
	if len(cells) > 1 {
		cred.Login = strings.EqualFold(trimStringBasedOnOS(strings.ToLower(cells[0]), false), "login")
		cred.Name = cells[1]
		cred.Address = cells[2]
		cred.Location = strings.ToLower(cells[3])
		var err error
		cred.TimeCome, err = time.Parse(DATEFORMATWITHTIME, cells[4])
		check(err)
		cred.TimeGone, err = time.Parse(DATEFORMATWITHTIME, cells[5])
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
