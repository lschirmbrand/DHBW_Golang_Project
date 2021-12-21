package checkinweb

import (
	"DHBW_Golang_Project/pkg/config"
	"DHBW_Golang_Project/pkg/location"
	"DHBW_Golang_Project/pkg/person"
	"DHBW_Golang_Project/pkg/token"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createServerValidationWrapper(validator token.Validator) *httptest.Server {
	return httptest.NewServer(
		tokenValidationWrapper(
			validator,
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Test")
			}))
}

func TestTokenValidationWrapperValid(t *testing.T) {

	ts := createServerValidationWrapper(func(t token.Token, l location.Location) bool { return true })

	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, "Test\n", string(body))

}

func TestTokenValidationWrapperNotValid(t *testing.T) {

	ts := createServerValidationWrapper(func(t token.Token, l location.Location) bool { return false })

	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, "Bad Request unvalid token\n", string(body))

}

func TestCheckinHandler(t *testing.T) {

	parseTemplates("test_assets/templates")

	req, err := http.NewRequest("GET", "http://localhost", nil)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), locationContextKey, location.Location("TestLocation"))
	ctx = context.WithValue(ctx, tokenContextKey, token.Token("TestToken"))

	recorder := httptest.NewRecorder()

	checkInHandler(recorder, req.WithContext(ctx))
	resp := recorder.Result()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCheckedInHandler(t *testing.T) {
	parseTemplates("test_assets/templates")
	config.ConfigureWeb()

	location.ReadLocations("test_assets/locations.xml")

	reader := strings.NewReader("firstName=Max&lastName=Mustermann&street=Musterstr.+12&plz=12345&city=Musterstadt&location=TestLocation")
	req, err := http.NewRequest("POST", "http://localhost", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	checkedInHandler(recorder, req)
	resp := recorder.Result()

	// http status should be ok
	assert.Equal(t, 200, resp.StatusCode)

	// body should contain name and location
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "Max")
	assert.Contains(t, string(body), "Mustermann")
	assert.Contains(t, string(body), "TestLocation")
}

func TestCheckedOutHandler(t *testing.T) {
	parseTemplates("test_assets/templates")
	config.ConfigureWeb()

	location.ReadLocations("test_assets/locations.xml")

	reader := strings.NewReader("firstName=Max&lastName=Mustermann&location=TestLocation")
	req, err := http.NewRequest("POST", "http://localhost", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	checkedOutHandler(recorder, req)
	resp := recorder.Result()

	// http status should be ok
	assert.Equal(t, 200, resp.StatusCode)

	// body should contain name and location
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "Max")
	assert.Contains(t, string(body), "Mustermann")
	assert.Contains(t, string(body), "TestLocation")
}

func TestValidateFormInput(t *testing.T) {
	assert.True(t, validateFormInput(person.P{
		Firstname: "Max",
		Lastname:  "Mustermann",
		Street:    "Musterstraße 12",
		PLZ:       "12345",
		City:      "Musterstadt",
	}))
	assert.False(t, validateFormInput(person.P{
		Firstname: "",
		Lastname:  "Mustermann",
		Street:    "Musterstraße 12",
		PLZ:       "12345",
		City:      "Musterstadt",
	}))
	assert.False(t, validateFormInput(person.P{
		Firstname: "Max",
		Lastname:  "",
		Street:    "Musterstraße 12",
		PLZ:       "12345",
		City:      "Musterstadt",
	}))
	assert.False(t, validateFormInput(person.P{
		Firstname: "Max",
		Lastname:  "Mustermann",
		Street:    "Musterstraße 12.",
		PLZ:       "12345",
		City:      "Musterstadt",
	}))
	assert.False(t, validateFormInput(person.P{
		Firstname: "Max",
		Lastname:  "Mustermann",
		Street:    "Musterstraße 12",
		PLZ:       "123456",
		City:      "Musterstadt",
	}))
	assert.False(t, validateFormInput(person.P{
		Firstname: "Max",
		Lastname:  "Mustermann",
		Street:    "Musterstraße 12",
		PLZ:       "12345",
		City:      "",
	}))
}
