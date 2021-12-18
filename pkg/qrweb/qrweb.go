package qrweb

import (
	"DHBW_Golang_Project/pkg/config"
	"DHBW_Golang_Project/pkg/location"
	"DHBW_Golang_Project/pkg/token"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/skip2/go-qrcode"
)

var (
	qrTemplate, locsTemplate *template.Template
	checkinUrls              map[location.Location]string
)

type qrCodePageData struct {
	RefreshTime int
	Location    string
	CheckInUrl  string
}

type locationsPageData struct {
	Location string
	Url      string
}

func Mux() http.Handler {

	parseTemplates(*config.TemplatePath)
	location.ReadLocations(*config.LocationFilePath)
	go reloadQR()

	mux := http.NewServeMux()
	mux.Handle("/qr-codes/", http.StripPrefix("/qr-codes/", http.FileServer(http.Dir(*config.QrCodePath))))
	// mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	mux.HandleFunc("/qr", handleQR)
	mux.HandleFunc("/locations", handleLocations)

	return mux
}

func parseTemplates(templateDir string) {
	qrTemplate = template.Must(template.ParseFiles(path.Join(templateDir, "qr-code.html")))
	locsTemplate = template.Must(template.ParseFiles(path.Join(templateDir, "locations.html")))
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

func handleLocations(w http.ResponseWriter, r *http.Request) {
	locs := location.Locations

	data := make([]locationsPageData, len(locs))

	for i, loc := range locs {
		data[i] = locationsPageData{
			Location: string(loc),
			Url:      "qr?location=" + url.QueryEscape(string(loc)),
		}
	}

	locsTemplate.Execute(w, data)
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
		url := fmt.Sprintf("https://localhost:%v/checkin?token=%v&location=%v",
			*config.CheckinPort,
			url.QueryEscape(string(token.CreateToken(loc))),
			url.QueryEscape(string(loc)),
		)
		if _, err := os.Stat(*config.QrCodePath); os.IsNotExist(err) {
			os.MkdirAll(*config.QrCodePath, 0755)
		}
		qrcode.WriteFile(url, qrcode.Medium, 256, path.Join(*config.QrCodePath, string(loc)+".jpg"))
		checkinUrls[loc] = url
	}
}
