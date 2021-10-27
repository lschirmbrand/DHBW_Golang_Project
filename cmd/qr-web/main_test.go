package main

import (
	"DHBW_Golang_Project/pkg/location"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestQrHandler(t *testing.T){
	parseTemplates("../../web/templates")

	locations,_ := location.ReadLocations( "../../assets/locations.xml")

	for _, actLocation := range locations {
		req, err := http.NewRequest("GET", "http://localhost:8142/qr?location="+string(actLocation), nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(handleQR)

		handler.ServeHTTP(rec, req)
		assert.Equal(t, 200, rec.Code)
	}
}

func TestMux(t *testing.T){
	mux := qrMux()
	server := httptest.NewServer(mux)

	rec, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)

	_, err = http.DefaultClient.Do(rec)
	assert.NoError(t, err)

}

func TestQrReload(t *testing.T){
	pathToPic("../../assets/qr-codes/")
	pathToLocations("../../assets/")
	prevUrl := ""
	StartVariable = StartVariableStruct{2,"4122"}
	setLocation("Italy")
	go reloadQR()

	time.Sleep(1 * time.Second)

	assert.NotEqual(t, prevUrl, CheckinUrl)
	prevUrl = CheckinUrl

	time.Sleep(2 * time.Second)

	assert.NotEqual(t, prevUrl, CheckinUrl)
	prevUrl = CheckinUrl
}

func TestCheckError(t *testing.T){
	var err1 error
	isError1 := checkError(err1)
	err2 := errors.New("error")
	isError2 := checkError(err2)

	assert.True(t, isError1)
	assert.False(t, isError2)

}

func TestValideLocation(t *testing.T){
	pathToLocations("../../assets/")

	valide1 := valideLocation("test")
	valide2 := valideLocation("Germany")

	assert.False(t, valide1)
	assert.True(t, valide2)
}
