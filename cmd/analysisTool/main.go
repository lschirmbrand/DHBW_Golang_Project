package main

import (
	"DHBW_Golang_Project/pkg/journal"
	"DHBW_Golang_Project/pkg/location"
	"flag"
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
	PATHTOLOGS                   = "logs/log-"
	PATHTOCSV                    = "logs/export-"
	DATEFORMAT                   = "2006-01-02"
	DATEFORMATWITHTIME           = "02-01-2006 15:04:05"
)

func main() {
	datePtr := flag.String("date", time.Now().Format(DATEFORMAT), "Date of the requested query. Format: YYYY-MM-DD")
	operationPtr := flag.String("operation", string(VISITOR), "Operation of the requested query. Format: Visitor or Location")
	queryPtr := flag.String("query", "", "The keyword of the requested query.")
	flag.Parse()

	args := flag.Args()
	if !requestedHelp(&args) {
		startAnalyticalToolDialog(datePtr, operationPtr, queryPtr)
	} else {
		fmt.Println("go test -date=<DATE> -operation=<VISITOR|LOCATION> -query=<QUERYKEYWORD>")
		fmt.Println("Standardvalue for:")
		fmt.Println("Date:\tDate today")
		fmt.Println("Operation:\tVisitor")
		fmt.Println("Query:\t<none>")
	}
}

func startAnalyticalToolDialog(datePtr *string, operationPtr *string, queryPtr *string) bool {
	var selectedOperation Operation
	if ok, fails := checkFlagFunctionality(datePtr, operationPtr, &selectedOperation, queryPtr); !ok {
		for i := range *fails {
			fmt.Println((*fails)[i])
			return false
		}
	} else {
		fileContent := readDataFromFile(buildFileLogPath(*datePtr))
		loggedCredits := contentToCredits(fileContent)
		var qryResults *[]string
		if strings.EqualFold(*operationPtr, string(VISITOR)) {
			qryResults = analyseLocationsByVisitor(*queryPtr, loggedCredits)
		} else {
			qryResults = analyseVisitorsByLocation(*queryPtr, loggedCredits)
		}

		if assertQueryExport(qryResults) {
			filePath := buildFileCSVPath(selectedOperation, *queryPtr)
			exportToCSVFile(qryResults, *queryPtr, selectedOperation, filePath)
			return true
		}
	}
	return false
}

func check(e error) bool {
	if e != nil {
		log.Fatalln(e)
	}
	return true
}

func analyseLocationsByVisitor(visitor string, data *[]journal.Credentials) *[]string {
	s := make([]string, 0)
	for _, entry := range *data {
		if strings.EqualFold(entry.Name, visitor) {
			s = append(s, string(entry.Location))
		}
	}
	return &s
}

func analyseVisitorsByLocation(location string, data *[]journal.Credentials) *[]string {
	s := make([]string, 0)
	for _, entry := range *data {
		if strings.EqualFold(string(entry.Location), location) {
			s = append(s, entry.Name)
		}
	}
	return &s
}

func contentToCredits(content *[]string) *[]journal.Credentials {
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
		cred.Location = location.Location(strings.ToLower(cells[3]))
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
