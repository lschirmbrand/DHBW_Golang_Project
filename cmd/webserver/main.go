package main

import (
	"DHBW_Golang_Project/pkg/checkinweb"
	"DHBW_Golang_Project/pkg/config"
	"DHBW_Golang_Project/pkg/qrweb"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.ConfigureWeb()

	finished := make(chan bool)
	// run web server for checkin/checkout
	go func(finished chan<- bool) {
		fmt.Printf("Starting server for checkin at port %v\n", *config.CheckinPort)

		err := http.ListenAndServeTLS(":"+fmt.Sprint(*config.CheckinPort), "assets/ssl/server.crt", "assets/ssl/server.key", checkinweb.Mux())
		if err != nil {
			log.Fatal(err)
		}

		finished <- true
	}(finished)

	// run web server for qr-codes
	go func(finished chan<- bool) {
		fmt.Printf("Starting server for qr-codes at port %v\n", *config.QRCodePort)

		err := http.ListenAndServeTLS(":"+fmt.Sprint(*config.QRCodePort), "assets/ssl/server.crt", "assets/ssl/server.key", qrweb.Mux())
		if err != nil {
			log.Fatal(err)
		}

		finished <- true
	}(finished)

	<-finished
	<-finished

}
