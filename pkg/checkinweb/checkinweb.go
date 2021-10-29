package checkinweb

import (
	"DHBW_Golang_Project/pkg/config"
	"DHBW_Golang_Project/pkg/location"
	"DHBW_Golang_Project/pkg/token"
	"context"
	"html/template"
	"net/http"
	"path"
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
	Name   string
	Street string
	PLZ    string
	City   string
}

type key string

const (
	locationKey key = "location"
	nameKey     key = "name"
	streetKey   key = "street"
	plzKey      key = "plz"
	cityKey     key = "city"
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

func checkInHandler(rw http.ResponseWriter, r *http.Request) {

	l := r.Context().Value(locationKey).(string)

	p := readPersonFromCookies(r)

	data := CheckInPageData{Person: *p, Location: location.Location(l)}

	checkInTemplate.Execute(rw, data)
}

func checkedInHandler(rw http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	p := Person{
		Name:   r.PostFormValue("name"),
		Street: r.PostFormValue("street"),
		PLZ:    r.PostFormValue("plz"),
		City:   r.PostFormValue("city"),
	}

	loc := r.PostFormValue("location")

	data := CheckedInPageData{
		Person:   p,
		Location: location.Location(loc),
		Time:     time.Now().Format(time.RFC3339),
	}

	savePersonToCookies(rw, &p)

	// journal.LogToJournal(journal.Credentials{
	// 	Name:    p.Name,
	// 	Address: location,
	// }, false)

	checkedInTemplate.Execute(rw, data)

}

func checkedOutHandler(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	data := CheckedoutPageData{
		Person:   Person{Name: r.PostFormValue("name")},
		Location: location.Location(r.PostFormValue("location")),
	}

	checkedOutTemplate.Execute(rw, data)
}

func savePersonToCookies(rw http.ResponseWriter, p *Person) {
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

	http.SetCookie(rw, &nameCookie)
	http.SetCookie(rw, &streetCookie)
	http.SetCookie(rw, &plzCookie)
	http.SetCookie(rw, &cityCookie)
}

func readPersonFromCookies(r *http.Request) *Person {
	p := Person{
		Name:   "",
		Street: "",
		PLZ:    "",
		City:   "",
	}

	name, err := r.Cookie(string(nameKey))
	if err == nil {
		p.Name = name.Value
	}

	street, err := r.Cookie(string(streetKey))
	if err == nil {
		p.Street = street.Value
	}

	plz, err := r.Cookie(string(plzKey))
	if err == nil {
		p.PLZ = plz.Value
	}

	city, err := r.Cookie(string(cityKey))
	if err == nil {
		p.City = city.Value
	}

	return &p

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
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
		}
	}
}
