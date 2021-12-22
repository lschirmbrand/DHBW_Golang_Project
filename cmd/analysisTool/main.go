package main

import (
	"DHBW_Golang_Project/internal/config"
	"DHBW_Golang_Project/internal/journal"
	"DHBW_Golang_Project/internal/location"
	"fmt"
	"log"
	"strings"
	"time"
)

type Operation string

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
			if !*config.Testcase {
				fmt.Println((*fails)[i])
			}
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
			// Query contacts and export them
			contacts := contactHandler(sessions)
			exportContacts(contacts)
		case string(VISITOR):
			// Query locations and export their visitors
			locations := visitorHandler(sessions)
			exportLocationsForVisitor(locations)
		default:
			// Query visitors and export their locations
			visitors := locationHandler(sessions)
			exportVisitorsForLocation(visitors)
		}
	}
	return false
}

func credentialsToSession(credentials *[]journal.Credentials) *[]session {
	/*
		The function iterates over all credentials of the log-file.
		This happens in chronological order, so that in case of
		multiple visits from a person at the same day, the checkins
		and checkouts won't get mixed up.
		In before this was implemented with workers, but the already
		described problem caused the change of structure.
	*/
	sessions := make([]session, 0)
	for _, e := range *credentials {
		if e.Checkin {
			/*
				If the matching checkout for the queried checkin was found,
				a session will be created with all information of the visit
			 */
			found := false
			for _, eout := range *credentials {
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
			/*
				If no matching checkout for the queried checkin was found,
				the session is still active and therefore will be treated
				for the analyser as the current time.
			*/
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

func contentToCredits(content *[]string) *[]journal.Credentials {
	/*
		Each row of the log-file will be parsed to a
		credential-entry
	 */
	data := make([]journal.Credentials, len(*content))
	for i, row := range *content {
		data[i] = splitDataRowToCells(row)
	}

	return &data
}

func splitDataRowToCells(row string) journal.Credentials {
	/*
		Each row will be split by commas, as the string
		separator. The linebreak Separator is getting
		removed. If a row results in more than one
		element, it contains content.
	*/
	var cred journal.Credentials
	row = strings.Trim(row, ";")
	cells := strings.Split(row, ",")
	if len(cells) > 1 {
		/*
			Each cell of the string represents a cell in the credentials.
			Therefore, the credentials gets filled.
		 */
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

func isOverlapping(entry1 *session, entry2 *session) bool {
	/*
		The function checks for an eventual overlap between two
		passed session. Returns true, if they overlap.
	 */
	return ((entry1.TimeCome.Before(entry2.TimeGone) && entry1.TimeGone.After(entry2.TimeCome)) || entry1.TimeCome.Equal(entry2.TimeCome)) && strings.EqualFold(string(entry1.Location), string(entry2.Location)) && !(strings.EqualFold(entry1.Name, entry2.Name))
}

func calculateOverlap(entry1 *session, entry2 *session) time.Duration {
	/*
		If an overlap of the two session was detected, the duration
		of the contact is important. Therefore, this duration has to
		will be calculated in this function.
	 */
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
	// The starttime gets subtracted from the endtime --> Duration
	return end.Sub(start)
}

func getOverlaps(queryEntry *session, entries *[]session) *[]contact {
	/*
		This function launches the contact verification, collects all
		contacts for one person (or their session). All contacts get returned.
	 */
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