package location

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type xmlLocations struct {
	XMLName   xml.Name   `xml:"locations"`
	Locations []Location `xml:"location"`
}

var (
	Locations []Location
)

type Location string

func ReadLocations(filepath string) error {
	file, err := os.Open(filepath)

	if err != nil {
		return err
	}

	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)

	if err != nil {
		return err
	}

	var xmlLocs xmlLocations

	xml.Unmarshal(byteValue, &xmlLocs)

	Locations = xmlLocs.Locations
	return nil
}

func Validate(expLocations Location) bool {
	for _, actLocation := range Locations {
		if actLocation == expLocations {
			return true
		}
	}
	return false
}
