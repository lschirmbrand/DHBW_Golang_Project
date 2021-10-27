package journal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestLogToJournal(t *testing.T) {
	var cred = Credentials{
		Address:  "Address",
		Name:     "name",
		Location: "Location",
		TimeCome: time.Now(),
		TimeGone: time.Now(),
	}
	LogInToJournal(&cred)
	filePath := "../../logs/log-temp-test-file.txt"
	data, e := os.ReadFile(filePath)
	defer os.Remove(filePath)
	check(e)

	assert.EqualValues(t, string(data), cred.Name+","+cred.Address+","+cred.Location+","+cred.TimeCome.Format(DATEFORMATWITHTIME)+","+cred.TimeGone.Format(DATEFORMATWITHTIME)+";\n")
}

func TestReturnCreditsToString(t *testing.T) {
	var cred = Credentials{
		Address:  "Address",
		Name:     "name",
		Location: "Location",
		TimeCome: time.Now(),
		TimeGone: time.Now(),
	}
	assert.EqualValues(t, returnCreditsToString(&cred, true), "in,"+cred.Name+","+cred.Address+","+cred.Location+","+cred.TimeCome.Format(DATEFORMATWITHTIME)+","+cred.TimeGone.Format(DATEFORMATWITHTIME)+";\n")
}

func TestLogTestExample(t *testing.T) {
	var cred = Credentials{
		Address:  "Address",
		Name:     "Name",
		Location: "Location",
		TimeCome: time.Now(),
		TimeGone: time.Now(),
	}
	for i := 0; i < 500; i++ {
		LogInToJournal(&cred)
	}
}
