package main

import (
	"log"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertQueryExport(t *testing.T) {
	res := make([]string, 0)
	assert.False(t, assertQueryExport(&res))
}

func TestTrimStringBasedOnOS(t *testing.T) {
	if runtime.GOOS == "windows" {
		res := trimStringBasedOnOS("teststring\r\n", true)
		assert.EqualValues(t, res, "teststring")
	} else {
		res := trimStringBasedOnOS("teststring\n", true)
		assert.EqualValues(t, res, "teststring")
	}
	res := trimStringBasedOnOS("\nteststring", false)
	assert.EqualValues(t, res, "teststring")
}

func checkErrorForTest(err error) {
	if err != nil {
		log.Fatalln(err)
	} else {
		return
	}
}
