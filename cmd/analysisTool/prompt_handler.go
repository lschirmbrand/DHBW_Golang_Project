package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func assertQueryExport(s *[]string) bool {
	qLen := queryLengthHandler(*s)
	if qLen > 0 {
		if exportHandler(qLen) {
			return true
		} else {
			fmt.Println("Results of query wont get exported. \nAborting.")
			return false
		}
	} else {
		fmt.Println("No results were found for the queried selector.")
		return false
	}
}

func exportHandler(length int) bool {
	fmt.Println("The requested query resulted in ", length, " elements.")
	fmt.Println("Do you want to export the query? [y/n]")
	reader := bufio.NewReader(os.Stdin)
	for {
		input, e := reader.ReadString('\n')
		check(e)
		ok, e := validateYesNoInput(input)
		check(e)
		if ok {
			return strings.EqualFold(trimStringBasedOnOS(input, true), "y")
		}
		fmt.Println("Input was incorrect, retry.")
	}
}

func requestedHelp(args *[]string) bool {
	if len(*args) > 0 {
		for i := range *args {
			if strings.EqualFold((*args)[i], "--help") {
				return true
			}
		}
	}
	return false
}

func queryLengthHandler(slice []string) int {
	return len(slice)
}

func trimStringBasedOnOS(text string, isSuffix bool) string {
	isWindows := runtime.GOOS == "windows"
	if isSuffix {
		if isWindows {
			text = strings.TrimSuffix(text, "\x0a\x0d")
			return strings.TrimSuffix(text, "\r\n")
		}
		text = strings.TrimSuffix(text, "\x0d")
		return strings.TrimSuffix(text, "\n")
	} else {
		text = strings.TrimPrefix(text, "\x0d")
		return strings.TrimPrefix(text, "\n")
	}
}
