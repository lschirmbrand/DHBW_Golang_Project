package main

import (
	"DHBW_Golang_Project/pkg/token"
	"context"
	"net/http"
	"text/template"
	"time"
)

type CheckInPageData struct {
	Location string
}

type CheckOutPageData struct {
	Location string
	Name     string
	Street   string
	PLZ      string
	Time     string
}

func checkinMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/checkin", tokenValidationWrapper(token.Validate, checkInHandler))
	mux.HandleFunc("/checkout", checkOutHandler)

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}

func checkInHandler(rw http.ResponseWriter, r *http.Request) {

	location := r.Context().Value("location").(string)

	tmpl := template.Must(template.ParseFiles("web/templates/checkin.html"))

	data := CheckInPageData{location}

	tmpl.Execute(rw, data)
}

func checkOutHandler(rw http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	tmpl := template.Must(template.ParseFiles("web/templates/checkout.html"))

	data := CheckOutPageData{
		Location: r.PostForm.Get("location"),
		Name:     r.PostFormValue("name"),
		Street:   r.PostForm.Get("street"),
		PLZ:      r.PostForm.Get("plz"),
		Time:     time.Now().Format(time.RFC3339),
	}

	tmpl.Execute(rw, data)
}

func tokenValidationWrapper(validator token.Validator, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := token.Token(r.URL.Query().Get("token"))
		if valid, location := validator(t); valid {
			ctx := context.WithValue(r.Context(), "location", location)

			handler(w, r.WithContext(ctx))
		} else {
			http.Error(w,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
		}
	}
}
