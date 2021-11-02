package config

import (
	"flag"
)

var (
	CheckinPort *int
	QRCodePort  *int

	RefreshTime    *int
	CookieLifetime *int

	TemplatePath     *string
	QrCodePath       *string
	LocationFilePath *string

	parsed bool
)

func Configure() {
	if !parsed {
		CheckinPort = flag.Int("checkinPort", 8443, "port of checkin server")
		QRCodePort = flag.Int("qrCodePort", 8444, "port of qr-code server")

		RefreshTime = flag.Int("refreshTime", 60, "refresh time for qr-codes in seconds")
		CookieLifetime = flag.Int("cookieLifeTime", 24*356, "Lifetime of Cookies in Hours")

		TemplatePath = flag.String("templatePath", "web/templates", "path to html-template directory")
		QrCodePath = flag.String("qrCodePath", "assets/qr-codes", "path to save qr-codes")
		LocationFilePath = flag.String("locationFilePath", "assets/locations.xml", "path to xml file with locations")
	}
	flag.Parse()
	parsed = true
}
