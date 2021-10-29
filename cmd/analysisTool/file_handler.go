package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func readDataFromFile(filePath string) *[]string {
	text, err := ioutil.ReadFile(filePath)
	check(err)
	out := strings.Split(string(text), "\n")
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	return &out
}

func exportToCSVFile(results []string, selector string, operation string) {
	filePath := PATHTOCSV + operation+"_"+selector+".csv"
	f, _ := os.Create(filePath)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(f)

	csvLineData := make([]string, len(results) + 1)
	csvLineData[0] = "Results for: " + selector + "\n"
	for i := 1; i< len(results)+ 1; i++ {
		csvLineData[i] = results[i-1]
	}

	w := csv.NewWriter(f)
	e := w.Write(csvLineData)
	w.Flush()
	check(e)

	promptFormatter(1)
	fmt.Println("The query-result was exported to: " + filePath)
}

func buildFilePath(date string) string {
	return PATHTOLOGS + date + ".txt"
}
