package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckFlagFunctionality(t *testing.T){
	datePtr := "2021-10-29"
	operationPtr := "Visitor"
	var selectedOperation Operation
	queryPtr := "QueryWord"

	res, fails := checkFlagFunctionality(&datePtr, &operationPtr, &selectedOperation, &queryPtr)
	assert.True(t, res)
	assert.EqualValues(t, 0, len(*fails))
	assert.EqualValues(t, string(selectedOperation), VISITOR)

	operationPtr = "Location"
	res, fails = checkFlagFunctionality(&datePtr, &operationPtr, &selectedOperation, &queryPtr)
	assert.True(t, res)
	assert.EqualValues(t, 0, len(*fails))
	assert.EqualValues(t, string(selectedOperation), LOCATION)

	datePtr = "2021-13-29"
	res, fails = checkFlagFunctionality(&datePtr, &operationPtr, &selectedOperation, &queryPtr)
	assert.False(t, res)
	assert.EqualValues(t, 1, len(*fails))

	operationPtr = "somethingDifferent"
	res, fails = checkFlagFunctionality(&datePtr, &operationPtr, &selectedOperation, &queryPtr)
	assert.False(t, res)
	assert.EqualValues(t, 2, len(*fails))

	queryPtr = ""
	res, fails = checkFlagFunctionality(&datePtr, &operationPtr, &selectedOperation, &queryPtr)
	assert.False(t, res)
	assert.EqualValues(t, 3, len(*fails))
}

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
	res, _ = validateOperationInput("person")
	assert.False(t, res)
	res, _ = validateOperationInput("visitor")
	assert.True(t, res)
	res, _ = validateOperationInput("location")
	assert.True(t, res)
	res, _ = validateOperationInput("Visitor")
	assert.True(t, res)
	res, _ = validateOperationInput("Location")
	assert.True(t, res)
	res, _ = validateOperationInput("viSItor")
	assert.True(t, res)
	res, _ = validateOperationInput("locATion")
	assert.True(t, res)
}

func TestYesNoValidator(t *testing.T){
	res, _ := validateYesNoInput("")
	assert.False(t, res)
	res, _ = validateYesNoInput("something")
	assert.False(t, res)
	res, _ = validateYesNoInput("yess")
	assert.False(t, res)
	res, _ = validateYesNoInput("noo")
	assert.False(t, res)
	res, _ = validateYesNoInput("0")
	assert.False(t, res)
	res, _ = validateYesNoInput("1")
	assert.False(t, res)
	res, _ = validateYesNoInput("es")
	assert.False(t, res)
	res, _ = validateYesNoInput("o")
	assert.False(t, res)
	res, _ = validateYesNoInput("n")
	assert.True(t, res)
	res, _ = validateYesNoInput("Y")
	assert.True(t, res)
	res, _ = validateYesNoInput("N")
	assert.True(t, res)
	res, _ = validateYesNoInput("Yes")
	assert.True(t, res)
	res, _ = validateYesNoInput("No")
	assert.True(t, res)
	res, _ = validateYesNoInput("yes")
	assert.True(t, res)
	res, _ = validateYesNoInput("no")
	assert.True(t, res)
	res, _ = validateYesNoInput("yEs")
	assert.True(t, res)
	res, _ = validateYesNoInput("nO")
	assert.True(t, res)
}

func TestValidateQueryInput(t *testing.T){
	res := validateQueryInput("")
	assert.False(t, res)
	res = validateQueryInput("something")
	assert.True(t, res)
	res = validateQueryInput("something else")
	assert.True(t, res)
	res = validateQueryInput("4206942")
	assert.True(t, res)
}
