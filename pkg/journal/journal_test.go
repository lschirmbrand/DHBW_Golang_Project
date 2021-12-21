package journal

import (
	"DHBW_Golang_Project/pkg/config"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testLogPath = "../../test-logs"
)

func TestLogToJournal(t *testing.T) {

	configure()
	defer cleanupTestLogs()

	var cred = Credentials{
		Checkin:  true,
		Address:  "Address",
		Name:     "Name",
		Location: "Location",
		Timestamp: time.Now(),
	}

	filePath := returnFilepath()
	logToJournal(&cred)
	data, e := os.ReadFile(filePath)
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
		Checkin:  true,
		Name:     "Name",
		Address:  "Address",
		Location: "Location",
		Timestamp: time.Now(),
	}
	for i := 0; i < 20; i++ {
		LogInToJournal(&cred)
		LogOutToJournal(&cred)
	}
}

func configure() {
	config.LogPath = &testLogPath
}

func cleanupTestLogs() {
	//err := os.RemoveAll(testLogPath)
	//check(err)
}
