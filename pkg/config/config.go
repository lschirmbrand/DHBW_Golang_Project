package config

import (
	"flag"
	"time"
)

const (
	DATEFORMAT         = "2006-01-02"
	DATEFORMATWITHTIME = "02-01-2006 15:04:05"

	defaultCheckinPort = 8443
	defaultQrCodePort  = 8444

	defaultRefreshTime    = 60
	defaultCookieLifetime = 24 * 356

	defaultTempaltePath     = "web/templates"
	defaultQrCodePath       = "assets/qr-codes"
	defaultLocationFilePath = "assets/location.xml"

	defaultOperation = "Visitor"
	defaultQuery     = ""

	defaultLogsPath = "logs"
)

var (
	CheckinPort *int
	QRCodePort  *int

	RefreshTime    *int
	CookieLifetime *int

	TemplatePath     *string
	QrCodePath       *string
	LocationFilePath *string

	LogPath *string

	Date      *string
	Operation *string
	Query     *string

	parsedWeb          bool = false
	parsedAnalysisTool bool = false
)

func ConfigureWeb() {
	if !parsedWeb {
		CheckinPort = flag.Int("checkinPort", defaultCheckinPort, "port of checkin server")
		QRCodePort = flag.Int("qrCodePort", defaultQrCodePort, "port of qr-code server")

		RefreshTime = flag.Int("refreshTime", defaultRefreshTime, "refresh time for qr-codes in seconds")
		CookieLifetime = flag.Int("cookieLifeTime", defaultCookieLifetime, "Lifetime of Cookies in Hours")

		TemplatePath = flag.String("templatePath", defaultTempaltePath, "path to html-template directory")
		QrCodePath = flag.String("qrCodePath", defaultQrCodePath, "path to save qr-codes")
		LocationFilePath = flag.String("locationFilePath", defaultLocationFilePath, "path to xml file with locations")

		configureMutual()

		flag.Parse()
		parsedWeb = true
	}
}

func ConfigureAnalysisTool() {
	if !parsedAnalysisTool {
		Date = flag.String("date", time.Now().Format("2006-01-02"), "Date of the requested query. Format: YYYY-MM-DD")
		Operation = flag.String("operation", defaultOperation, "Operation of the requested query. Format: Visitor or Location")
		Query = flag.String("query", defaultQuery, "The keyword of the requested query.")

		configureMutual()

		flag.Parse()
		parsedAnalysisTool = true
	}
}

func configureMutual() {
	LogPath = flag.String("logPath", defaultLogsPath, "path to logfile directory")

}
