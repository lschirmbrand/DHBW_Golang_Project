package checkinweb

import (
	"DHBW_Golang_Project/pkg/config"
	"DHBW_Golang_Project/pkg/journal"
	"DHBW_Golang_Project/pkg/location"
	"DHBW_Golang_Project/pkg/person"
	"DHBW_Golang_Project/pkg/token"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"time"
)

type CheckInPageData struct {
	Person       person.P
	Location     string
	Token        string
	InvalidInput bool
}

type CheckedInPageData struct {
	Person   person.P
	Location string
	Time     string
}

type CheckedoutPageData struct {
	Person   person.P
	Location string
}

type contextKey string

const (
	locationContextKey = contextKey("location")
	tokenContextKey    = contextKey("token")
)

const (
	locationKey = "location"
	tokenKey    = "token"
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

	// read location and token from request context
	loc := r.Context().Value(locationContextKey).(location.Location)
	tok := r.Context().Value(tokenContextKey).(token.Token)

	invalid := r.URL.Query().Has("invalid_input")

	// read saved person from cookies
	pers := person.ReadFromCookies(r)

	data := CheckInPageData{
		Person:       *pers,
		Location:     string(loc),
		Token:        string(tok),
		InvalidInput: invalid,
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
	p := person.P{
		Firstname: r.PostFormValue(person.FirstNameKey),
		Lastname:  r.PostFormValue(person.LastNameKey),
		Street:    r.PostFormValue(person.StreetKey),
		PLZ:       r.PostFormValue(person.PlzKey),
		City:      r.PostFormValue(person.CityKey),
	}

	loc := location.Location(r.PostFormValue(locationKey))

	// validate location
	if !location.Validate(loc) {
		http.Error(rw,
			http.StatusText(http.StatusBadRequest)+"not a valid location",
			http.StatusBadRequest)
	}

	person.SaveToCookies(rw, &p)

	// validate Person input

	if !validateFormInput(p) {
		token := r.PostFormValue(tokenKey)
		url := fmt.Sprintf("/checkin?location=%v&token=%v&invalid_input", url.QueryEscape(string(loc)), url.QueryEscape(token))

		http.Redirect(rw, r, url, http.StatusSeeOther)
		return
	}

	data := CheckedInPageData{
		Person:   p,
		Location: string(loc),
		Time:     time.Now().Format(time.RFC3339),
	}

	address := fmt.Sprintf("%v, %v %v", p.Street, p.PLZ, p.City)
	name := fmt.Sprintf("%v %v", p.Firstname, p.Lastname)

	journal.LogInToJournal(&journal.Credentials{
		Login:    true,
		Name:     name,
		Address:  address,
		Location: location.Location(data.Location),
		TimeCome: time.Now(),
	})

	checkedInTemplate.Execute(rw, data)

}

func checkedOutHandler(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	p := person.P{
		Firstname: r.PostFormValue(person.FirstNameKey),
		Lastname:  r.PostFormValue(person.LastNameKey),
		Street:    r.PostFormValue(person.StreetKey),
		PLZ:       r.PostFormValue(person.PlzKey),
		City:      r.PostFormValue(person.CityKey),
	}

	data := CheckedoutPageData{
		Person:   p,
		Location: r.PostFormValue(locationKey),
	}

	address := fmt.Sprintf("%v, %v %v", p.Street, p.PLZ, p.City)
	name := fmt.Sprintf("%v %v", p.Firstname, p.Lastname)

	journal.LogOutToJournal(&journal.Credentials{
		Login:    false,
		Name:     name,
		Address:  address,
		Location: location.Location(data.Location),
		TimeGone: time.Now(),
	})

	checkedOutTemplate.Execute(rw, data)
}

func tokenValidationWrapper(validator token.Validator, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		t, _ := url.QueryUnescape(r.URL.Query().Get(tokenKey))
		l, _ := url.QueryUnescape(r.URL.Query().Get(locationKey))

		tok := token.Token(t)
		loc := location.Location(l)

		if valid := validator(tok, loc); valid {
			ctx := context.WithValue(r.Context(), locationContextKey, loc)
			ctx = context.WithValue(ctx, tokenContextKey, tok)
			handler(w, r.WithContext(ctx))
		} else {
			http.Error(w,
				http.StatusText(http.StatusBadRequest)+" unvalid token",
				http.StatusBadRequest)
		}
	}
}

func validateFormInput(p person.P) bool {

	namePattern := regexp.MustCompile(`^[\wÄÖÜäöüß\-\s]+$`)
	streetPattern := regexp.MustCompile(`^[\wÄÖÜäöüß\-\s.]+ [0-9]+[A-Za-z]*$`)
	plzPattern := regexp.MustCompile("^[0-9]{5}$")

	first := namePattern.MatchString(p.Firstname)
	last := namePattern.MatchString(p.Lastname)
	city := namePattern.MatchString(p.City)
	street := streetPattern.MatchString(p.Street)
	plz := plzPattern.MatchString(p.PLZ)

	return first && last && city && street && plz
}
