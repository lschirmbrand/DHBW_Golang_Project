package journal

import (
	"DHBW_Golang_Project/pkg/config"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogToJournal(t *testing.T) {
	var cred = Credentials{
		Login:    true,
		Address:  "Address",
		Name:     "Name",
		Location: "Location",
		TimeCome: time.Now(),
		TimeGone: time.Now(),
	}

	modifyFlagsForTestCase(true, true)

	filePath := returnFilepath()
	logToJournal(&cred)
	data, e := os.ReadFile(filePath)

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatalln(err)
		}
	}(filePath)
	check(e)
	assert.EqualValues(t, string(data), "CHECKIN,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.TimeCome.Format(DATEFORMATWITHTIME)+","+cred.TimeGone.Format(DATEFORMATWITHTIME)+";\n")
}

func TestLogInToJournal(t *testing.T) {
	var cred = Credentials{
		Address:  "Address",
		Name:     "Name",
		Location: "Location",
		TimeCome: time.Now(),
		TimeGone: time.Now(),
	}

	modifyFlagsForTestCase(true, true)

	LogInToJournal(&cred)

	filePath := returnFilepath()
	data, e := os.ReadFile(filePath)

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatalln(err)
		}
	}(filePath)
	check(e)
	assert.EqualValues(t, string(data), "CHECKIN,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.TimeCome.Format(DATEFORMATWITHTIME)+","+cred.TimeGone.Format(DATEFORMATWITHTIME)+";\n")
}

func TestLogOutToJournal(t *testing.T) {
	var cred = Credentials{
		Address:  "Address",
		Name:     "Name",
		Location: "Location",
		TimeCome: time.Now(),
		TimeGone: time.Now(),
	}
	modifyFlagsForTestCase(true, true)

	LogOutToJournal(&cred)

	filePath := returnFilepath()
	data, e := os.ReadFile(filePath)

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatalln(err)
		}
	}(filePath)
	check(e)
	assert.EqualValues(t, string(data), "CHECKOUT,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.TimeCome.Format(DATEFORMATWITHTIME)+","+cred.TimeGone.Format(DATEFORMATWITHTIME)+";\n")
}

func TestReturnCreditsToString(t *testing.T) {
	var cred = Credentials{
		Address:  "Address",
		Name:     "Name",
		Location: "Location",
		TimeCome: time.Now(),
		TimeGone: time.Now(),
	}
	assert.EqualValues(t, buildCredits(&cred), "CHECKOUT,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.TimeCome.Format(DATEFORMATWITHTIME)+","+cred.TimeGone.Format(DATEFORMATWITHTIME)+";\n")
}

func TestLogTestExample(t *testing.T) {
	modifyFlagsForTestCase(true, false)
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

func modifyFlagsForTestCase(filePath bool, fileName bool) {
	if filePath {
		*config.LogPath = "../../logs"
	}
	if fileName {
		*LogFilename = "testcase"
	}
	if filePath || fileName {
		flag.Parse()
	}
}
