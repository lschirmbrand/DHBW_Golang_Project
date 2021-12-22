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

type LocationStore struct {
	Locations []Location
}

func NewLocationStore(filepath string) *LocationStore {
	locs, _ := readLocations(filepath)
	return &LocationStore{
		Locations: locs,
	}
}

func readLocations(filepath string) ([]Location, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return []Location{}, err
	}

	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)

	if err != nil {
		return []Location{}, err
	}

	var xmlLocs xmlLocations

	xml.Unmarshal(byteValue, &xmlLocs)

	return xmlLocs.Locations, nil
}

func (st LocationStore) Validate(expLocations Location) bool {
	for _, actLocation := range st.Locations {
		if actLocation == expLocations {
			return true
		}
	}
	return false
}