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
	/*
		Testfunction that checks, whether the credentials will be exported
		correctly for the case "Checkin"
	*/
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
	/*
		Testfunction that checks, whether the credentials will be exported
		correctly for the case "Checkout"
	*/
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
	/*
		Testfunction that checks, whether the credentials will be transformed correctly
		into the string, which would be exported afterwards.
	*/
	var cred = Credentials{
		Address:   "Address",
		Name:      "Name",
		Location:  "Location",
		Timestamp: time.Now(),
	}
	assert.EqualValues(t, buildCreditString(&cred), "CHECKOUT,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.Timestamp.Format(DATEFORMATWITHTIME)+";\n")
}

func cleanupTestLogs() {
	err := os.RemoveAll(testLogPath)
	check(err)
}
