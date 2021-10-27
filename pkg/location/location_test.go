package location

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadLocationsFromReader(t *testing.T) {
	locations, err := ReadLocations("test_assets/locations_test.xml")

	assert.NoError(t, err)
	assert.Equal(t, []Location{
		Location("Test1"),
		Location("Test2"),
		Location("Test3"),
	}, locations)

}
