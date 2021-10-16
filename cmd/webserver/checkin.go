package main

import (
	"DHBW_Golang_Project/pkg/token"
	"net/http"
	"text/template"
	"time"
)

type CheckInPageData struct {
	PlaceName string
}

type CheckOutPageData struct {
	PlaceName string
	Name      string
	Street    string
	PLZ       string
	Time      string
}

func checkinMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/checkin", checkInHandler)
	mux.HandleFunc("/checkout", checkOutHandler)

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}

func checkInHandler(rw http.ResponseWriter, r *http.Request) {

	t := token.Token(r.URL.Query().Get("token"))
	isValid := token.VerifyToken(t)

	if !isValid {
		rw.Write([]byte("invalid token"))
	}

	placeName := r.URL.Query().Get("place")

	tmpl := template.Must(template.ParseFiles("web/templates/checkin.html"))

	data := CheckInPageData{PlaceName: placeName}

	tmpl.Execute(rw, data)
}

func checkOutHandler(rw http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	tmpl := template.Must(template.ParseFiles("web/templates/checkout.html"))

	data := CheckOutPageData{
		PlaceName: r.PostForm.Get("place"),
		Name:      r.PostFormValue("name"),
		Street:    r.PostForm.Get("street"),
		PLZ:       r.PostForm.Get("plz"),
		Time:      time.Now().Format(time.RFC3339),
	}

	tmpl.Execute(rw, data)
}
