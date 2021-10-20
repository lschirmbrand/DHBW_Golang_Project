package journal

import (
	"log"
	"os"
	"strings"
	"time"
)

type Credentials struct {
	Login    bool
	Name     string
	Address  string
	Location string
	TimeCome time.Time
	TimeGone time.Time
}

const PATHTOLOGS = "../../logs/log-"
const DATEFORMAT = "2006-01-02"
const DATEFORMATWITHTIME = "02-01-2006 15:04:05"

func check(e error) bool {
	if e != nil {
		log.Fatal(e)
		return false
	}
	return true
}

func LogOutToJournal(cred *Credentials) bool {
	ok := logToJournal(cred, false)
	return ok
}

func LogInToJournal(cred *Credentials) bool {
	ok := logToJournal(cred, true)
	return ok
}

func logToJournal(cred *Credentials, login bool) bool {
	logmsg := returnCreditsToString(cred, login)
	filePath := returnFilename()
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
	if !check(e) {
		return false
	}

	return true
}

func returnFilename() string {
	currentTime := time.Now()
	return PATHTOLOGS + currentTime.Format(DATEFORMAT) + ".txt"
}

func returnCreditsToString(credits *Credentials, isLogin bool) string {
	var sb strings.Builder
	if isLogin {
		sb.WriteString("LOGIN")
	} else {
		sb.WriteString("LOGOUT")
	}
	sb.WriteString(",")
	sb.WriteString(credits.Name)
	sb.WriteString(",")
	sb.WriteString(credits.Address)
	sb.WriteString(",")
	sb.WriteString(credits.Location)
	sb.WriteString(",")
	sb.WriteString(credits.TimeCome.Format(DATEFORMATWITHTIME))
	sb.WriteString(",")
	sb.WriteString(credits.TimeGone.Format(DATEFORMATWITHTIME))
	sb.WriteString(";\n")
	return sb.String()
}