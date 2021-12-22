package checkinweb

import (
	"DHBW_Golang_Project/internal/person"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSavePersonToCookies(t *testing.T) {

	store := NewCookieStore(60)

	recorder := httptest.NewRecorder()

	p := person.P{
		Firstname: "Max",
		Lastname:  "Mustermann",
		Street:    "Teststr. 5",
		PLZ:       "12345",
		City:      "Musterstadt",
	}

	store.SaveToCookies(recorder, &p)

	cookies := recorder.Result().Cookies()

	firstNameCookie := cookies[0]
	assert.Equal(t, p.Firstname, decodeFromBase64(firstNameCookie.Value))

	lastNameCookie := cookies[1]
	assert.Equal(t, p.Lastname, decodeFromBase64(lastNameCookie.Value))

	streetCookie := cookies[2]
	assert.Equal(t, p.Street, decodeFromBase64(streetCookie.Value))

	plzCookie := cookies[3]
	assert.Equal(t, p.PLZ, decodeFromBase64(plzCookie.Value))

	cityCookie := cookies[4]
	assert.Equal(t, p.City, decodeFromBase64(cityCookie.Value))

}

func TestReadPersonFromCookies(t *testing.T) {

	req := httptest.NewRequest("GET", "http://localhost", nil)

	p := person.P{
		Firstname: "Max",
		Lastname:  "Mustermann",
		Street:    "Teststr. 5",
		PLZ:       "12345",
		City:      "Musterstadt",
	}

	store := NewCookieStore(60)

	firstNameCookie := http.Cookie{
		Name:  store.FirstNameKey,
		Value: encodeToBase64(p.Firstname),
	}
	lastNameCookie := http.Cookie{
		Name:  store.LastNameKey,
		Value: encodeToBase64(p.Lastname),
	}
	streetCookie := http.Cookie{
		Name:  store.StreetKey,
		Value: encodeToBase64(p.Street),
	}
	plzCookie := http.Cookie{
		Name:  store.PlzKey,
		Value: encodeToBase64(p.PLZ),
	}
	cityCookie := http.Cookie{
		Name:  store.CityKey,
		Value: encodeToBase64(p.City),
	}

	req.AddCookie(&firstNameCookie)
	req.AddCookie(&lastNameCookie)
	req.AddCookie(&streetCookie)
	req.AddCookie(&plzCookie)
	req.AddCookie(&cityCookie)

	p1 := store.ReadFromCookies(req)

	assert.Equal(t, p, *p1)
}
