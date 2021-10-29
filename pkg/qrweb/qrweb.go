package qrweb

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
	qrTemplate  *template.Template
	checkinUrls map[location.Location]string
)

type qrCodePageData struct {
	RefreshTime int
	Location    string
	CheckInUrl  string
}

func Mux() http.Handler {

	parseTemplates(*config.TemplatePath)
	location.ReadLocations(*config.LocationFilePath)
	go reloadQR()

	mux := http.NewServeMux()
	mux.Handle("/qr-codes/", http.StripPrefix("/qr-codes/", http.FileServer(http.Dir(*config.QrCodePath))))
	// mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	mux.HandleFunc("/qr", handleQR)

	return mux
}

func parseTemplates(templateDir string) {
	qrTemplate = template.Must(template.ParseFiles(path.Join(templateDir, "qr-code.html")))
}

func handleQR(w http.ResponseWriter, r *http.Request) {
	loc := location.Location(r.URL.Query().Get("location"))

	if location.Validate(loc) {

		data := qrCodePageData{
			RefreshTime: *config.RefreshTime,
			Location:    string(loc),
			CheckInUrl:  checkinUrls[loc],
		}
		qrTemplate.Execute(w, data)
	}
}

func reloadQR() {
	checkinUrls = make(map[location.Location]string)

	createUrl()
	ticker := time.NewTicker(time.Duration(*config.RefreshTime * 1000000000))
	for range ticker.C {
		createUrl()
	}
}

func createUrl() {
	for _, loc := range location.Locations {
		url := fmt.Sprintf("https://localhost:%v/checkin?token=%v&location=%v", *config.CheckinPort, token.CreateToken(loc), loc)
		qrcode.WriteFile(url, qrcode.Medium, 256, path.Join(*config.QrCodePath, string(loc)+".jpg"))
		checkinUrls[loc] = url
	}
}
