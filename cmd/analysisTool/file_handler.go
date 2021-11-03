package main

import (
	"DHBW_Golang_Project/pkg/config"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

func exportToCSVFile(results *[]string, selector string, operation Operation, filePath string) {
	f, e := os.Create(filePath)
	check(e)
	defer f.Close()

	csvHeader := make([]string, 1)
	csvHeader[0] = "Results for: " + selector

	w := csv.NewWriter(f)
	e = w.Write(csvHeader)
	check(e)
	e = w.Write(*(results))
	w.Flush()
	check(e)

	fmt.Println("The query-result was exported to: " + filePath)
}

func buildFileLogPath(date string) string {
	return path.Join(*config.LogPath, "log-"+date+".txt")
}

func buildFileCSVPath(operation Operation, selector string) string {
	return path.Join(*config.LogPath, "export-"+string(operation)+"_"+selector+".csv")
}
