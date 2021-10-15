package main

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"log"
	"net/http"
	"time"
)

var portQR = "8142"




/*func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET"{
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/hello", mainHandler)

	fmt.Print("Starting server at Port 8844\n")
	if err := http.ListenAndServe(":8844", nil); err != nil{
		log.Fatal(err)
	}
}*/


func reloadQR() {
	c := true
	ticker := time.NewTicker(5 * time.Second)
	for _ = range ticker.C {

		if c {
			fmt.Println("google")
			err := qrcode.WriteFile("https://google.org", qrcode.Medium, 256, "./pic/qr_code.jpg")

			if err != nil {
				log.Fatal(err)
			}
			c = false
		} else {
			fmt.Println("example")
			err := qrcode.WriteFile("https://example.org", qrcode.Medium, 256, "./pic/qr_code.jpg")

			if err != nil {
				log.Fatal(err)
			}
			c = true
		}
	}
}

func OpenServer(){
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.Handle("/pic/", http.StripPrefix("/pic/", http.FileServer(http.Dir("./pic"))))

	fmt.Printf("Starting server at port %v\n", portQR)
	if err := http.ListenAndServe(":"+portQR, nil); err != nil {
		log.Fatal(err)
	}
}

func main(){

		go reloadQR()

		OpenServer()

}