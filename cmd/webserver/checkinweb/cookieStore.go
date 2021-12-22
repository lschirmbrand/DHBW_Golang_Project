package checkinweb

import (
	"DHBW_Golang_Project/internal/person"
	"encoding/base64"
	"net/http"
	"time"
)

const ()

type CookieStore struct {
	cookieLifetime time.Duration
	FirstNameKey   string
	LastNameKey    string
	StreetKey      string
	PlzKey         string
	CityKey        string
}

func NewCookieStore(cookieLifeSeconds int) *CookieStore {
	return &CookieStore{
		cookieLifetime: time.Hour * time.Duration(cookieLifeSeconds),
		FirstNameKey:   "firstName",
		LastNameKey:    "lastName",
		StreetKey:      "street",
		PlzKey:         "plz",
		CityKey:        "city",
	}
}

func (cs CookieStore) SaveToCookies(rw http.ResponseWriter, p *person.P) {

	firstNameCookie := http.Cookie{
		Name:    cs.FirstNameKey,
		Value:   encodeToBase64(p.Firstname),
		Expires: time.Now().Add(cs.cookieLifetime),
	}
	lastNameCookie := http.Cookie{
		Name:    cs.LastNameKey,
		Value:   encodeToBase64(p.Lastname),
		Expires: time.Now().Add(cs.cookieLifetime),
	}
	streetCookie := http.Cookie{
		Name:    cs.StreetKey,
		Value:   encodeToBase64(p.Street),
		Expires: time.Now().Add(cs.cookieLifetime),
	}
	plzCookie := http.Cookie{
		Name:    cs.PlzKey,
		Value:   encodeToBase64(p.PLZ),
		Expires: time.Now().Add(cs.cookieLifetime),
	}
	cityCookie := http.Cookie{
		Name:    cs.CityKey,
		Value:   encodeToBase64(p.City),
		Expires: time.Now().Add(cs.cookieLifetime),
	}

	http.SetCookie(rw, &firstNameCookie)
	http.SetCookie(rw, &lastNameCookie)
	http.SetCookie(rw, &streetCookie)
	http.SetCookie(rw, &plzCookie)
	http.SetCookie(rw, &cityCookie)
}

func (cs CookieStore) ReadFromCookies(r *http.Request) *person.P {
	p := person.P{
		Firstname: "",
		Lastname:  "",
		Street:    "",
		PLZ:       "",
		City:      "",
	}

	firstName, err := r.Cookie(cs.FirstNameKey)
	if err == nil {
		p.Firstname = decodeFromBase64(firstName.Value)
	}

	lastName, err := r.Cookie(cs.LastNameKey)
	if err == nil {
		p.Lastname = decodeFromBase64(lastName.Value)
	}

	street, err := r.Cookie(cs.StreetKey)
	if err == nil {
		p.Street = decodeFromBase64(street.Value)
	}

	plz, err := r.Cookie(cs.PlzKey)
	if err == nil {
		p.PLZ = decodeFromBase64(plz.Value)
	}

	city, err := r.Cookie(cs.CityKey)
	if err == nil {
		p.City = decodeFromBase64(city.Value)
	}

	return &p
}

func encodeToBase64(str string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(str))
}

func decodeFromBase64(encoded string) string {
	decoded, _ := base64.RawStdEncoding.DecodeString(encoded)

	return string(decoded)
}
