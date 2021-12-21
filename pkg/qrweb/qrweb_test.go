package qrweb

import (
	"DHBW_Golang_Project/pkg/location"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var st *location.LocationStore

func SetupTest() {
	st = location.NewLocationStore("./test_assets/locations.xml")

	Setup(st, QrMuxCfg{
		TemplatePath: "../../web/templates",
		QrCodePath:   "./test_assets/qr-codes",
		RefreshTime:  2,
		CheckInPort:  8443,
	})
}

func TestQrHandler(t *testing.T) {

	SetupTest()

	for _, actLocation := range st.Locations {
		req, err := http.NewRequest("GET", "http://localhost:8142/qr?location="+string(actLocation), nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(handleQR)

		handler.ServeHTTP(rec, req)
		assert.Equal(t, 200, rec.Code)
	}
}

func TestMux(t *testing.T) {
	SetupTest()

	mux := Mux()
	server := httptest.NewServer(mux)

	rec, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)

	_, err = http.DefaultClient.Do(rec)
	assert.NoError(t, err)

}

func TestQrReload(t *testing.T) {
	SetupTest()

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
