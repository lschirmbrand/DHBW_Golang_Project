package journal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLogToJournal(t *testing.T) {
	var cred = credentials{
		address: "address",
		name:    "name",
	}
	LogToJournal(cred, true)
	filePath := "../logs/log-temp-test-file.txt"
	data, e := os.ReadFile(filePath)
	defer os.Remove(filePath)
	check(e)

	assert.EqualValues(t, string(data), cred.address + "," + cred.name + ";\n")
}

func TestCheck
