package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestQrHandler(t *testing.T){
	req, err := http.NewRequest("GET", "http://localhost", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(handleQR)

	handler.ServeHTTP(rec, req)
	assert.Equal(t, 200, rec.Code)
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
	prevUrl := ""
	StartVariable = StartVariableStruct{2,"4122"}
	go reloadQR()

	time.Sleep(1 * time.Second)

	assert.NotEqual(t, prevUrl, CheckinUrl)
	prevUrl = CheckinUrl

	time.Sleep(2 * time.Second)

	assert.NotEqual(t, prevUrl, CheckinUrl)
	prevUrl = CheckinUrl
}
