package main

import (
	"DHBW_Golang_Project/internal/config"
	"DHBW_Golang_Project/internal/location"
	"encoding/csv"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
	Erstellt von: 	4775194
	Created by:		4775194

	also: 9514094, 8864957
*/

func TestReadDataFromFile(t *testing.T) {
	/*
		Testfunction validates that the imported content will
		be interpreted correctly
	*/
	config.ConfigureAnalysisTool()
	*config.LogPath = testlogPath

	defer func(path string) {
		err := os.RemoveAll(path)
		checkErrorForTest(err)
	}(*config.LogPath)

	in := "value1-x-y-z;\nvalue2.,!?;\nvalue3\t;\n"
	expected := strings.Split(in, "\n")
	//filePath := filepath.Join("../../logs/temporaryForTest.txt")
	if _, err := os.Stat(*config.LogPath); os.IsNotExist(err) {
		os.MkdirAll(*config.LogPath, 0755)
	}
	filePath := buildFileLogPath(*config.Date)
	f, _ := os.Create(filePath)
	_, e := f.WriteString(in)
	checkErrorForTest(e)

	err := f.Close()
	checkErrorForTest(err)

	out := *readDataFromFile(filePath)
	for i := 0; i < len(out)-1; i++ {
		assert.EqualValues(t, expected[i], out[i])
	}
}

func TestExportToCSVFile(t *testing.T) {
	/*
		Testfunction validates that the exported content will
		be interpreted and formatted correctly
	*/
	config.ConfigureAnalysisTool()
	*config.LogPath = testExportPath

	defer cleanupTestLogs()

	results := []string{
		"location1", "location2", "location3", "location4", "location5",
	}

	*config.Operation = string(VISITOR)
	*config.Query = "TestSelector"
	*config.Testcase = true

	// Tests use path relative from own path
	filePath := buildFileCSVPath()
	csvHeader := createCSVHeader()
	writeSessionsToCSV(&results, csvHeader, filePath)
	f, err := os.Open(filePath)
	checkErrorForTest(err)

	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1
	content, err := csvReader.ReadAll()
	checkErrorForTest(err)

	for j := 0; j < len(content[0]); j++ {
		assert.EqualValues(t, results[j], content[1][j])
	}

	err = f.Close()
	checkErrorForTest(err)
}

func TestBuildFileLogPath(t *testing.T) {
	// Testfunction that validates, that the logpath will be created correctly
	config.ConfigureAnalysisTool()
	*config.LogPath = "./logs"
	in := *config.Date
	out := buildFileLogPath(in)
	expected := "logs/logs-" + in + ".txt"
	assert.EqualValues(t, expected, out)
}

func TestBuildFileCSVPath(t *testing.T) {
	// Testfunction that validates, that the exportpath will be created correctly
	config.ConfigureAnalysisTool()
	*config.LogPath = "./logs"
	*config.Operation = "operation"
	*config.Query = "selector"
	out := buildFileCSVPath()
	expected := "logs/export-operation_selector.csv"
	assert.EqualValues(t, expected, out)
}

func TestCreateCSVHeader(t *testing.T) {
	/*
		Testfunction that validates, that the header/caption of the csv export
		will have to fitting content and format
	 */
	config.ConfigureAnalysisTool()
	*config.Query = "Selector"

	*config.Operation = string(LOCATION)
	out := createCSVHeader()
	assert.EqualValues(t, 1, len(*out))
	assert.EqualValues(t, "Results for the query: "+string(LOCATION)+" = "+*config.Query, (*out)[0])

	*config.Operation = string(VISITOR)
	out = createCSVHeader()
	assert.EqualValues(t, 1, len(*out))
	assert.EqualValues(t, "Results for the query: "+string(VISITOR)+" = "+*config.Query, (*out)[0])

	*config.Operation = string(CONTACT)
	out = createCSVHeader()
	assert.EqualValues(t, 1, len(*out))
	assert.EqualValues(t, string(CONTACT)+"s for the user: "+*config.Query, (*out)[0])
}

func TestToSlice(t *testing.T) {
	name := "Dummyname"
	location := location.Location("Dummylocation")
	duration := 2 * time.Hour

	contact := contact{
		session: session{
			Name:     name,
			Location: location,
		},
		duration: duration,
	}

	out := contact.toSlice()
	assert.EqualValues(t, 3, len(*out))
	assert.EqualValues(t, name, (*out)[0])
	assert.EqualValues(t, location, (*out)[1])
	assert.EqualValues(t, duration.String(), (*out)[2])
}

func TestTrimStringBasedOnOS(t *testing.T) {
	/*
		For the interpretation of files it's important to split/trim
		accordingly/fitting for the used OS
		The testfunction simulates the difference of the os
	 */
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
