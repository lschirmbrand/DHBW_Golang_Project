package main

import (
	"DHBW_Golang_Project/pkg/config"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testExportPath = "./test-export"
)

func TestAssertQueryExport(t *testing.T) {
	res := make([]string, 0)
	assert.False(t, assertQueryExport(&res))
}

func TestTrimStringBasedOnOS(t *testing.T) {
	if runtime.GOOS == "windows" {
		res := trimStringBasedOnOS("teststring\r\n", true)
		assert.EqualValues(t, res, "teststring")
	} else {
		res := trimStringBasedOnOS("teststring\n", true)
		assert.EqualValues(t, res, "teststring")
	}
	res := trimStringBasedOnOS("\nteststring", false)
	assert.EqualValues(t, res, "teststring")
}

func checkErrorForTest(err error) {
	if err != nil {
		log.Fatalln(err)
	} else {
		return
	}
}

func TestExportContacts(t *testing.T){

	defer cleanupTestLogs()

	config.ConfigureAnalysisTool()
	*config.LogPath = testExportPath
	*config.Testcase = true
	*config.Query = "Selector"
	*config.Operation = string(CONTACT)

	contactA := contact{
		session: sessionB,
		duration: 60*time.Hour,
	}

	contacts := make([]contact, 1)
	contacts[0] = contactA

	exportContacts(&contacts)
	text, err := ioutil.ReadFile(buildFileCSVPath())
	check(err)

	assert.EqualValues(t, "Contacts for the user: Selector\nNameB,Location,60h0m0s\n", string(text))
}

func cleanupTestLogs() {
	err := os.RemoveAll(testExportPath)
	check(err)
}
