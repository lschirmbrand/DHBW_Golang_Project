package config

import "flag"

var (
	CheckinPort *int
	QRCodePort  *int

	RefreshTime *int

	TemplatePath *string
	QrCodePath   *string
)

func Configure() {
	CheckinPort = flag.Int("checkinPort", 8443, "port of checkin server")
	QRCodePort = flag.Int("qrCodePort", 8444, "port of qr-code server")

	RefreshTime = flag.Int("refreshTime", 60, "refresh time for qr-codes in seconds")

	TemplatePath = flag.String("templatePath", "web/templates", "path to html-template directory")
	QrCodePath = flag.String("qrCodePath", "assets/qr-codes", "path to save qr-codes")

	flag.Parse()
}
