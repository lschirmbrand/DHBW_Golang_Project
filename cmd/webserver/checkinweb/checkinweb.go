package checkinweb

import (
	"DHBW_Golang_Project/internal/journal"
	"DHBW_Golang_Project/internal/location"
	"DHBW_Golang_Project/internal/person"
	"DHBW_Golang_Project/internal/token"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"time"
)

/*
	Erstellt von: 	8864957
	Created by:		8864957

	also: 4775194, 9514094
*/

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

	jour          journal.Journal
	cookieStore   *CookieStore
	locationStore *location.LocationStore
)

type CheckInMuxCfg struct {
	TempaltePath   string
	CookieLifetime int
}

func Setup(j journal.Journal, locStore *location.LocationStore, cfg *CheckInMuxCfg) {
	jour = j
	parseTemplates(cfg.TempaltePath)
	cookieStore = NewCookieStore(cfg.CookieLifetime)
	locationStore = locStore
}

func Mux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/checkin", tokenValidationWrapper(token.Validate, checkInHandler))
	mux.HandleFunc("/checkedin", checkedInHandler)
	mux.HandleFunc("/checkedout", checkedOutHandler)

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
	pers := cookieStore.ReadFromCookies(r)

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
		Firstname: r.PostFormValue(cookieStore.FirstNameKey),
		Lastname:  r.PostFormValue(cookieStore.LastNameKey),
		Street:    r.PostFormValue(cookieStore.StreetKey),
		PLZ:       r.PostFormValue(cookieStore.PlzKey),
		City:      r.PostFormValue(cookieStore.CityKey),
	}

	loc := location.Location(r.PostFormValue(locationKey))

	// validate location
	if !locationStore.Validate(loc) {
		http.Error(rw,
			http.StatusText(http.StatusBadRequest)+"not a valid location",
			http.StatusBadRequest)
	}

	cookieStore.SaveToCookies(rw, &p)

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

	address := fmt.Sprintf("%v %v %v", p.Street, p.PLZ, p.City)
	name := fmt.Sprintf("%v %v", p.Firstname, p.Lastname)

	jour.LogIn(&journal.Credentials{
		Checkin:   true,
		Name:      name,
		Address:   address,
		Location:  location.Location(data.Location),
		Timestamp: time.Now(),
	})

	checkedInTemplate.Execute(rw, data)

}

func checkedOutHandler(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	p := person.P{
		Firstname: r.PostFormValue(cookieStore.FirstNameKey),
		Lastname:  r.PostFormValue(cookieStore.LastNameKey),
		Street:    r.PostFormValue(cookieStore.StreetKey),
		PLZ:       r.PostFormValue(cookieStore.PlzKey),
		City:      r.PostFormValue(cookieStore.CityKey),
	}

	data := CheckedoutPageData{
		Person:   p,
		Location: r.PostFormValue(locationKey),
	}

	address := fmt.Sprintf("%v %v %v", p.Street, p.PLZ, p.City)
	name := fmt.Sprintf("%v %v", p.Firstname, p.Lastname)

	jour.LogOut(&journal.Credentials{
		Checkin:   false,
		Name:      name,
		Address:   address,
		Location:  location.Location(data.Location),
		Timestamp: time.Now(),
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
