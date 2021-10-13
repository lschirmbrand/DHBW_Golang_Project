package journal

import (
	"os"
	"time"
)

type location struct {
}

type credentials struct {
	name    string
	address string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func LogToJournal(cred credentials) bool {
	log := cred.address + "," + cred.name + ";\n"
	f, e := os.OpenFile(returnFilename(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	check(e)

	defer f.Close()

	_, e2 := f.WriteString(log)
	check(e2)

	return true
}

func returnFilename() string {
	currentTime := time.Now()
	return "./logs/" + currentTime.Format("01-02-2006") + ".txt"
}
