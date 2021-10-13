package main

import (
	"github.com/stretchr/testify/assert"
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
