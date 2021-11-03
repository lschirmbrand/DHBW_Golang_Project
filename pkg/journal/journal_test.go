package journal

import (
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogToJournal(t *testing.T) {
	var cred = Credentials{
		Login:     true,
		Address:   "Address",
		Name:      "Name",
		Location:  "Location",
		Timestamp: time.Now(),
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
	assert.EqualValues(t, string(data), "CHECKIN,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.Timestamp.Format(DATEFORMATWITHTIME)+";\n")
	resetFlags()
}

func TestLogInToJournal(t *testing.T) {
	var cred = Credentials{
		Address:   "Address",
		Name:      "Name",
		Location:  "Location",
		Timestamp: time.Now(),
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
	assert.EqualValues(t, string(data), "CHECKIN,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.Timestamp.Format(DATEFORMATWITHTIME)+";\n")
	resetFlags()
}

func TestLogOutToJournal(t *testing.T) {
	var cred = Credentials{
		Address:   "Address",
		Name:      "Name",
		Location:  "Location",
		Timestamp: time.Now(),
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
	assert.EqualValues(t, string(data), "CHECKOUT,"+cred.Name+","+cred.Address+","+string(cred.Location)+","+cred.Timestamp.Format(DATEFORMATWITHTIME)+";\n")
	resetFlags()
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
	modifyFlagsForTestCase(true, false)
	var cred = Credentials{
		Login:     true,
		Name:      "Name",
		Address:   "Address",
		Location:  "Location",
		Timestamp: time.Now(),
	}
	for i := 0; i < 250; i++ {
		cred.Login = true
		LogInToJournal(&cred)
		cred.Login = false
		LogOutToJournal(&cred)
	}
}

func modifyFlagsForTestCase(filePath bool, fileName bool) {
	if filePath {
		*LogPath = "../../" + PATHTOLOGS
	}
	if fileName {
		*LogFilename = "testcase"
	}
	if filePath || fileName {
		flag.Parse()
	}
}

func resetFlags() {
	*LogPath = PATHTOLOGS
	*LogFilename = time.Now().Format(DATEFORMAT)
	flag.Parse()
}
