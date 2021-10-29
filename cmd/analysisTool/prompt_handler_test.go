package main

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"runtime"
	"testing"
)

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

func TestDateInputHandler(t *testing.T) {
	var reader bufio.Reader
	in := "2021-10-28\n"
	writer := bufio.NewWriter(os.Stdin)
	accessed := make(chan bool)
	go func() {
		<-accessed
		num, err := writer.WriteString(in)
		fmt.Println(num)
		checkError(err)
		err = writer.Flush()
		checkError(err)
		fmt.Println("flushed")
	}()

	out := dateInputHandler(&reader, accessed)
	assert.EqualValues(t, out, in)
}

func checkError(err error){
	if err != nil {
		return
	} else {
		log.Fatalln(err)
	}
}
