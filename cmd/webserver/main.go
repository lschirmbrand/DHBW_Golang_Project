package main

import (
	"DHBW_Golang_Project/cmd/webserver/checkinweb"
	"DHBW_Golang_Project/cmd/webserver/qrweb"
	"DHBW_Golang_Project/internal/config"
	"DHBW_Golang_Project/internal/journal"
	"DHBW_Golang_Project/internal/location"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	config.ConfigureWeb()

	jour := journal.NewLogFileJournal(*config.LogPath)
	locationStore := location.NewLocationStore(*config.LocationFilePath)

	checkinMuxCfg := checkinweb.CheckInMuxCfg{
		TempaltePath:   *config.TemplatePath,
		CookieLifetime: *config.CookieLifetime,
	}

	qrMuxCfg := qrweb.QrMuxCfg{
		TemplatePath: *config.TemplatePath,
		QrCodePath:   *config.QrCodePath,
		RefreshTime:  *config.RefreshTime,
		CheckInPort:  *config.CheckinPort,
	}

	var wg sync.WaitGroup

	// run web server for checkin/checkout
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("Starting server for checkin at port %v\n", *config.CheckinPort)

		addr := fmt.Sprintf("localhost:%v", *config.CheckinPort)

		checkinweb.Setup(jour, locationStore, &checkinMuxCfg)

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

		qrweb.Setup(locationStore, qrMuxCfg)

		err := http.ListenAndServeTLS(addr, *config.CertificateFilePath, *config.KeyFilePath, qrweb.Mux())
		if err != nil {
			log.Fatal(err)
		}

	}()

	wg.Wait()
}
