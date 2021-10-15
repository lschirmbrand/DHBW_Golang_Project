package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Analyse Tool")
	fmt.Println("-------------")
	fmt.Println("Enter Date in format XX-YY-ZZZZ: ")

	for {
		text, _ := reader.ReadString('\n')
		ok, err := validateDateInput(text)
		check(err)
		if ok {
			break
		}
		fmt.Println("Format wrong or pointless. Retry.")
	}

	fmt.Println("Would you like to analyse:")
	fmt.Println("Locations for a person \t[1]")
	fmt.Println("Visitor for a location \t[2]")

	for {
		text, _ := reader.ReadString('\n')
		fmt.Println(text)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func validateDateInput(date string) (bool, error) {
	return regexp.Match("^(0[1-9]|[12][0-9]|3[01])[-](0[1-9]|1[012])[-](19|20)", []byte(date))
}
