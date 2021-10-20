package journal

import (
	"os"
	"strings"
	"time"
)

type Credentials struct {
	Name     string
	Address  string
	Location string
	TimeCome time.Time
	TimeGone time.Time
}

func check(e error) bool {
	if e != nil {
		panic(e)
		return false
	}
	return true
}

func LogToJournal(cred *Credentials, isTest bool) bool {
	log := returnCreditsToString(cred)
	var filePath string
	if !isTest {
		filePath = returnFilename()
	} else {
		filePath = "../../logs/log-temp-test-file.txt"
	}
	f, e := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if !check(e) {
		return false
	}

	defer f.Close()

	_, e2 := f.WriteString(log)
	if !check(e2) {
		return false
	}

	return true
}

func returnFilename() string {
	currentTime := time.Now()
	return "../../logs/log-" + currentTime.Format("02-01-2006") + ".txt"
}

func returnCreditsToString(credits *Credentials) string {
	var sb strings.Builder
	sb.WriteString(credits.Name)
	sb.WriteString(",")
	sb.WriteString(credits.Address)
	sb.WriteString(",")
	sb.WriteString(credits.Location)
	sb.WriteString(",")
	sb.WriteString(credits.TimeCome.Format("02-01-2006 15:04:05"))
	sb.WriteString(",")
	sb.WriteString(credits.TimeGone.Format("02-01-2006 15:04:05"))
	sb.WriteString(";\n")
	return sb.String()
}
