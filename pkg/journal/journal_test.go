package journal

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

func TestLogToJournal(t *testing.T) {
	var cred = Credentials{
		Login:    true,
		Address:  "Address",
		Name:     "name",
		Location: "Location",
		TimeCome: time.Now(),
		TimeGone: time.Now(),
	}
	LogInToJournal(&cred)
	filePath := "../../logs/log-temp-test-file.txt" //<- Schlägt fehl, aufgrund momentaner Architektur, wird gefixxt
	data, e := os.ReadFile(filePath)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatalln(err)
		}
	}(filePath)
	check(e)

	assert.EqualValues(t, string(data), "LOGIN,"+cred.Name+","+cred.Address+","+cred.Location+","+cred.TimeCome.Format(DATEFORMATWITHTIME)+","+cred.TimeGone.Format(DATEFORMATWITHTIME)+";\n")
}

func TestReturnCreditsToString(t *testing.T) {
	var cred = Credentials{
		Address:  "Address",
		Name:     "name",
		Location: "Location",
		TimeCome: time.Now(),
		TimeGone: time.Now(),
	}
	assert.EqualValues(t, returnCreditsToString(&cred, true), "LOGIN,"+cred.Name+","+cred.Address+","+cred.Location+","+cred.TimeCome.Format(DATEFORMATWITHTIME)+","+cred.TimeGone.Format(DATEFORMATWITHTIME)+";\n")
}

func TestLogTestExample(t *testing.T) {
	var cred = Credentials{
		Login:    true,
		Name:     "Name",
		Address:  "Address",
		Location: "Location",
		TimeCome: time.Now(),
		TimeGone: time.Now(),
	}
	for i := 0; i < 500; i++ {
		LogInToJournal(&cred)
	}
}
