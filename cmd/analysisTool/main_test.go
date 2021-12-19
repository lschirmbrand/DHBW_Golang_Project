package main

import (
	"DHBW_Golang_Project/pkg/config"
	"DHBW_Golang_Project/pkg/location"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContentToArray(t *testing.T) {
	var content = strings.Split("CHECKIN,name1,address1,location1,20-10-2021 09:44:25;\n"+
		"CHECKOUT,name1,address1,location1,20-10-2021 09:44:41;", "\n")
	sessions := *credentialsToSession(contentToCredits(&content))
	// assert.EqualValues(t, sessions[0].Login, true)
	assert.EqualValues(t, sessions[0].Name, "name1")
	assert.EqualValues(t, sessions[0].Address, "address1")
	assert.EqualValues(t, sessions[0].Location, "location1")
	assert.EqualValues(t, sessions[0].TimeCome.Format(config.DATEFORMATWITHTIME), "20-10-2021 09:44:25")
	assert.EqualValues(t, sessions[0].TimeGone.Format(config.DATEFORMATWITHTIME), "20-10-2021 09:44:41")
	// assert.EqualValues(t, contentArray[1].Login, true)
	// assert.EqualValues(t, contentArray[1].Name, "name2")
	// assert.EqualValues(t, contentArray[1].Address, "address2")
	// assert.EqualValues(t, contentArray[1].Location, "location2")
	// assert.EqualValues(t, contentArray[1].Timestamp.Format(DATEFORMATWITHTIME), "20-10-2021 09:44:41")
}

func BenchmarkPerformanceOfData(b *testing.B) {
	fileContent := "LOGIN,name,address,location,20-10-2021 09:44:25,20-10-2021 09:44:25;\nLOGIN,name,address,location,20-10-2021 09:44:41,20-10-2021 09:44:41;\nLOGIN,name,address,location,20-10-2021 10:07:13,20-10-2021 10:07:13;\nLOGIN,name,address,location,20-10-2021 10:07:18,20-10-2021 10:07:18;\nLOGIN,name,address,location,20-10-2021 10:07:28,20-10-2021 10:07:28;\nLOGIN,name,address,location,20-10-2021 10:07:33,20-10-2021 10:07:33;\nLOGIN,name,address,location,20-10-2021 10:07:33,20-10-2021 10:07:33;"
	for n := 0; n < b.N; n++ {
		content := strings.Split(fileContent, "\n")
		contentToCredits(&content)
	}
}

func TestStartAnalyticalToolDialog(t *testing.T) {
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

/*func TestAnalyseLocationsByVisitor(t *testing.T) {
	creds := make([]journal.Credentials, 0)
	visitor := "Gustav Gans"
	res := analyseLocationsByVisitor(visitor, &creds)
	assert.EqualValues(t, 0, len(*res))

	var searchedCred = journal.Credentials{
		Name:     "Gustav Gans",
		Location: "Entenhausen",
	}

	creds = append(creds, searchedCred)
	res = analyseLocationsByVisitor(visitor, &creds)
	assert.EqualValues(t, 1, len(*res))
	assert.EqualValues(t, "Entenhausen", (*res)[0])

	var notSearchedCred = journal.Credentials{
		Name:     "Donald Duck",
		Location: "Entenhausen",
	}

	creds = append(creds, notSearchedCred)
	res = analyseLocationsByVisitor(visitor, &creds)
	assert.EqualValues(t, 1, len(*res))
	assert.EqualValues(t, "Entenhausen", (*res)[0])
}*/

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

/*
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
 */

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

	var notSearchedCred = session {
		Name:     "Bambis Mutter",
		Location: "Friedhof",
	}

	creds = append(creds, notSearchedCred)
	res = analyseVisitorsByLocation(location, &creds)
	assert.EqualValues(t, 1, len(*res))
	assert.EqualValues(t, "Gustav Gans", (*res)[0])
}

func TestIsOverlapping(t *testing.T) {
	startA := time.Date(2021, time.January, 20, 20, 0, 0, 0, time.Local)
	endA := time.Date(2021, time.January, 20, 23, 0, 0, 0, time.Local)
	startB := time.Date(2021, time.January, 20, 19, 0, 0, 0, time.Local)
	endB := time.Date(2021, time.January, 20, 22, 30, 0, 0, time.Local)
	nameA := "Name"
	nameB := "AuchName"
	locationA := "Location"
	locationB := "Location"

	var sessionA = session{
		Name:     nameA,
		Address:  "",
		Location: location.Location(locationA),
		TimeCome: startA,
		TimeGone: endA,
	}

	var sessionB = session{
		Name:     nameB,
		Address:  "",
		Location: location.Location(locationB),
		TimeCome: startB,
		TimeGone: endB,
	}

	assert.True(t, isOverlapping(&sessionA, &sessionB))
	// assert.False(t, isOverlapping(startB, endB, startA, endA))
}