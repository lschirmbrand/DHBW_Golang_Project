package journal

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testLogPath = "./test-logs"
)

func TestLogIn(t *testing.T) {

	defer cleanupTestLogs()

	var cred = Credentials{
		Address:   "Address",
		Name:      "Name",
		Location:  "Location",
		Timestamp: time.Now(),
	}

	jour := NewLogFileJournal(testLogPath)

	jour.LogIn(&cred)

	data, e := os.ReadFile(jour.logFile)

	check(e)
	assert.EqualValues(t, string(data), "CHECKIN,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.Timestamp.Format(DATEFORMATWITHTIME)+";\n")
}

func TestLogOut(t *testing.T) {

	defer cleanupTestLogs()

	var cred = Credentials{
		Address:   "Address",
		Name:      "Name",
		Location:  "Location",
		Timestamp: time.Now(),
	}

	jour := NewLogFileJournal(testLogPath)

	jour.LogOut(&cred)

	data, e := os.ReadFile(jour.logFile)

	check(e)
	assert.EqualValues(t, string(data), "CHECKOUT,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.Timestamp.Format(DATEFORMATWITHTIME)+";\n")
}

func TestReturnCreditsToString(t *testing.T) {
	var cred = Credentials{
		Address:   "Address",
		Name:      "Name",
		Location:  "Location",
		Timestamp: time.Now(),
	}
	assert.EqualValues(t, buildCredits(&cred), "CHECKOUT,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.Timestamp.Format(DATEFORMATWITHTIME)+";\n")
}

func TestLogTestExample(t *testing.T) {

	defer cleanupTestLogs()

	jour := NewLogFileJournal(testLogPath)

	var cred = Credentials{
		Checkin:   true,
		Name:      "Name",
		Address:   "Address",
		Location:  "Location",
		Timestamp: time.Now(),
	}
	for i := 0; i < 20; i++ {
		jour.LogIn(&cred)
		jour.LogOut(&cred)
	}
}

func cleanupTestLogs() {
	err := os.RemoveAll(testLogPath)
	check(err)
}
