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

type Location string

func ReadLocations(filepath string) ([]Location, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	var locations xmlLocations

	xml.Unmarshal(byteValue, &locations)

	return locations.Locations, nil
}
