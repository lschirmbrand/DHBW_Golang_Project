package qrweb

import (
	"DHBW_Golang_Project/pkg/config"
	"DHBW_Golang_Project/pkg/location"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQrHandler(t *testing.T) {
	config.Configure()
	parseTemplates("../../web/templates")
	location.ReadLocations("../../assets/locations.xml")

	for _, actLocation := range location.Locations {
		req, err := http.NewRequest("GET", "http://localhost:8142/qr?location="+string(actLocation), nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(handleQR)

		handler.ServeHTTP(rec, req)
		assert.Equal(t, 200, rec.Code)
	}
}

func TestMux(t *testing.T) {
	config.Configure()
	// overwrite template dir path
	templatePath := "../../web/templates"
	config.TemplatePath = &templatePath

	mux := Mux()
	server := httptest.NewServer(mux)

	rec, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)

	_, err = http.DefaultClient.Do(rec)
	assert.NoError(t, err)

}

func TestQrReload(t *testing.T) {
	// overwrite path to location file
	locationFilePath := "../../assets/locations.xml"
	config.LocationFilePath = &locationFilePath

	// overwrite refreshTime
	refreshTime := 2
	config.RefreshTime = &refreshTime

	go reloadQR()

	time.Sleep(1 * time.Second)

	// copy current checkinUrls
	previous := make(map[location.Location]string)
	for loc, url := range checkinUrls {
		previous[loc] = url
	}

	time.Sleep(2 * time.Second)

	// checkinUrls should be different
	for loc, url := range checkinUrls {
		assert.NotEqual(t, previous[loc], url)
	}
}
