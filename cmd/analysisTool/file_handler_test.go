package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestExportToCSVFile(t *testing.T){

}

func TestBuildFilePath(t *testing.T){

}