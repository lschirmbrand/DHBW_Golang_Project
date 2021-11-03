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

type session struct {
	Name     string
	Address  string
	Location location.Location
	TimeCome time.Time
	TimeGone time.Time
}

type contact struct {
	cred     session
	duration time.Duration
}

const (
	LOCATION           Operation = "Location"
	VISITOR            Operation = "Visitor"
	CONTACT            Operation = "Contact"
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
		startAnalyticalTool(datePtr, operationPtr, queryPtr)
	} else {
		fmt.Println("go test -date=<DATE> -operation=<VISITOR|LOCATION> -query=<QUERYKEYWORD>")
		fmt.Println("Standardvalue for:")
		fmt.Println("Date:\tDate today")
		fmt.Println("Operation:\tVisitor")
		fmt.Println("Query:\t<none>")
	}
}

func startAnalyticalTool(datePtr *string, operationPtr *string, queryPtr *string) bool {
	var selectedOperation Operation
	if ok, fails := checkFlagFunctionality(datePtr, operationPtr, &selectedOperation, queryPtr); !ok {
		for i := range *fails {
			fmt.Println((*fails)[i])
			return false
		}
	} else {
		fileContent := readDataFromFile(buildFileLogPath(*datePtr))
		sessions := credentialsToSession(contentToCredits(fileContent))
		if strings.EqualFold(*operationPtr, string(CONTACT)) {
			contacts := make([]contact, 0)
			for _, entry := range *sessions {
				if strings.EqualFold(entry.Name, *queryPtr) {
					newContacts := getOverlaps(&entry, sessions)
					contacts = append(contacts, *newContacts...)
				}

			}

			fmt.Println(contacts)

		} else {
			var qryResults *[]string
			if strings.EqualFold(*operationPtr, string(VISITOR)) {
				qryResults = analyseLocationsByVisitor(*queryPtr, sessions)
			} else {
				qryResults = analyseVisitorsByLocation(*queryPtr, sessions)
			}

			if assertQueryExport(qryResults) {
				filePath := buildFileCSVPath(selectedOperation, *queryPtr)
				exportToCSVFile(qryResults, *queryPtr, selectedOperation, filePath)
				return true
			}
		}
	}
	return false
}

func credentialsToSession(creds *[]journal.Credentials) *[]session {
	sessions := make([]session, 0)
	for _, e := range *creds {
		if e.Login {
			found := false
			for _, eout := range *creds {
				if !eout.Login {
					if e.Name == eout.Name && e.Address == eout.Address && e.Location == eout.Location {
						sessions = append(sessions, session{
							e.Name,
							e.Address,
							e.Location,
							e.Timestamp,
							eout.Timestamp,
						})
						found = true
						break
					}
				}
			}
			if !found {
				sessions = append(sessions, session{
					e.Name,
					e.Address,
					e.Location,
					e.Timestamp,
					time.Now(),
				})
			}
		}
	}
	// fmt.Println(sessions)
	return &sessions
}

func check(e error) bool {
	if e != nil {
		log.Fatalln(e)
	}
	return true
}

func analyseLocationsByVisitor(visitor string, data *[]session) *[]string {
	s := make([]string, 0)
	for _, entry := range *data {
		if strings.EqualFold(entry.Name, visitor) {
			s = append(s, string(entry.Location))
		}
	}
	return &s
}

func analyseVisitorsByLocation(location string, data *[]session) *[]string {
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
		cred.Login = strings.EqualFold(trimStringBasedOnOS(strings.ToLower(cells[0]), false), "checkin")
		cred.Name = cells[1]
		cred.Address = cells[2]
		cred.Location = location.Location(strings.ToLower(cells[3]))
		var err error
		cred.Timestamp, err = time.Parse(DATEFORMATWITHTIME, cells[4])
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

func isOverlapping(entry_1 *session, entry_2 *session) bool {
	return (entry_1.TimeCome.Before(entry_2.TimeGone) && entry_1.TimeGone.After(entry_2.TimeCome)) && strings.EqualFold(entry_1.Name, entry_2.Name)
}

func calculateOverlap(entry_1 *session, entry_2 *session) time.Duration {
	var start time.Time
	var end time.Time

	// Set starttime of contact
	if entry_1.TimeCome.After(entry_2.TimeCome) {
		start = entry_2.TimeCome
	} else {
		start = entry_1.TimeCome
	}

	// Set endtime of contact
	if entry_1.TimeGone.After(entry_2.TimeGone) {
		start = entry_2.TimeGone
	} else {
		start = entry_1.TimeGone
	}

	return end.Sub(start)
}

func getOverlaps(queryEntry *session, entries *[]session) *[]contact {
	contacts := make([]contact, 0)

	for _, entry := range *entries {
		if !isOverlapping(queryEntry, &entry) {
			continue
		}
		newContact := entry
		contacts = append(contacts, contact{
			newContact,
			calculateOverlap(&newContact, queryEntry),
		})
	}

	return &contacts
}
