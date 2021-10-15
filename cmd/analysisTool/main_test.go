package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDateValidator(t *testing.T) {
	res, _ := validateDateInput("01-01-1111")
	assert.False(t, res)
	res, _ = validateDateInput("01-13-2021")
	assert.False(t, res)
	res, _ = validateDateInput("32-10-2021")
	assert.False(t, res)
	res, _ = validateDateInput("13-10-2021")
	assert.True(t, res)
	res, _ = validateDateInput("13.10.2021")
	assert.False(t, res)
	res, _ = validateDateInput("13/10/2021")
	assert.False(t, res)
}

func TestOperationValidator(t *testing.T) {
	res, _ := validateOperationInput("")
	assert.False(t, res)
	res, _ = validateOperationInput("a")
	assert.False(t, res)
	res, _ = validateOperationInput("A")
	assert.False(t, res)
	res, _ = validateOperationInput("0")
	assert.False(t, res)
	res, _ = validateOperationInput("3")
	assert.False(t, res)
	res, _ = validateOperationInput("11")
	assert.False(t, res)
	res, _ = validateOperationInput("01")
	assert.False(t, res)
	res, _ = validateOperationInput("1")
	assert.True(t, res)
	res, _ = validateOperationInput("2")
	assert.True(t, res)

}

func TestTrimStringBasedOnOS(t *testing.T) {
	if runtime.GOOS == "windows" {
		res := trimStringBasedOnOS("teststring\r\n")
		assert.EqualValues(t, res, "teststring")
	} else {
		res := trimStringBasedOnOS("teststring\n")
		assert.EqualValues(t, res, "teststring")
	}
}

func TestContentToArray(t *testing.T) {
	var content = strings.Split("address1,name1;\naddress2,name2;\naddress3,name3;\n", "\n")
	contentArray := *contentToArray(content)
	assert.EqualValues(t, contentArray[0][0], "address1")
	assert.EqualValues(t, contentArray[0][1], "name1")
	assert.EqualValues(t, contentArray[1][0], "address2")
	assert.EqualValues(t, contentArray[1][1], "name2")
	assert.EqualValues(t, contentArray[2][0], "address3")
	assert.EqualValues(t, contentArray[2][1], "name3")
}

func TestReadDataFromFile(t *testing.T) {
	in := "value1-x-y-z;\nvalue2.,!?;\nvalue3\t;\n"
	expected := strings.Split(in, "\n")
	filePath := filepath.Join("../../logs/temporaryForTest.txt")
	f, _ := os.Create(filePath)
	f.WriteString(in)
	defer os.Remove(filePath)
	defer f.Close()

	out := *readDataFromFile(filePath)
	for i := 0; i < len(out)-1; i++ {
		assert.EqualValues(t, expected[i], out[i])
	}
}
