package config

import (
	"flag"
)

var (
	Hostname    *string
	CheckinPort *int
	QRCodePort  *int

	CertificateFilePath *string
	KeyFilePath         *string

	RefreshTime    *int
	CookieLifetime *int

	TemplatePath     *string
	QrCodePath       *string
	LocationFilePath *string

	parsed bool
)

func Configure() {
	if !parsed {
		Hostname = flag.String("hostname", "localhost", "hostname of application")
		CheckinPort = flag.Int("checkinPort", 8443, "port of checkin server")
		QRCodePort = flag.Int("qrCodePort", 8444, "port of qr-code server")

		CertificateFilePath = flag.String("certificateFile", "assets/ssl/cert.pem", "path to ssl certificate")
		KeyFilePath = flag.String("keyFile", "assets/ssl/key.pem", "path to ssl key")

		RefreshTime = flag.Int("refreshTime", 60, "refresh time for qr-codes in seconds")
		CookieLifetime = flag.Int("cookieLifeTime", 24*356, "Lifetime of Cookies in Hours")

		TemplatePath = flag.String("templatePath", "web/templates", "path to html-template directory")
		QrCodePath = flag.String("qrCodePath", "assets/qr-codes", "path to save qr-codes")
		LocationFilePath = flag.String("locationFile", "assets/locations.xml", "path to xml file with locations")
	}
	flag.Parse()
	parsed = true
}
