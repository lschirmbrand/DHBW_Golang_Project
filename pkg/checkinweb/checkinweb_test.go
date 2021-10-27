package checkinweb

import (
	"DHBW_Golang_Project/pkg/location"
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

	assert.Equal(t, "Bad Request\n", string(body))

}

func TestSavePersonToCookies(t *testing.T) {
	recorder := httptest.NewRecorder()

	p := Person{
		Name:   "Max Mustermann",
		Street: "Teststr. 5",
		PLZ:    "12345",
		City:   "Musterstadt",
	}

	savePersonToCookies(recorder, &p)

	cookies := recorder.Result().Cookies()

	nameCookie := cookies[0]
	assert.Equal(t, p.Name, nameCookie.Value)

	streetCookie := cookies[1]
	assert.Equal(t, p.Street, streetCookie.Value)

	plzCookie := cookies[2]
	assert.Equal(t, p.PLZ, plzCookie.Value)

	cityCookie := cookies[3]
	assert.Equal(t, p.City, cityCookie.Value)

}

func TestReadPersonFromCookies(t *testing.T) {

	req := httptest.NewRequest("GET", "http://localhost", nil)

	p := Person{
		Name:   "Max Mustermann",
		Street: "Teststr. 5",
		PLZ:    "12345",
		City:   "Musterstadt",
	}

	nameCookie := http.Cookie{
		Name:  string(nameKey),
		Value: p.Name,
	}
	streetCookie := http.Cookie{
		Name:  string(streetKey),
		Value: p.Street,
	}
	plzCookie := http.Cookie{
		Name:  string(plzKey),
		Value: p.PLZ,
	}
	cityCookie := http.Cookie{
		Name:  string(cityKey),
		Value: p.City,
	}

	req.AddCookie(&nameCookie)
	req.AddCookie(&streetCookie)
	req.AddCookie(&plzCookie)
	req.AddCookie(&cityCookie)

	p1 := readPersonFromCookies(req)

	assert.Equal(t, p, *p1)
}

func TestCheckinHandler(t *testing.T) {

	parseTemplates("test_assets/templates")

	req, err := http.NewRequest("GET", "http://localhost", nil)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), locationKey, location.Location("TestLocation"))

	recorder := httptest.NewRecorder()

	checkInHandler(recorder, req.WithContext(ctx))
	resp := recorder.Result()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCheckedInHandler(t *testing.T) {
	parseTemplates("test_assets/templates")

	reader := strings.NewReader("name=Max+Mustermann&street=Musterstr.+12&plz=12345&city=Musterstadt&location=TestLocation")
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
	assert.Equal(t, "Max Mustermann,TestLocation\n", string(body))

	// cookies should be set
	cookies := resp.Cookies()

	assert.Equal(t, "Max Mustermann", cookies[0].Value)
	assert.Equal(t, "Musterstr. 12", cookies[1].Value)
	assert.Equal(t, "12345", cookies[2].Value)
	assert.Equal(t, "Musterstadt", cookies[3].Value)
}

func TestCheckedOutHandler(t *testing.T) {
	parseTemplates("test_assets/templates")

	reader := strings.NewReader("name=Max+Mustermann&location=TestLocation")
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
	assert.Equal(t, "Max Mustermann,TestLocation\n", string(body))
}