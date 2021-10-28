package main

import (
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReadDataFromFile(t *testing.T) {
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
	results := []string{
		"location1", "location2", "location3", "location4", "location5",
	}
	selector := "TestSelector"
	operation := Operation("TestOperation")

	// Tests use path relative from own path
	filePath := "../../" + buildFileCSVPath(operation, selector)
	exportToCSVFile(&results, selector, operation, filePath)
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

	assert.EqualValues(t, "Results for: "+selector, content[0][0])

	for j := 0; j < len(content[0]); j++ {
		assert.EqualValues(t, results[j], content[1][j])
	}
}

func TestBuildFileLogPath(t *testing.T) {
	in := time.Now().Format(DATEFORMAT)
	out := buildFileLogPath(in)
	expected := "logs/log-" + in + ".txt"
	assert.EqualValues(t, expected, out)
}

func TestBuildFileCSVPath(t *testing.T) {
	out := buildFileCSVPath("operation", "selector")
	expected := "logs/export-operation_selector.csv"
	assert.EqualValues(t, expected, out)
}