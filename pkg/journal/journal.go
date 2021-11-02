package journal

import (
	"DHBW_Golang_Project/pkg/location"
	"flag"
	"log"
	"os"
	"strings"
	"time"
)

type Credentials struct {
	Login    bool
	Name     string
	Address  string
	Location location.Location
	TimeCome time.Time
	TimeGone time.Time
}

const PATHTOLOGS = "logs"
const DATEFORMAT = "2006-01-02"
const DATEFORMATWITHTIME = "02-01-2006 15:04:05"

var (
	LogPath     *string = flag.String("filepath", PATHTOLOGS, "The filepath to the log-Directory.")
	LogFilename *string = flag.String("filename", time.Now().Format(DATEFORMAT), "The filename of the log-file.")
)

func check(e error) bool {
	if e != nil {
		log.Fatal(e)
		return false
	}
	return true
}

func LogOutToJournal(cred *Credentials) bool {
	parseFlags()
	cred.Login = false
	ok := logToJournal(cred)
	return ok
}

func LogInToJournal(cred *Credentials) bool {
	parseFlags()
	cred.Login = true
	ok := logToJournal(cred)
	return ok
}

func logToJournal(cred *Credentials) bool {
	logmsg := buildCredits(cred)
	filePath := returnFilepath()
	if _, err := os.Stat(*LogPath); !os.IsNotExist(err) {
		f, e := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

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
	} else {
		check(err)
	}
	return false
}

func returnFilepath() string {
	return *LogPath + "/logs-" + *LogFilename + ".txt"
}

func buildCredits(credits *Credentials) string {
	var sb strings.Builder
	switch credits.Login {
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
	sb.WriteString(credits.TimeCome.Format(DATEFORMATWITHTIME))
	sb.WriteString(",")
	sb.WriteString(credits.TimeGone.Format(DATEFORMATWITHTIME))
	sb.WriteString(";\n")
	return sb.String()
}

func parseFlags() {
	flag.Parse()
}
