package main

import (
	"net/http"
	"text/template"
)

type CheckInPageData struct {
	PlaceName string
}

func checkinMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/checkin", func(rw http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("resources/templates/checkin.html"))

		data := CheckInPageData{PlaceName: "LOLOLOL"}

		tmpl.Execute(rw, data)
	})

	fs := http.FileServer(http.Dir("resources/static"))
	mux.Handle("/static", http.StripPrefix("/static", fs))

	return mux
}
