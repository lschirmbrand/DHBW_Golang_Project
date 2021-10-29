package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"runtime"
	"testing"
)

func TestAssertQueryExport(t *testing.T){
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

func TestRequestHelp(t *testing.T){
	args := make([]string, 0)
	assert.False(t, requestedHelp(&args))

	args = append(args, "--something")
	assert.False(t, requestedHelp(&args))

	args = append(args, "--help")
	assert.True(t, requestedHelp(&args))

	args = append(args, "--something")
	assert.True(t, requestedHelp(&args))
}

func checkErrorForTest(err error){
	if err != nil {
		log.Fatalln(err)
	} else {
		return
	}
}
