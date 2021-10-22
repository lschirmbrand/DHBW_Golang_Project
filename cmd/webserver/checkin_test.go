package main

import (
	"DHBW_Golang_Project/pkg/token"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createServerValidationWrapper(validator token.Validator) *httptest.Server {
	return httptest.NewServer(
		tokenValidationWrapper(
			validator,
			func(w http.ResponseWriter, r *http.Request) {
				location := r.Context().Value(locationContextKey).(string)
				fmt.Fprintln(w, location)
			}))
}

func TestTokenValidationWrapperValid(t *testing.T) {

	ts := createServerValidationWrapper(func(t token.Token) (bool, string) { return true, "lol" })

	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, "lol\n", string(body))

}

func TestTokenValidationWrapperNotValid(t *testing.T) {

	ts := createServerValidationWrapper(func(t token.Token) (bool, string) { return false, "lol" })

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
		Name:  "name",
		Value: p.Name,
	}
	streetCookie := http.Cookie{
		Name:  "street",
		Value: p.Street,
	}
	plzCookie := http.Cookie{
		Name:  "plz",
		Value: p.PLZ,
	}
	cityCookie := http.Cookie{
		Name:  "city",
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
	parseTemplates("../../web/templates")

	req, err := http.NewRequest("GET", "http://localhost", nil)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), locationContextKey, "TestLocation")

	recorder := httptest.NewRecorder()

	checkInHandler(recorder, req.WithContext(ctx))
	resp := recorder.Result()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCheckOutHandler(t *testing.T) {
	parseTemplates("../../web/templates")

	req, err := http.NewRequest("GET", "http://localhost", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	checkOutHandler(recorder, req)
	resp := recorder.Result()
	assert.Equal(t, 200, resp.StatusCode)
}
