package main

import (
	"DHBW_Golang_Project/pkg/location"
	"DHBW_Golang_Project/pkg/token"
	"flag"
	"fmt"
	"github.com/skip2/go-qrcode"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

type StartVariableStruct struct {
	RefreshTime int
	Port        string
}

type htmlVariableStruct struct {
	RefreshTime int
	Location string
}

var CheckinUrl string

var actualLocation location.Location

var qrTemplate *template.Template

var pathToQr string

var pathToLocation string

var StartVariable StartVariableStruct

var htmlVariable htmlVariableStruct

func main() {

	flag.Parse()

	parseTemplates("web/templates")
	pathToPic("assets/qr-codes/")
	pathToLocations("assets/")

	go reloadQR()

	err := http.ListenAndServeTLS(":"+StartVariable.Port, "./assets/ssl/server.crt", "./assets/ssl/server.key", qrMux())

	checkError(err)
}

func setLocation(loc location.Location) {
	actualLocation = loc
}

func parseTemplates(templateDir string) {
	qrTemplate = template.Must(template.ParseFiles(path.Join(templateDir, "qr-code.html")))
}

func init() {
	var rT int
	var p string
	flag.IntVar(&rT, "refreshTime", 2, "Start variable for refresh time.")
	flag.StringVar(&p, "port", "8142", "Start variable for URL-Link")
	StartVariable = StartVariableStruct{RefreshTime: rT, Port: p}
	htmlVariable = htmlVariableStruct{RefreshTime: rT}
}

func reloadQR() {
	createUrl()
	ticker := time.NewTicker(time.Duration(StartVariable.RefreshTime * 1000000000))
	for range ticker.C {
		createUrl()
	}
}

func createUrl() {
	locations,_ := location.ReadLocations(pathToLocation + "locations.xml")
	for _, loc := range locations{
		CheckinUrl = "localhost:" + StartVariable.Port + "?token=" + string(token.CreateToken(loc))

		err := qrcode.WriteFile(CheckinUrl, qrcode.Medium, 256, pathToQr + string(loc)+".jpg")

		checkError(err)
	}
}

func checkError(err error) bool {
	if err != nil{
		log.Println(err)
		return false
	}else{
		return true
	}
}

func pathToPic(pictureDir string){
	pathToQr = pictureDir
}

func pathToLocations(locationsDir string){
	pathToLocation = locationsDir
}

func handleQR(w http.ResponseWriter, r *http.Request) {
	actualLocation = location.Location(r.URL.Query()["location"][0])

	if valideLocation(actualLocation) {
		htmlVariable.Location = string(actualLocation)

		err := qrTemplate.Execute(w, htmlVariable)
		checkError(err)
	}
}



func valideLocation(expLocations location.Location) bool{
	locations,_ := location.ReadLocations(pathToLocation + "locations.xml")
	for _, actLocation := range locations{
		if actLocation == expLocations{
			return true
		}
	}
	return false
}

func qrMux() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/qr-codes/", http.StripPrefix("/qr-codes/", http.FileServer(http.Dir("assets/qr-codes/"))))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	mux.HandleFunc("/qr", handleQR)
	fmt.Printf("Starting server at port %v\n", StartVariable.Port)
	fmt.Printf("Refresh time: %v\n", StartVariable.RefreshTime)

	return mux
}