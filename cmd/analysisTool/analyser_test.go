package main

import (
	"DHBW_Golang_Project/internal/config"
	"DHBW_Golang_Project/internal/location"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

/*
	Erstellt von: 	4775194
	Created by:		4775194

	also: 9514094, 8864957
*/

var (
	testExportPath = "./test-export"
)

func TestContactHandler(t *testing.T) {
	/*
		Testfunction validates, that the output of two overlapping
		sessions will be correctly interpreted as a contact
	 */
	config.ConfigureAnalysisTool()
	*config.Operation = string(CONTACT)
	*config.Query = "NameA"
	sessions := make([]session, 2)
	sessions[0] = sessionA
	sessions[1] = sessionB

	out := contactHandler(&sessions)
	assert.EqualValues(t, 1, len(*out))
	assert.EqualValues(t, sessionB.Name, (*out)[0].session.Name)
	assert.EqualValues(t, sessionB.Location, (*out)[0].session.Location)
	assert.EqualValues(t, sessionB.Address, (*out)[0].session.Address)
	assert.EqualValues(t, sessionB.TimeCome, (*out)[0].session.TimeCome)
	assert.EqualValues(t, sessionB.TimeGone, (*out)[0].session.TimeGone)
	assert.EqualValues(t, 1*time.Hour, (*out)[0].duration)
}

func TestVisitorHandler(t *testing.T) {
	/*
		Testfunction validates that all Locations of a queried
		person will be interpreted and logged correctly
	*/
	config.ConfigureAnalysisTool()
	*config.Operation = string(VISITOR)
	*config.Query = sessionA.Name
	sessions := sessionsToSlice()

	newLocation := "Different Location"
	sessions = append(sessions, session{
		Name:     nameA,
		Address:  addressA,
		Location: location.Location(newLocation),
	})

	out := visitorHandler(&sessions)
	assert.EqualValues(t, 2, len(*out))
	assert.EqualValues(t, sessionA.Location, (*out)[0])
	assert.EqualValues(t, newLocation, (*out)[1])
}

func TestLocationHandler(t *testing.T) {
	/*
		Testfunction validates that all visitors of a queried
		location will be interpreted and logged correctly
	*/
	config.ConfigureAnalysisTool()
	*config.Testcase = true
	*config.Operation = string(VISITOR)
	*config.Query = locationA
	sessions := sessionsToSlice()
	newVisitor := "Different Visitor"
	sessions = append(sessions, session{
		Name:     newVisitor,
		Address:  addressA,
		Location: location.Location(locationA),
	})
	out := locationHandler(&sessions)
	assert.EqualValues(t, 3, len(*out))
	assert.EqualValues(t, sessionA.Name, (*out)[0])
	assert.EqualValues(t, sessionB.Name, (*out)[1])
	assert.EqualValues(t, newVisitor, (*out)[2])
}

func sessionsToSlice() []session {
	// Helperfunction that returns a slice of the available sessions
	sessions := make([]session, 2)
	sessions[0] = sessionA
	sessions[1] = sessionB
	return sessions
}

func TestAssertQueryExport(t *testing.T) {
	res := make([]string, 0)
	assert.False(t, assertQueryExport(&res))
}

func checkErrorForTest(err error) {
	if err != nil {
		log.Fatalln(err)
	} else {
		return
	}
}

func TestExportContacts(t *testing.T){
	/*
		Testfunction validates that the export of the contacts
		for the queried person will be exported correctly
	*/
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

func TestExportLocations(t *testing.T){
	/*
		Testfunction validates that the export of the locations
		for the queried visitors will be exported correctly
	*/
	defer cleanupTestLogs()

	config.ConfigureAnalysisTool()
	*config.LogPath = testExportPath
	*config.Testcase = true
	*config.Query = "NameA"
	*config.Operation = string(VISITOR)
	sessions := sessionsToSlice()
	locations := visitorHandler(&sessions)
	exportLocationsForVisitor(locations)
	text, err := ioutil.ReadFile(buildFileCSVPath())
	check(err)

	assert.EqualValues(t, "Results for the query: Visitor = NameA\nLocation\n", string(text))
}

func TestExportVisitors(t *testing.T){
	/*
		Testfunction validates that the export of the visitors
		for the queried location will be exported correctly
	*/
	defer cleanupTestLogs()

	config.ConfigureAnalysisTool()
	*config.LogPath = testExportPath
	*config.Testcase = true
	*config.Query = "Location"
	*config.Operation = string(LOCATION)
	sessions := sessionsToSlice()
	visitors := locationHandler(&sessions)
	exportVisitorsForLocation(visitors)
	text, err := ioutil.ReadFile(buildFileCSVPath())
	check(err)

	assert.EqualValues(t, "Results for the query: Location = Location\nNameA,NameB\n", string(text))
}

func cleanupTestLogs() {
	err := os.RemoveAll(testExportPath)
	check(err)
}