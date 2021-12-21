package journal

import (
	"DHBW_Golang_Project/pkg/location"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

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
	return &LogFileJournal{
		logDir:  logDirName,
		logFile: path.Join(logDirName, "logs-"+time.Now().Format(DATEFORMAT)+".txt"),
	}
}

func check(e error) bool {
	if e != nil {
		log.Fatal(e)
		return false
	}
	return true
}

func (j LogFileJournal) LogOut(cred *Credentials) bool {
	cred.Checkin = false
	ok := j.log(cred)
	return ok
}

func (j LogFileJournal) LogIn(cred *Credentials) bool {
	cred.Checkin = true
	ok := j.log(cred)
	return ok
}

func (j LogFileJournal) log(cred *Credentials) bool {
	logmsg := buildCredits(cred)
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

func buildCredits(credits *Credentials) string {
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
