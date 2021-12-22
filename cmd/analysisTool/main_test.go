package main

import (
	"DHBW_Golang_Project/internal/config"
	"DHBW_Golang_Project/internal/location"
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

var (
	testlogPath = "./test-log"
)

func TestCredentialsToSession(t *testing.T) {
	/*
		Testfunction that validates, that the imported content was transformed correctly
		into credentials and after that transformed correctly into sessions
	 */
	var content = strings.Split("CHECKIN,name1,address1,location1,20-10-2021 09:44:25;\n"+
		"CHECKOUT,name1,address1,location1,20-10-2021 09:44:41;", "\n")
	sessions := *credentialsToSession(contentToCredits(&content))
	assert.EqualValues(t, sessions[0].Name, "name1")
	assert.EqualValues(t, sessions[0].Address, "address1")
	assert.EqualValues(t, sessions[0].Location, "location1")
	assert.EqualValues(t, sessions[0].TimeCome.Format(config.DATEFORMATWITHTIME), "20-10-2021 09:44:25")
	assert.EqualValues(t, sessions[0].TimeGone.Format(config.DATEFORMATWITHTIME), "20-10-2021 09:44:41")
}

func TestAnalyseLocationsByVisitor(t *testing.T) {
	/*
		Testfunction that validates, that the locations visited by a person will be
		interpreted correctly
	 */
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
	/*
		Testfunction that validates, that the visitors that visited a queried
		location will be interpreted correctly
	*/
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
	// Testfunction, that checks whether an overlap of contacts gets detected
	assert.True(t, isOverlapping(&sessionA, &sessionB))
}

func TestCalculateOverlap(t *testing.T) {
	/*
		Testfunction that checks, whether the duration of the overlap of at
		least two sessions will be calculated correctly
	 */
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
	/*
		Testfunction that checks if an overlap of contacts get detected and
		returns the duration. If so, the duration will also get checked
	 */
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
