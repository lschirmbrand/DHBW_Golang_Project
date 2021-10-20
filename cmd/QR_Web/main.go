package main

import (
	"flag"
	"fmt"
	"github.com/skip2/go-qrcode"
	"html/template"
	"log"
	"net/http"
	"time"
)

type startVariableStruct struct {
	RefreshTime int
	UrlLink string
}

var startVariable startVariableStruct
var templateQR = template.Must(
	template.New("qrCode").Parse(`<!DOCTYPE html>
		<html lang="en">
			<head>
    			<meta charset="UTF-8">
    			<meta http-equiv="refresh" content="{{.RefreshTime}}; URL=http://localhost:{{.UrlLink}}/">
			</head>
			<body>
				<img src="../pic/qr_code.jpg" alt="QR-Code" width="256" height="256">
			</body>
		</html>`))

func main() {
	var rT int
	var uL string
	flag.IntVar(&rT, "refreshTime", 60, "Start variable for refresh time.")
	flag.StringVar(&uL, "urlLink", "8142", "Start variable for URL-Link")
	flag.Parse()

	startVariable = startVariableStruct{RefreshTime: rT, UrlLink: uL}

	go reloadQR()
	OpenServer()

}


func reloadQR() {
	ticker := time.NewTicker(time.Duration(startVariable.RefreshTime * 1000000000))
	for _ = range ticker.C {
		tokenURL := "localhost:" + startVariable.UrlLink + "?token=" +  string(randToken())
		err := qrcode.WriteFile(tokenURL, qrcode.Medium, 256, "./pic/qr_code.jpg")

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(tokenURL)
	}
}

func randToken() []byte{
	b := make([]byte, 8)
	_, err := rand.Read(b)

	if err != nil{
		log.Fatal(err)
	}

	return b
}

func handleQR(w http.ResponseWriter, r *http.Request){
	templateQR.ExecuteTemplate(w, "qrCode", startVariable)
}

func OpenServer() {

	http.Handle("/pic/", http.StripPrefix("/pic/", http.FileServer(http.Dir("./pic"))))

	http.HandleFunc("/", handleQR)
	fmt.Printf("Starting server at port %v\n", startVariable.UrlLink)
	if err := http.ListenAndServe(":"+startVariable.UrlLink, nil); err != nil {
		log.Fatal(err)
	}
}


