package journal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogToJournal(t *testing.T) {
	var cred = Credentials{
		Address: "address",
		Name:    "name",
	}
	LogToJournal(cred, true)
	filePath := "../logs/log-temp-test-file.txt"
	data, e := os.ReadFile(filePath)
	defer os.Remove(filePath)
	check(e)

	assert.EqualValues(t, string(data), cred.Address+","+cred.Name+";\n")
}
