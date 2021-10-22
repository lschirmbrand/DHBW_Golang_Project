package main

import (
	"DHBW_Golang_Project/pkg/token"
	"context"
	"html/template"
	"net/http"
	"path"
	"time"
)

type CheckInPageData struct {
	Person
	Location string
}

type CheckOutPageData struct {
	Person
	Location string
	Time     string
}

type Person struct {
	Name   string
	Street string
	PLZ    string
	City   string
}

type contextKey string

const (
	locationContextKey contextKey = "location"
)

var checkInTemplate *template.Template
var checkOutTemplate *template.Template

func checkinMux() http.Handler {
	mux := http.NewServeMux()

	parseTemplates("web/templates")

	mux.HandleFunc("/checkin", tokenValidationWrapper(token.Validate, checkInHandler))
	mux.HandleFunc("/checkout", checkOutHandler)

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}

func parseTemplates(templateDir string) {
	checkInTemplate = template.Must(template.ParseFiles(path.Join(templateDir, "checkin.html")))
	checkOutTemplate = template.Must(template.ParseFiles(path.Join(templateDir, "checkOut.html")))
}

func checkInHandler(rw http.ResponseWriter, r *http.Request) {

	location := r.Context().Value(locationContextKey).(string)

	p := readPersonFromCookies(r)
	data := CheckInPageData{Person: *p, Location: location}

	checkInTemplate.Execute(rw, data)
}

func checkOutHandler(rw http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	p := Person{
		Name:   r.PostFormValue("name"),
		Street: r.PostFormValue("street"),
		PLZ:    r.PostFormValue("plz"),
		City:   r.PostFormValue("city"),
	}

	data := CheckOutPageData{
		Person:   p,
		Location: r.PostForm.Get("location"),
		Time:     time.Now().Format(time.RFC3339),
	}

	savePersonToCookies(rw, &p)

	checkOutTemplate.Execute(rw, data)

}

func savePersonToCookies(rw http.ResponseWriter, p *Person) {
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

	name, err := r.Cookie("name")
	if err == nil {
		p.Name = name.Value
	}

	street, err := r.Cookie("street")
	if err == nil {
		p.Street = street.Value
	}

	plz, err := r.Cookie("plz")
	if err == nil {
		p.PLZ = plz.Value
	}

	city, err := r.Cookie("city")
	if err == nil {
		p.City = city.Value
	}

	return &p

}

func tokenValidationWrapper(validator token.Validator, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := token.Token(r.URL.Query().Get("token"))
		if valid, location := validator(t); valid {
			ctx := context.WithValue(r.Context(), locationContextKey, location)

			handler(w, r.WithContext(ctx))
		} else {
			http.Error(w,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
		}
	}
}
