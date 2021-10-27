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
	res, _ := validateDateInput("111-01-01")
	assert.False(t, res)
	res, _ = validateDateInput("2021-13-01")
	assert.False(t, res)
	res, _ = validateDateInput("2021-10-32")
	assert.False(t, res)
	res, _ = validateDateInput("2021-10-13")
	assert.True(t, res)
	res, _ = validateDateInput("2021.10.13")
	assert.False(t, res)
	res, _ = validateDateInput("2021/10/13")
	assert.False(t, res)
	res, _ = validateDateInput("2021-10-22")
	assert.True(t, res)
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
		res := trimStringBasedOnOS("teststring\r\n", true)
		assert.EqualValues(t, res, "teststring" )
	} else {
		res := trimStringBasedOnOS("teststring\n", true)
		assert.EqualValues(t, res, "teststring")
	}
	res := trimStringBasedOnOS("\nteststring", false)
	assert.EqualValues(t, res, "teststring" )
}

func TestContentToArray(t *testing.T) {
	var content = strings.Split("in,name1,address1,location1,20-10-2021 09:44:25,20-10-2021 09:44:25;\nin,name2,address2,location2,20-10-2021 09:44:41,20-10-2021 09:44:41;\n", "\n")
	contentArray := *contentToArray(&content)
	assert.EqualValues(t, contentArray[0].Login, true)
	assert.EqualValues(t, contentArray[0].Name, "name1")
	assert.EqualValues(t, contentArray[0].Address, "address1")
	assert.EqualValues(t, contentArray[0].Location, "location1")
	assert.EqualValues(t, contentArray[0].TimeCome.Format(DATEFORMATWITHTIME), "20-10-2021 09:44:25")
	assert.EqualValues(t, contentArray[0].TimeGone.Format(DATEFORMATWITHTIME), "20-10-2021 09:44:25")
	assert.EqualValues(t, contentArray[1].Login, true)
	assert.EqualValues(t, contentArray[1].Name, "name2")
	assert.EqualValues(t, contentArray[1].Address, "address2")
	assert.EqualValues(t, contentArray[1].Location, "location2")
	assert.EqualValues(t, contentArray[1].TimeCome.Format(DATEFORMATWITHTIME), "20-10-2021 09:44:41")
	assert.EqualValues(t, contentArray[1].TimeGone.Format(DATEFORMATWITHTIME), "20-10-2021 09:44:41")
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

func BenchmarkPerformanceOfData(b *testing.B) {
	fileContent := "name,address,location,20-10-2021 09:44:25,20-10-2021 09:44:25;\nname,address,location,20-10-2021 09:44:41,20-10-2021 09:44:41;\nname,address,location,20-10-2021 10:07:13,20-10-2021 10:07:13;\nname,address,location,20-10-2021 10:07:18,20-10-2021 10:07:18;\nname,address,location,20-10-2021 10:07:28,20-10-2021 10:07:28;\nname,address,location,20-10-2021 10:07:33,20-10-2021 10:07:33;\nname,address,location,20-10-2021 10:07:33,20-10-2021 10:07:33;"
	for n := 0; n < b.N; n++ {
		content := strings.Split(fileContent, "\n")
		contentToArray(&content)
	}
}
