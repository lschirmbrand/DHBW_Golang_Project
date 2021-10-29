package location

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadLocationsFromReader(t *testing.T) {
	err := ReadLocations("test_assets/locations_test.xml")

	assert.NoError(t, err)
	assert.Equal(t, []Location{
		Location("Test1"),
		Location("Test2"),
		Location("Test3"),
	}, Locations)

}

func TestValide(t *testing.T) {
	// overwrite path to location file
	Locations = []Location{"TestLocation"}

	assert.True(t, Validate("TestLocation"))
	assert.False(t, Validate("Test"))
}
