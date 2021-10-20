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

type startVariableStruct struct {
	RefreshTime int
	Port string
}

var startVariable startVariableStruct
var templateQR = template.Must(
	template.New("qrCode").Parse(`<!DOCTYPE html>
		<html lang="en">
			<head>
    			<meta charset="UTF-8">
    			<meta http-equiv="refresh" content="{{.RefreshTime}}; URL=localhost:{{.UrlLink}}/">
			</head>
			<body>
				<img src="../pic/qr_code.jpg" alt="QR-Code" width="256" height="256">
			</body>
		</html>`))

func main() {
	t :=token.CreateToken("DE")
	s := token.CreateToken("IT")
	b:=token.VerifyToken(t)
	b=token.VerifyToken(s)
	print(b)
	var rT int
	var p string
	flag.IntVar(&rT, "refreshTime", 60, "Start variable for refresh time.")
	flag.StringVar(&p, "port", "8142", "Start variable for URL-Link")
	flag.Parse()

	startVariable = startVariableStruct{RefreshTime: rT, Port: p}

	go reloadQR()

	if err := http.ListenAndServeTLS(":"+startVariable.Port, "assets/ssl/server.crt", "assets/ssl/server.key", checkinMux()); err != nil {
		log.Fatal(err)
	}

}


func reloadQR() {
	ticker := time.NewTicker(time.Duration(startVariable.RefreshTime * 1000000000))
	for _ = range ticker.C {
		tokenURL := "localhost:" + startVariable.Port + "?token=" +  fmt.Sprint(token.CreateToken("DE"))
		err := qrcode.WriteFile(tokenURL, qrcode.Medium, 256, "./pic/qr_code.jpg")

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(tokenURL)
	}
}



func handleQR(w http.ResponseWriter, r *http.Request){
	templateQR.ExecuteTemplate(w, "qrCode", startVariable)
}

func checkinMux() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/pic/", http.StripPrefix("/pic/", http.FileServer(http.Dir("./pic"))))

	mux.HandleFunc("/", handleQR)
	fmt.Printf("Starting server at port %v\n", startVariable.Port)
	if err := http.ListenAndServe(":"+startVariable.Port, nil); err != nil {
		log.Fatal(err)
	}

	return mux
}


