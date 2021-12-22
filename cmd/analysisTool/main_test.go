package main

import (
	"DHBW_Golang_Project/internal/config"
	"DHBW_Golang_Project/internal/location"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	startA    = time.Date(2021, time.January, 20, 20, 0, 0, 0, time.Local)
	endA      = time.Date(2021, time.January, 20, 22, 0, 0, 0, time.Local)
	startB    = time.Date(2021, time.January, 20, 19, 0, 0, 0, time.Local)
	endB      = time.Date(2021, time.January, 20, 21, 0, 0, 0, time.Local)
	nameA     = "NameA"
	nameB     = "NameB"
	locationA = "Location"
	locationB = "Location"
	addressA  = "AddressA"
	addressB  = "AddressB"
)
var sessionA = session{
	Name:     nameA,
	Address:  addressA,
	Location: location.Location(locationA),
	TimeCome: startA,
	TimeGone: endA,
}

var sessionB = session{
	Name:     nameB,
	Address:  addressB,
	Location: location.Location(locationB),
	TimeCome: startB,
	TimeGone: endB,
}

func TestContentToArray(t *testing.T) {
	var content = strings.Split("CHECKIN,name1,address1,location1,20-10-2021 09:44:25;\n"+
		"CHECKOUT,name1,address1,location1,20-10-2021 09:44:41;", "\n")
	sessions := *credentialsToSession(contentToCredits(&content))
	assert.EqualValues(t, sessions[0].Name, "name1")
	assert.EqualValues(t, sessions[0].Address, "address1")
	assert.EqualValues(t, sessions[0].Location, "location1")
	assert.EqualValues(t, sessions[0].TimeCome.Format(config.DATEFORMATWITHTIME), "20-10-2021 09:44:25")
	assert.EqualValues(t, sessions[0].TimeGone.Format(config.DATEFORMATWITHTIME), "20-10-2021 09:44:41")
}

func TestStartAnalyticalToolDialog(t *testing.T) {
	*config.Testcase = true
	filePath := "../../" + buildFileLogPath(time.Now().Format(config.DATEFORMAT))
	_, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	checkErrorForTest(err)

	date := ""
	config.Date = &date
	assert.False(t, startAnalyticalToolDialog())

	date = "../../" + time.Now().Format(config.DATEFORMAT)
	operation := "Visitor"
	config.Operation = &operation
	query := "abcdefghijklmnopqrstuvwxyz"
	config.Query = &query
	assert.False(t, startAnalyticalToolDialog())
}

func TestAnalyseLocationsByVisitor(t *testing.T) {
	creds := make([]session, 0)
	visitor := "Gustav Gans"
	res := analyseLocationsByVisitor(visitor, &creds)
	assert.EqualValues(t, 0, len(*res))

	var searchedCred = session{
		Name:     "Gustav Gans",
		Location: "Entenhausen",
	}

	creds = append(creds, searchedCred)
	res = analyseLocationsByVisitor(visitor, &creds)
	assert.EqualValues(t, 1, len(*res))
	assert.EqualValues(t, "Entenhausen", (*res)[0])

	var notSearchedCred = session{
		Name:     "Donald Duck",
		Location: "Entenhausen",
	}

	creds = append(creds, notSearchedCred)
	res = analyseLocationsByVisitor(visitor, &creds)
	assert.EqualValues(t, 1, len(*res))
	assert.EqualValues(t, "Entenhausen", (*res)[0])
}

func TestAnalyseVisitorsByLocation(t *testing.T) {
	creds := make([]session, 0)
	location := "Entenhausen"
	res := analyseVisitorsByLocation(location, &creds)
	assert.EqualValues(t, 0, len(*res))

	var searchedCred = session{
		Name:     "Gustav Gans",
		Location: "Entenhausen",
	}

	creds = append(creds, searchedCred)
	res = analyseVisitorsByLocation(location, &creds)
	assert.EqualValues(t, 1, len(*res))
	assert.EqualValues(t, "Gustav Gans", (*res)[0])

	var notSearchedCred = session{
		Name:     "Bambis Mutter",
		Location: "Friedhof",
	}

	creds = append(creds, notSearchedCred)
	res = analyseVisitorsByLocation(location, &creds)
	assert.EqualValues(t, 1, len(*res))
	assert.EqualValues(t, "Gustav Gans", (*res)[0])
}

func TestIsOverlapping(t *testing.T) {
	assert.True(t, isOverlapping(&sessionA, &sessionB))
}

func TestCalculateOverlap(t *testing.T) {
	overlap := calculateOverlap(&sessionA, &sessionB)
	assert.EqualValues(t, 1*time.Hour, overlap)

	sessionA.TimeCome = startB
	sessionA.TimeGone = endB
	sessionB.TimeCome = startA
	sessionA.TimeGone = endA

	overlap = calculateOverlap(&sessionA, &sessionB)
	assert.EqualValues(t, 1*time.Hour, overlap)

	sessionA.TimeCome = startB
	sessionA.TimeGone = endA
	sessionB.TimeCome = startA
	sessionA.TimeGone = endB

	overlap = calculateOverlap(&sessionA, &sessionB)
	assert.EqualValues(t, 1*time.Hour, overlap)
}

func TestGetOverlaps(t *testing.T) {
	sessionA.TimeCome = startA
	sessionA.TimeGone = endA
	sessionB.TimeCome = startB
	sessionA.TimeGone = endB

	sessions := sessionsToSlice()
	contacts := getOverlaps(&sessionA, &sessions)
	assert.EqualValues(t, 1, len(*contacts))
	assert.EqualValues(t, 1*time.Hour, (*contacts)[0].duration)
	assert.EqualValues(t, sessionB, (*contacts)[0].session)
}

func TestContactHandler(t *testing.T) {
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
	sessions := make([]session, 2)
	sessions[0] = sessionA
	sessions[1] = sessionB
	return sessions
}
