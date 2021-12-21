package main

import (
	"DHBW_Golang_Project/pkg/config"
	"DHBW_Golang_Project/pkg/location"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadDataFromFile(t *testing.T) {
	config.ConfigureAnalysisTool()
	in := "value1-x-y-z;\nvalue2.,!?;\nvalue3\t;\n"
	expected := strings.Split(in, "\n")
	filePath := filepath.Join("../../logs/temporaryForTest.txt")
	f, _ := os.Create(filePath)
	_, e := f.WriteString(in)
	check(e)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatalln(err)
		}
	}(filePath)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(f)

	out := *readDataFromFile(filePath)
	for i := 0; i < len(out)-1; i++ {
		assert.EqualValues(t, expected[i], out[i])
	}
}

func TestExportToCSVFile(t *testing.T) {
	config.ConfigureAnalysisTool()
	results := []string{
		"location1", "location2", "location3", "location4", "location5",
	}

	*config.Operation = string(VISITOR)
	*config.Query = "TestSelector"

	// Tests use path relative from own path
	filePath := "../../" + buildFileCSVPath()
	csvHeader := createCSVHeader()
	writeSessionsToCSV(&results, csvHeader, filePath)
	f, err := os.Open(filePath)
	checkErrorForTest(err)

	defer func(name string) {
		err := f.Close()
		checkErrorForTest(err)
		err = os.Remove(name)
		checkErrorForTest(err)
	}(filePath)

	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1
	content, err := csvReader.ReadAll()
	checkErrorForTest(err)

	for j := 0; j < len(content[0]); j++ {
		assert.EqualValues(t, results[j], content[1][j])
	}
}

func TestBuildFileLogPath(t *testing.T) {
	config.ConfigureAnalysisTool()
	in := time.Now().Format(config.DATEFORMAT)
	out := buildFileLogPath(in)
	expected := "logs/logs-" + in + ".txt"
	assert.EqualValues(t, expected, out)
}

func TestBuildFileCSVPath(t *testing.T) {
	config.ConfigureAnalysisTool()
	*config.Operation = "operation"
	*config.Query = "selector"
	out := buildFileCSVPath()
	expected := "logs/export-operation_selector.csv"
	assert.EqualValues(t, expected, out)
}

func TestCreateCSVHeader(t *testing.T){
	config.ConfigureAnalysisTool()
	*config.Query = "Selector"

	*config.Operation = string(LOCATION)
	out := createCSVHeader()
	assert.EqualValues(t, 1, len(*out))
	assert.EqualValues(t, "Results for the query: " + string(LOCATION) + " = " + *config.Query, (*out)[0])

	*config.Operation = string(VISITOR)
	out = createCSVHeader()
	assert.EqualValues(t, 1, len(*out))
	assert.EqualValues(t, "Results for the query: " + string(VISITOR) + " = " + *config.Query, (*out)[0])

	*config.Operation = string(CONTACT)
	out = createCSVHeader()
	assert.EqualValues(t, 1, len(*out))
	assert.EqualValues(t, string(CONTACT) + "s for the user: " + *config.Query, (*out)[0])
}

func TestToSlice(t *testing.T){
	name := "Dummyname"
	location := location.Location("Dummylocation")
	duration := 2*time.Hour

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
