package config

import (
	"flag"
	"time"
)

const (
	DATEFORMATWITHTIME = "02-01-2006 15:04:05"

	defaultCheckinPort = 8443
	defaultQrCodePort  = 8444

	defaultRefreshTime    = 60
	defaultCookieLifetime = 24 * 356

	defaultTemplatePath     = "web/templates"
	defaultQrCodePath       = "assets/qr-codes"
	defaultLocationFilePath = "assets/locations.xml"

	defaultCertificateFilePath = "assets/ssl/cert.pem"
	defaultKeyFilePath         = "assets/ssl/key.pem"

	defaultOperation = "Visitor"
	defaultQuery     = ""

	defaultLogsPath = "logs"
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

	LogPath *string

	Date      *string
	Operation *string
	Query     *string
	Testcase  *bool

	parsedWeb          = false
	parsedAnalysisTool = false
)

func ConfigureWeb() {
	if !parsedWeb {
		CheckinPort = flag.Int("checkinPort", defaultCheckinPort, "port of checkin server")
		QRCodePort = flag.Int("qrCodePort", defaultQrCodePort, "port of qr-code server")

		RefreshTime = flag.Int("refreshTime", defaultRefreshTime, "refresh time for qr-codes in seconds")
		CookieLifetime = flag.Int("cookieLifeTime", defaultCookieLifetime, "Lifetime of Cookies in Hours")

		TemplatePath = flag.String("templatePath", defaultTemplatePath, "path to html-template directory")
		QrCodePath = flag.String("qrCodePath", defaultQrCodePath, "path to save qr-codes")
		LocationFilePath = flag.String("locationFilePath", defaultLocationFilePath, "path to xml file with locations")

		CertificateFilePath = flag.String("certificateFilePath", defaultCertificateFilePath, "path to ssl certificate")
		KeyFilePath = flag.String("keyFilePath", defaultKeyFilePath, "path to ssl key")

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
		Testcase = flag.Bool("testcase", false, "Determines, whether output gets printed in the terminal.")
		LogPath = flag.String("logpath", defaultLogsPath, "The default path of the log-files.")

		configureMutual()

		flag.Parse()
		parsedAnalysisTool = true
	}
}

func configureMutual() {
	LogPath = flag.String("logPath", defaultLogsPath, "path to logfile directory")
}
