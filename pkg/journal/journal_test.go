package journal

import (
	"DHBW_Golang_Project/pkg/config"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testLogPath string = "../../test-logs"
)

func TestLogToJournal(t *testing.T) {

	configure()
	defer cleanupTestLogs()

	var cred = Credentials{
		Login:    true,
		Address:  "Address",
		Name:     "Name",
		Location: "Location",
		Timestamp: time.Now(),
	}

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
	assert.EqualValues(t, string(data), "CHECKIN,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.Timestamp.Format(DATEFORMATWITHTIME)+";\n")
}

func TestLogInToJournal(t *testing.T) {

	configure()
	defer cleanupTestLogs()

	var cred = Credentials{
		Address:  "Address",
		Name:     "Name",
		Location: "Location",
		Timestamp: time.Now(),
	}

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
	assert.EqualValues(t, string(data), "CHECKIN,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.Timestamp.Format(DATEFORMATWITHTIME)+";\n")
}

func TestLogOutToJournal(t *testing.T) {

	configure()
	defer cleanupTestLogs()

	var cred = Credentials{
		Address:  "Address",
		Name:     "Name",
		Location: "Location",
		Timestamp: time.Now(),
	}

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
	assert.EqualValues(t, string(data), "CHECKOUT,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.Timestamp.Format(DATEFORMATWITHTIME)+";\n")
}

func TestReturnCreditsToString(t *testing.T) {
	var cred = Credentials{
		Address:  "Address",
		Name:     "Name",
		Location: "Location",
		Timestamp: time.Now(),
	}
	assert.EqualValues(t, buildCredits(&cred), "CHECKOUT,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.Timestamp.Format(DATEFORMATWITHTIME)+";\n")
}

func TestLogTestExample(t *testing.T) {

	configure()
	defer cleanupTestLogs()

	var cred = Credentials{
		Login:    true,
		Name:     "Name",
		Address:  "Address",
		Location: "Location",
		Timestamp: time.Now(),
	}
	for i := 0; i < 500; i++ {
		LogInToJournal(&cred)
	}
}

func configure() {

	config.LogPath = &testLogPath
}

func cleanupTestLogs() {
	os.RemoveAll(testLogPath)
}
