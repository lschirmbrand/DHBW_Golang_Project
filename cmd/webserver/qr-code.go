package main

import (
	"DHBW_Golang_Project/pkg/config"
	"DHBW_Golang_Project/pkg/location"
	"DHBW_Golang_Project/pkg/token"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/skip2/go-qrcode"
)

var (
	qrTemplate     *template.Template
	pathToQr       string
	pathToLocation string
	CheckinUrl     string
)

type qrCodePageData struct {
	RefreshTime int
	Location    string
}

func qrMux() http.Handler {

	parseQrTemplates("web/templates")
	pathToPic("assets/qr-codes/")
	pathToLocations("assets/")

	go reloadQR()

	mux := http.NewServeMux()
	mux.Handle("/qr-codes/", http.StripPrefix("/qr-codes/", http.FileServer(http.Dir("assets/qr-codes/"))))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	mux.HandleFunc("/qr", handleQR)

	return mux
}

func parseQrTemplates(templateDir string) {
	qrTemplate = template.Must(template.ParseFiles(path.Join(templateDir, "qr-code.html")))
}

func handleQR(w http.ResponseWriter, r *http.Request) {
	actualLocation := location.Location(r.URL.Query().Get("location"))

	if valideLocation(actualLocation) {

		data := qrCodePageData{
			RefreshTime: *config.RefreshTime,
			Location:    string(actualLocation),
		}
		qrTemplate.Execute(w, data)
	}
}

func valideLocation(expLocations location.Location) bool {
	locations, _ := location.ReadLocations(pathToLocation + "locations.xml")
	for _, actLocation := range locations {
		if actLocation == expLocations {
			return true
		}
	}
	return false
}

func pathToPic(pictureDir string) {
	pathToQr = pictureDir
}

func pathToLocations(locationsDir string) {
	pathToLocation = locationsDir
}

func reloadQR() {
	createUrl()
	ticker := time.NewTicker(time.Duration(*config.RefreshTime * 1000000000))
	for range ticker.C {
		createUrl()
	}
}

func createUrl() {
	locations, _ := location.ReadLocations(pathToLocation + "locations.xml")
	for _, loc := range locations {
		url := fmt.Sprintf("https://localhost:%v/checkin?token=%v&location=%v", *config.CheckinPort, token.CreateToken(loc), loc)
		qrcode.WriteFile(url, qrcode.Medium, 256, pathToQr+string(loc)+".jpg")
	}
}