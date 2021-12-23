package journal

import (
	"DHBW_Golang_Project/internal/location"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

/*
	Erstellt von: 	4775194
	Created by:		4775194

	also: 8864957, 9514094
*/

type Journal interface {
	LogIn(cred *Credentials) bool
	LogOut(cred *Credentials) bool
}

type Credentials struct {
	Checkin   bool
	Name      string
	Address   string
	Location  location.Location
	Timestamp time.Time
}

const DATEFORMAT = "2006-01-02"
const DATEFORMATWITHTIME = "02-01-2006 15:04:05"

type LogFileJournal struct {
	logDir  string
	logFile string
}

func NewLogFileJournal(logDirName string) *LogFileJournal {
	/*
		Public accessible function that returns a structure, which contains
		the log-directory and the log-path.
	 */
	return &LogFileJournal{
		logDir:  logDirName,
		logFile: path.Join(logDirName, "logs-"+time.Now().Format(DATEFORMAT)+".txt"),
	}
}

func (j LogFileJournal) LogOut(cred *Credentials) bool {
	/*
		Public accessible function, which calls the function "log" and
		only sets the "Checkin" value to false.
	*/
	cred.Checkin = false
	ok := j.log(cred)
	return ok
}

func (j LogFileJournal) LogIn(cred *Credentials) bool {
	/*
		Public accessible function, which calls the function "log" and
		only sets the "Checkin" value to true.
	 */
	cred.Checkin = true
	ok := j.log(cred)
	return ok
}

func (j LogFileJournal) log(cred *Credentials) bool {
	/*
		The function is an internal function, which is only accessible inside
		the package. It gets called from the two functions LogIn and LogOut,
		which only change the "Checkin" value of the credentials element
	*/
	logmsg := buildCreditString(cred)
	if _, err := os.Stat(j.logDir); os.IsNotExist(err) {
		os.MkdirAll(j.logDir, 0755)
	}

	f, e := os.OpenFile(j.logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if !check(e) {
		return false
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(f)

	_, e = f.WriteString(logmsg)
	return check(e)
}

func buildCreditString(credits *Credentials) string {
	/*
		Passed credit gets transformed into a string with a string-builder,
		which will be returned and exported to the log-file.
	 */
	var sb strings.Builder
	switch credits.Checkin {
	case true:
		sb.WriteString("CHECKIN")
	default:
		sb.WriteString("CHECKOUT")
	}
	sb.WriteString(",")
	sb.WriteString(credits.Name)
	sb.WriteString(",")
	sb.WriteString(credits.Address)
	sb.WriteString(",")
	sb.WriteString(string(credits.Location))
	sb.WriteString(",")
	sb.WriteString(credits.Timestamp.Format(DATEFORMATWITHTIME))
	sb.WriteString(";\n")
	return sb.String()
}

func check(e error) bool {
	if e != nil {
		log.Fatal(e)
		return false
	}
	return true
}
