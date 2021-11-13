package checkinweb

import (
	"DHBW_Golang_Project/pkg/config"
	"DHBW_Golang_Project/pkg/journal"
	"DHBW_Golang_Project/pkg/location"
	"DHBW_Golang_Project/pkg/token"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"regexp"
	"time"
)

type CheckInPageData struct {
	Person
	Location location.Location
}

type CheckedInPageData struct {
	Person
	Location location.Location
	Time     string
}

type CheckedoutPageData struct {
	Person
	Location location.Location
}

type Person struct {
	Firstname string
	Lastname  string
	Street    string
	PLZ       string
	City      string
}

const (
	locationKey  = "location"
	firstNameKey = "firstName"
	lastNameKey  = "lastName"
	streetKey    = "street"
	plzKey       = "plz"
	cityKey      = "city"
)

var (
	checkInTemplate    *template.Template
	checkedInTemplate  *template.Template
	checkedOutTemplate *template.Template
)

func Mux() http.Handler {
	mux := http.NewServeMux()

	parseTemplates(*config.TemplatePath)

	mux.HandleFunc("/checkin", tokenValidationWrapper(token.Validate, checkInHandler))
	mux.HandleFunc("/checkedin", checkedInHandler)
	mux.HandleFunc("/checkedout", checkedOutHandler)

	// fs := http.FileServer(http.Dir("web/static"))
	// mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}

func parseTemplates(templateDir string) {
	checkInTemplate = template.Must(template.ParseFiles(path.Join(templateDir, "checkin.html")))
	checkedInTemplate = template.Must(template.ParseFiles(path.Join(templateDir, "checkedin.html")))
	checkedOutTemplate = template.Must(template.ParseFiles(path.Join(templateDir, "checkedOut.html")))
}

// handler function for /checkin route
func checkInHandler(rw http.ResponseWriter, r *http.Request) {

	// only GET allowed
	if r.Method != http.MethodGet {
		http.Error(rw,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
	}

	// read location from request context
	loc := r.Context().Value(locationKey).(location.Location)

	// read saved person from cookies
	pers := readPersonFromCookies(r)

	data := CheckInPageData{
		Person:   *pers,
		Location: location.Location(loc),
	}

	checkInTemplate.Execute(rw, data)
}

func checkedInHandler(rw http.ResponseWriter, r *http.Request) {

	// only POST allowed
	if r.Method != http.MethodPost {
		http.Error(rw,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
	}

	// read Person and location from Post Form
	p := Person{
		Firstname: r.PostFormValue(firstNameKey),
		Lastname:  r.PostFormValue(lastNameKey),
		Street:    r.PostFormValue(streetKey),
		PLZ:       r.PostFormValue(plzKey),
		City:      r.PostFormValue(cityKey),
	}

	loc := location.Location(r.PostFormValue(locationKey))

	if !location.Validate(loc) {
		http.Error(rw,
			http.StatusText(http.StatusBadRequest)+"not a valid location",
			http.StatusBadRequest)
	}

	data := CheckedInPageData{
		Person:   p,
		Location: location.Location(loc),
		Time:     time.Now().Format(time.RFC3339),
	}

	savePersonToCookies(rw, &p)

	address := fmt.Sprintf("%v, %v %v", p.Street, p.PLZ, p.City)
	name := fmt.Sprintf("%v %v", p.Firstname, p.Lastname)

	journal.LogInToJournal(&journal.Credentials{
		Login:    true,
		Name:     name,
		Address:  address,
		Location: data.Location,
		TimeCome: time.Now(),
	})

	checkedInTemplate.Execute(rw, data)

}

func checkedOutHandler(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	p := Person{
		Firstname: r.PostFormValue(firstNameKey),
		Lastname:  r.PostFormValue(lastNameKey),
		Street:    r.PostFormValue(streetKey),
		PLZ:       r.PostFormValue(plzKey),
		City:      r.PostFormValue(cityKey),
	}

	data := CheckedoutPageData{
		Person:   p,
		Location: location.Location(r.PostFormValue(locationKey)),
	}

	address := fmt.Sprintf("%v, %v %v", p.Street, p.PLZ, p.City)
	name := fmt.Sprintf("%v %v", p.Firstname, p.Lastname)

	journal.LogOutToJournal(&journal.Credentials{
		Login:    false,
		Name:     name,
		Address:  address,
		Location: data.Location,
		TimeGone: time.Now(),
	})

	checkedOutTemplate.Execute(rw, data)
}

func savePersonToCookies(rw http.ResponseWriter, p *Person) {

	lifetime := time.Hour * time.Duration(*config.CookieLifetime)

	firstNameCookie := http.Cookie{
		Name:    firstNameKey,
		Value:   encodeToBase64(p.Firstname),
		Expires: time.Now().Add(lifetime),
	}
	lastNameCookie := http.Cookie{
		Name:    lastNameKey,
		Value:   encodeToBase64(p.Lastname),
		Expires: time.Now().Add(lifetime),
	}
	streetCookie := http.Cookie{
		Name:    streetKey,
		Value:   encodeToBase64(p.Street),
		Expires: time.Now().Add(lifetime),
	}
	plzCookie := http.Cookie{
		Name:    plzKey,
		Value:   encodeToBase64(p.PLZ),
		Expires: time.Now().Add(lifetime),
	}
	cityCookie := http.Cookie{
		Name:    cityKey,
		Value:   encodeToBase64(p.City),
		Expires: time.Now().Add(lifetime),
	}

	http.SetCookie(rw, &firstNameCookie)
	http.SetCookie(rw, &lastNameCookie)
	http.SetCookie(rw, &streetCookie)
	http.SetCookie(rw, &plzCookie)
	http.SetCookie(rw, &cityCookie)
}

func encodeToBase64(str string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(str))
}

func readPersonFromCookies(r *http.Request) *Person {
	p := Person{
		Firstname: "",
		Lastname:  "",
		Street:    "",
		PLZ:       "",
		City:      "",
	}

	firstName, err := r.Cookie(firstNameKey)
	if err == nil {
		p.Firstname = decodeFromBase64(firstName.Value)
	}

	lastName, err := r.Cookie(firstNameKey)
	if err == nil {
		p.Firstname = decodeFromBase64(lastName.Value)
	}

	street, err := r.Cookie(streetKey)
	if err == nil {
		p.Street = decodeFromBase64(street.Value)
	}

	plz, err := r.Cookie(plzKey)
	if err == nil {
		p.PLZ = decodeFromBase64(plz.Value)
	}

	city, err := r.Cookie(cityKey)
	if err == nil {
		p.City = decodeFromBase64(city.Value)
	}

	return &p
}

func decodeFromBase64(encoded string) string {
	decoded, _ := base64.RawStdEncoding.DecodeString(encoded)

	return string(decoded)
}

func tokenValidationWrapper(validator token.Validator, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := token.Token(r.URL.Query().Get("token"))
		l := location.Location(r.URL.Query().Get("location"))
		if valid := validator(t, l); valid {
			ctx := context.WithValue(r.Context(), locationKey, l)

			handler(w, r.WithContext(ctx))
		} else {
			http.Error(w,
				http.StatusText(http.StatusBadRequest)+"not a valid token",
				http.StatusBadRequest)
		}
	}
}

func validateFormInput(p Person) bool {

	regexp.Match("[0-9]{5}", []byte(p.PLZ))

	return true
}
