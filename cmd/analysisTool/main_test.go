package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestContentToArray(t *testing.T) {
	var content = strings.Split("LOGIN,name1,address1,location1,20-10-2021 09:44:25,20-10-2021 09:44:25;\n"+
		"LOGIN,name2,address2,location2,20-10-2021 09:44:41,20-10-2021 09:44:41;", "\n")
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

func BenchmarkPerformanceOfData(b *testing.B) {
	fileContent := "LOGIN,name,address,location,20-10-2021 09:44:25,20-10-2021 09:44:25;\nLOGIN,name,address,location,20-10-2021 09:44:41,20-10-2021 09:44:41;\nLOGIN,name,address,location,20-10-2021 10:07:13,20-10-2021 10:07:13;\nLOGIN,name,address,location,20-10-2021 10:07:18,20-10-2021 10:07:18;\nLOGIN,name,address,location,20-10-2021 10:07:28,20-10-2021 10:07:28;\nLOGIN,name,address,location,20-10-2021 10:07:33,20-10-2021 10:07:33;\nLOGIN,name,address,location,20-10-2021 10:07:33,20-10-2021 10:07:33;"
	for n := 0; n < b.N; n++ {
		content := strings.Split(fileContent, "\n")
		contentToArray(&content)
	}
}

func TestAnalyseLocationsByVisitor(t *testing.T){

}

func TestAnalyseVisitorsByLocation(t *testing.T){

}

func TestAssertQueryExport(t *testing.T){

}
