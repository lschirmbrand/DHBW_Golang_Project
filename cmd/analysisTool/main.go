package main

import (
	"DHBW_Golang_Project/internal/config"
	"DHBW_Golang_Project/internal/journal"
	"DHBW_Golang_Project/internal/location"
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

type contact struct {
	session  session
	duration time.Duration
}

type session struct {
	Name     string
	Address  string
	Location location.Location
	TimeCome time.Time
	TimeGone time.Time
}

const (
	LOCATION Operation = "Location"
	VISITOR  Operation = "Visitor"
	CONTACT  Operation = "Contact"
)

func main() {
	/*
		Setting standard flags and overwrite
		those, which were passed as arguments
	*/
	config.ConfigureAnalysisTool()
	startAnalyticalToolDialog()
}

func startAnalyticalToolDialog() bool {
	// Checks whether the Flags are set up right
	if ok, fails := checkFlagFunctionality(); !ok {
		for i := range *fails {
			fmt.Println((*fails)[i])
			return false
		}
	} else {
		/*
			If the flags are set up right, the content
			of the Logfile will be imported
		*/
		fileContent := readDataFromFile(buildFileLogPath(*config.Date))
		/*
			The imported content, formatted as logs will
			be parsed to the session-structure
		*/
		sessions := credentialsToSession(contentToCredits(fileContent))
		/*
			Switching through the possible Use-Cases
			for the operations
		*/
		switch *config.Operation {
		case string(CONTACT):
			contacts := contactHandler(sessions)
			exportContacts(contacts)
		case string(VISITOR):
			locations := visitorHandler(sessions)
			exportLocations(locations)
		default:
			visitors := locationHandler(sessions)
			exportVisitors(visitors)
		}
	}
	return false
}

func contactHandler(sessions *[]session) *[]contact {
	contacts := make([]contact, 0)
	for _, entry := range *sessions {
		if strings.EqualFold(entry.Name, *config.Query) {
			newContacts := getOverlaps(&entry, sessions)
			contacts = append(contacts, *newContacts...)
		}
	}
	return &contacts
}

func visitorHandler(sessions *[]session) *[]string {
	qryResults := analyseLocationsByVisitor(*config.Query, sessions)
	return qryResults
}

func locationHandler(sessions *[]session) *[]string {
	qryResults := analyseVisitorsByLocation(*config.Query, sessions)
	return qryResults
}

func credentialsToSession(creds *[]journal.Credentials) *[]session {
	sessions := make([]session, 0)
	for _, e := range *creds {
		if e.Checkin {
			found := false
			for _, eout := range *creds {
				if !eout.Checkin {
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
		cred.Checkin = strings.EqualFold(trimStringBasedOnOS(strings.ToLower(cells[0]), false), "checkin")
		cred.Name = cells[1]
		cred.Address = cells[2]
		cred.Location = location.Location(strings.ToLower(cells[3]))
		var err error
		cred.Timestamp, err = time.Parse(config.DATEFORMATWITHTIME, trimStringBasedOnOS(cells[4], true))
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

func isOverlapping(entry1 *session, entry2 *session) bool {
	return ((entry1.TimeCome.Before(entry2.TimeGone) && entry1.TimeGone.After(entry2.TimeCome)) || entry1.TimeCome.Equal(entry2.TimeCome)) && strings.EqualFold(string(entry1.Location), string(entry2.Location)) && !(strings.EqualFold(entry1.Name, entry2.Name))
}

func calculateOverlap(entry1 *session, entry2 *session) time.Duration {
	var start time.Time
	var end time.Time

	// Set starttime of contact
	if entry1.TimeCome.After(entry2.TimeCome) {
		start = entry1.TimeCome
	} else {
		start = entry2.TimeCome
	}

	// Set endtime of contact
	if entry1.TimeGone.After(entry2.TimeGone) {
		end = entry2.TimeGone
	} else {
		end = entry1.TimeGone
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
