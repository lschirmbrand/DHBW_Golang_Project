package main

import (
	"DHBW_Golang_Project/pkg/checkinweb"
	"DHBW_Golang_Project/pkg/config"
	"DHBW_Golang_Project/pkg/qrweb"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	config.Configure()

	var wg sync.WaitGroup

	// run web server for checkin/checkout
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("Starting server for checkin at port %v\n", *config.CheckinPort)

		addr := fmt.Sprintf("localhost:%v", *config.CheckinPort)

		err := http.ListenAndServeTLS(addr, *config.CertificateFilePath, *config.KeyFilePath, checkinweb.Mux())
		if err != nil {
			log.Fatal(err)
		}

	}()

	// run web server for qr-codes
	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Printf("Starting server for qr-codes at port %v\n", *config.QRCodePort)

		addr := fmt.Sprintf("localhost:%v", *config.QRCodePort)

		err := http.ListenAndServeTLS(addr, *config.CertificateFilePath, *config.KeyFilePath, qrweb.Mux())
		if err != nil {
			log.Fatal(err)
		}

	}()

	wg.Wait()
}
