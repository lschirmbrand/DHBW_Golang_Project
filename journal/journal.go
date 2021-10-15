package journal

import (
	"os"
	"time"
)

type credentials struct {
	name    string
	address string
}

func check(e error) bool {
	if e != nil {
		panic(e)
		return false
	}
	return true
}

func LogToJournal(cred credentials, isTest bool) bool {
	log := cred.address + "," + cred.name + ";\n"
	var filePath string
	if !isTest {
		filePath = returnFilename()
	} else {
		filePath = "../logs/log-temp-test-file.txt"
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
	return "../logs/log-" + currentTime.Format("02-01-2006") + ".txt"
}
