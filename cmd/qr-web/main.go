package main

import (
	"DHBW_Golang_Project/pkg/token"
	"flag"
	"fmt"
	"github.com/skip2/go-qrcode"
	"html/template"
	"log"
	"net/http"
	"time"
)

var templateQR = template.Must(
	template.New("qrCode").Parse(`<!DOCTYPE html>
		<html lang="en">
			<head>
    			<meta charset="UTF-8">
    			<meta http-equiv="refresh" content="{{.RefreshTime}}">
			</head>
			<body>
				<img src="/pic/qr_code.jpg" alt="QR-Code" width="256" height="256">
			</body>
		</html>`))

type StartVariableStruct struct {
	RefreshTime int
	Port        string
}

var CheckinUrl string

var StartVariable StartVariableStruct

func main() {
	flag.Parse()
	go reloadQR()

	if err := http.ListenAndServeTLS(":"+StartVariable.Port, "../../assets/ssl/server.crt", "../../assets/ssl/server.key", qrMux()); err != nil {
		log.Fatal(err)
	}
}

func init() {
	var rT int
	var p string
	flag.IntVar(&rT, "refreshTime", 60, "Start variable for refresh time.")
	flag.StringVar(&p, "port", "8142", "Start variable for URL-Link")
	StartVariable = StartVariableStruct{RefreshTime: rT, Port: p}
}

func reloadQR() {
	createUrl()
	ticker := time.NewTicker(time.Duration(StartVariable.RefreshTime * 1000000000))
	for _ = range ticker.C {
		createUrl()
	}
}

func createUrl() {
	CheckinUrl = "localhost:" + StartVariable.Port + "?token=" + fmt.Sprint(token.CreateToken("DE"))
	err := qrcode.WriteFile(CheckinUrl, qrcode.Medium, 256, "pic/qr_code.jpg")

	if err != nil {
		log.Fatal(err)
	}
}

func handleQR(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../../web/templates/qr-code.html"))

	if err := tmpl.Execute(w, StartVariable); err != nil {
		log.Fatal(err)
	}
}

func qrMux() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/pic/", http.StripPrefix("/pic/", http.FileServer(http.Dir("./pic"))))

	mux.HandleFunc("/", handleQR)
	fmt.Printf("Starting server at port %v\n", StartVariable.Port)

	return mux
}
