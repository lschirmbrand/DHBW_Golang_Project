package location

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadLocationsFromReader(t *testing.T) {
	st := NewLocationStore("test_assets/locations_test.xml")

	assert.Equal(t, []Location{
		Location("Test1"),
		Location("Test2"),
		Location("Test3"),
	}, st.Locations)

}

func TestValide(t *testing.T) {

	st := NewLocationStore("test_assets/locations_test.xml")

	assert.True(t, st.Validate("Test1"))
	assert.False(t, st.Validate("Test4"))
}
