package journal

import "testing"

func TestLogToJournal(t *testing.T) {
	var cred = credentials{
		address: "address",
		name:    "name",
	}
	LogToJournal(cred)
}
