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
	Port string
}

var startVariable StartVariableStruct


func main() {
	var rT int
	var p string
	flag.IntVar(&rT, "refreshTime", 60, "Start variable for refresh time.")
	flag.StringVar(&p, "port", "8142", "Start variable for URL-Link")
	flag.Parse()

	startVariable = StartVariableStruct{RefreshTime: rT, Port: p}

	go reloadQR()

	if err := http.ListenAndServeTLS(":"+startVariable.Port, "../../assets/ssl/server.crt", "../../assets/ssl/server.key", qrMux()); err != nil {
		log.Fatal(err)
	}
	//openServer()
}

func reloadQR() {
	ticker := time.NewTicker(time.Duration(startVariable.RefreshTime * 1000000000))
	for _ = range ticker.C {
		checkinUrl := "localhost:" + startVariable.Port + "?token=" +  fmt.Sprint(token.CreateToken("DE"))
		err := qrcode.WriteFile(checkinUrl, qrcode.Medium, 256, "pic/qr_code.jpg")

		if err != nil {
			log.Fatal(err)
		}
	}
}



func handleQR(w http.ResponseWriter, r *http.Request){
	tmpl := template.Must(template.ParseFiles("../../web/templates/qr-code.html"))

	if err := tmpl.Execute(w, startVariable); err != nil{
		log.Fatal(err)
	}
}

func qrMux() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/pic/", http.StripPrefix("/pic/", http.FileServer(http.Dir("./pic"))))

	mux.HandleFunc("/", handleQR)
	fmt.Printf("Starting server at port %v\n", startVariable.Port)

	return mux
}







