package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	checkinPort *int
)

func init() {
	checkinPort = flag.Int("checkinPort", 8443, "port of checkin server")
}

func main() {
	flag.Parse()

	fmt.Printf("Starting server for checkin at port %v\n", *checkinPort)

	if err := http.ListenAndServeTLS(":"+fmt.Sprint(*checkinPort), "assets/ssl/server.crt", "assets/ssl/server.key", checkinMux()); err != nil {
		log.Fatal(err)
	}

}
