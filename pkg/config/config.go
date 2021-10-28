package config

import "flag"

var (
	CheckinPort *int
	QRCodeport  *int

	RefreshTime *int
)

func Configure() {
	CheckinPort = flag.Int("checkinPort", 8443, "port of checkin server")
	QRCodeport = flag.Int("qrCodePort", 8444, "port of qr-code server")

	RefreshTime = flag.Int("refreshTime", 60, "refresh time for qr-codes in seconds")

	flag.Parse()
}
